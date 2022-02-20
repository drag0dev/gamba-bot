package cCSGOCASES

import (
	"database/sql"
	"drag0dev/gamba-bot/scrapers"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var DB_NAME_USERS string

func emitCodesToUsers(db *sql.DB, codes [][]string, s *discordgo.Session, errChan chan error) {
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
            s.ChannelMessageSend(userChannel.ID, fmt.Sprintf("CSGOCASES CODE: %s (%s)", code[0], code[1]))
        }
    }

    errChan <- nil
    return
}

func StartCSGOCASES (db *sql.DB, s *discordgo.Session){
    err := godotenv.Load(".env")
    if err != nil {
        log.Printf("Error opening .env file, %s\n", err)
        return
    }

    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")

    for {
        log.Print("Scraping csgocases")
        cErr := make(chan error)
        cCodes := make(chan [][]string)
        cDone := make(chan bool)
        var errState bool = false

        go csgocases.Scrape(db, cErr, cCodes, cDone)

        err := <- cErr
        codes := <- cCodes
        done := <- cDone

        if err != nil{
            log.Printf("csgocases scraper encountered error: %s", err)
            errState = true
        }else if len(codes)>0 && done{
            log.Println("Emitting codes")
            cEmittError := make(chan error)

            go emitCodesToUsers(db, codes, s, cEmittError)

            err := <- cEmittError

            if err != nil{
                log.Printf("error encounterd during emission: %s", err)
                errState = true
            }
        }

        if !errState{ // if there was no error wait for 15minutes before next check
            time.Sleep(900 * time.Second)
        }else{
            time.Sleep(300 * time.Second) // if there was error wait 5 minutes before retrying
        }
    }

}
