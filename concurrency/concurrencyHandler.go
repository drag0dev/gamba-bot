package conScraping

import (
	"database/sql"
	"drag0dev/gamba-bot/scraping"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var DB_NAME_USERS string
var DB_NAME_SITES string
var DB_NAME_CHANNELS string

func updateNewestId(db *sql.DB, id string, website string, errChan chan error){
    log.Printf(`"%s", Updating new id`, website)
    var updateStm string = fmt.Sprintf(`UPDATE %s SET lastId = '%s' WHERE name = '%s';`, DB_NAME_SITES, id, website)

    _, err := db.Exec(updateStm)

    if err != nil{
        errChan <- err
    }
    close(errChan)
}


func emitCodesToUsers(db *sql.DB, codes [][]string, s *discordgo.Session, errChan chan error, site string) {
    var offset int = 0
    for {
        var selectStm string = fmt.Sprintf(`SELECT user_id FROM "%s" LIMIT 100 OFFSET %d;`, DB_NAME_USERS, offset)
        rows, err := db.Query(selectStm)

        if err != nil{
            errChan <- err
            return
        }

        defer rows.Close()

        var idsArray []string
        for rows.Next(){
            var temp string
            err = rows.Scan(&temp)
            if err != nil{
                errChan <- err
                return
            }
            idsArray = append(idsArray, temp)
        }

        err = rows.Err()

        if err != nil{
            errChan <- err
            return
        }

        if len(idsArray) == 0{
            break
        }

        for _, id := range idsArray{
            userChannel, err := s.UserChannelCreate(id)
            if err != nil{
                log.Printf(`Error encountered during creation of channel in emit codes: %s`, err)
                errChan <- err
                return
            }

            for _, code := range codes{
                // catching error in order to clean up logs
                _, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf("%s CODE: %s (%s)", strings.ToUpper(site), code[0], code[1]))
            }
        }

        offset += 100
    }
    close(errChan)
    return
}

func emitCodesToChannels(db *sql.DB, codes [][]string, s *discordgo.Session, errChan chan error, site string){
    var offset int = 0
    for {
        var selectStm string = fmt.Sprintf(`SELECT channel_id FROM "%s" LIMIT 100 OFFSET %d;`, DB_NAME_CHANNELS, offset)
        rows, err := db.Query(selectStm)

        if err != nil{
            errChan <- err
            return
        }

        defer rows.Close()

        var idsArray []string
        for rows.Next(){
            var temp string
            err = rows.Scan(&temp)
            if err != nil{
                errChan <- err
                return
            }
            idsArray = append(idsArray, temp)
        }

        err = rows.Err()

        if err != nil{
            errChan <- err
            return
        }

        if len(idsArray)==0{
            break
        }

        for _, id := range idsArray{
            for _, code := range codes{
                // catching error in order to clean up logs
                _, err = s.ChannelMessageSend(id, fmt.Sprintf("%s CODE: %s (%s)", strings.ToUpper(site), code[0], code[1]))
            }
        }

        offset += 100
    }

    close(errChan)
    return
}

func StartScraping (db *sql.DB, s *discordgo.Session, website string){
    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")
    DB_NAME_SITES = os.Getenv("DB_NAME_WEBSITES")
    DB_NAME_CHANNELS = os.Getenv("DB_NAME_CHANNELS")
    var userEmitted bool = false
    var channelEmitted bool = false

    for {
        log.Printf(`Scraping "%s"`, strings.ToUpper(website))
        cErr := make(chan error)
        cCodes := make(chan [][]string)
        cDone := make(chan bool)
        cId := make(chan string)
        var errState bool = false

        go scraping.Scrape(db, cErr, cCodes, cDone, cId, website)

        done := <- cDone
        err := <- cErr
        codes := <- cCodes
        newestId := <- cId

        if err != nil{
            log.Printf(`"%s" scraper encountered error: %s`,strings.ToUpper(website), err)
            errState = true
        }

        if len(codes)>0 && done && !userEmitted && !errState{ // emit to users
            log.Printf(`Emitting codes to users for "%s"`, strings.ToUpper(website))
            cEmitUsersError := make(chan error)

            go emitCodesToUsers(db, codes, s, cEmitUsersError, website)

            err := <- cEmitUsersError

            if err != nil{
                log.Printf(`"%s", error encounterd during user emission: %s`, website, err)
                errState = true
            }else{
                userEmitted = true
            }
        }


        if len(codes)>0 && done && !channelEmitted && !errState{ // emit to channels
            log.Printf(`Emitting codes to channels for "%s"`, strings.ToUpper(website))
            cEmitChannelError := make(chan error)

            go emitCodesToChannels(db, codes, s, cEmitChannelError, website)

            err := <- cEmitChannelError

            if err != nil{
                log.Printf(`"%s", error encounterd during channel emission: %s`, website, err)
                errState = true
            }else{
                channelEmitted = true
            }
        }

        if !errState && len(codes)>0{
            // only change the newestid if codes were emitted succesfully
            // in case they didn't scraping is gonna be redone until codes can be emitted
            cUpdateError := make(chan error)
            go updateNewestId(db, newestId, website, cUpdateError)
            err := <- cUpdateError
            if err != nil{
                log.Printf(`"%s" error updating newest id: %s`, website, err)
                errState = true
            }else{
                channelEmitted = false
                userEmitted = false
            }
        }

        if !errState{ // if there was no error wait for 3minutes before next check
            log.Printf(`Scraping successful for "%s", sleeping for 3min`, strings.ToUpper(website))
            time.Sleep(180 * time.Second)
        }else{
            log.Printf(`Error getting codes for "%s", sleeping for 2min`, strings.ToUpper(website))
            time.Sleep(120 * time.Second) // if there was error wait 2 minutes before retrying
        }
    }

}
