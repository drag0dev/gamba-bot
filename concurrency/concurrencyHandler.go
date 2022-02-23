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

func emitCodesToUsers(db *sql.DB, codes [][]string, s *discordgo.Session, errChan chan error, site string) {
    var selectStm string = fmt.Sprintf(`SELECT * from "%s";`, DB_NAME_USERS)
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

    for _, id := range idsArray{
        userChannel, err := s.UserChannelCreate(id)
        if err != nil{
            errChan <- err
            return
        }
        for _, code := range codes{
            s.ChannelMessageSend(userChannel.ID, fmt.Sprintf("%s CODE: %s (%s)", strings.ToUpper(site), code[0], code[1]))
        }
    }

    errChan <- nil
    return
}

func StartScraping (db *sql.DB, s *discordgo.Session, website string){

    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")

    for {
        log.Printf(`Scraping "%s"`, strings.ToUpper(website))
        cErr := make(chan error)
        cCodes := make(chan [][]string)
        cDone := make(chan bool)
        var errState bool = false

        go scraping.Scrape(db, cErr, cCodes, cDone, website)

        done := <- cDone
        err := <- cErr
        codes := <- cCodes

        if err != nil{
            log.Printf(`"%s" scraper encountered error: %s`,strings.ToUpper(website), err)
            errState = true
        }else if len(codes)>0 && done{
            log.Printf(`Emitting codes for "%s"`, strings.ToUpper(website))
            cEmittError := make(chan error)

            go emitCodesToUsers(db, codes, s, cEmittError, website)

            err := <- cEmittError

            if err != nil{
                log.Printf("error encounterd during emission: %s", err)
                errState = true
            }
        }

        if !errState{ // if there was no error wait for 15minutes before next check
            log.Printf(`Scraping successful for "%s", sleeping for 15min`, strings.ToUpper(website))
            time.Sleep(900 * time.Second)
        }else{
            log.Printf(`Error getting codes for "%s", sleeping for 5min`, strings.ToUpper(website))
            time.Sleep(300 * time.Second) // if there was error wait 5 minutes before retrying
        }
    }

}
