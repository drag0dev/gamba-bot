package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	_ "github.com/lib/pq"

	"drag0dev/gamba-bot/concurrency"

	"github.com/bwmarrin/discordgo"
)

var Token string
var DB_URL, DB_NAME_USERS string
var db *sql.DB

func init(){
    Token = os.Getenv("DG_TOKEN")
    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")
    DB_URL = os.Getenv("DB_URL")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate){
    if m.Author.ID == s.State.User.ID{
        return
    }
    if m.Content == "!subscribe"{
        go handleSubscribe(s, m)
    }else if m.Content == "!unsubscribe"{
        go handleUnsubscribe(s, m)
    }
}

func handleSubscribe(s *discordgo.Session, m *discordgo.MessageCreate){
    var existsStm = fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE user_id='%s');`, DB_NAME_USERS, m.Author.ID)
    var count bool

    _ = db.QueryRow(existsStm).Scan(&count)
    if !count{
        var insertStm = fmt.Sprintf(`INSERT INTO %s("user_id") VALUES('%s');`, DB_NAME_USERS, m.Author.ID)
        _, err := db.Exec(insertStm)

        if err != nil{
            s.ChannelMessageSend(m.ChannelID ,"There was a problem subscribing user " + m.Author.Username)
            log.Printf(`Error encountered during subscribe handling: %s`, err)
            return
        }

        _, err =s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s sucessfully subscribed!", m.Author.Username))
        if err != nil{
            log.Printf("Error encountered during subscribe sucessful message: %s", err)
        }

    }else if count{
        _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s already subscribed!", m.Author.Username))
        if err != nil{
            log.Printf("Error encountered during message send user already subscribed: %s", err)
        }

        return
    }else{
        s.ChannelMessageSend(m.ChannelID, "Internal error, please try again later!")
        return
    }
}

func handleUnsubscribe(s *discordgo.Session, m *discordgo.MessageCreate){
    var existsStm string = fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE user_id='%s');`, DB_NAME_USERS, m.Author.ID)
    var count bool

    _ = db.QueryRow(existsStm).Scan(&count)

    if !count{
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Cannot unsubscribe user %s!", m.Author.Username))
        return
    }else if count{
        var deleteStm string = fmt.Sprintf(`DELETE FROM %s WHERE user_id='%s';`, DB_NAME_USERS, m.Author.ID)
        _, err := db.Exec(deleteStm)

        if err != nil{
            s.ChannelMessageSend(m.ChannelID, "Interal server erorr, please try again!")
            log.Printf(`Error encountered during unsubscribing handling: %s`, err)
            return
        }else{
            _, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s has been unsubscribed!", m.Author.Username))
            if err != nil{
                log.Printf("Error encountered during user unsubscribed: %s", err)
            }
            return
        }

    }else{
        s.ChannelMessageSend(m.ChannelID, "Internal error, please try again!")
        return
    }
}


func main (){
    var err error
    db, err = sql.Open("postgres", DB_URL)

    if err != nil{
        log.Printf("Problem opening connection do the db, %s\n", err)
        return
    }

    dgSession, err := discordgo.New("Bot " + Token)
    if err != nil{
        log.Printf("Error making a new session, %s\n", err)
        return
    }

    dgSession.AddHandler(messageHandler)
    dgSession.Identify.Intents = discordgo.IntentsGuildMessages

    err = dgSession.Open()

    if err != nil {
        log.Printf("Error opening connection to Discord. %s\n", err)
        os.Exit(1)
    }

    go conScraping.StartScraping(db, dgSession, "csgocases")
    go conScraping.StartScraping(db, dgSession, "keydrop")

    log.Printf(`Now running. Press CTRL-C to exit.`)
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

   dgSession.Close()
}
