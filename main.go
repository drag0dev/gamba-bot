package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"drag0dev/gamba-bot/concurrency"

	"github.com/bwmarrin/discordgo"
)

var Token string
var DB_URL, DB_NAME_USERS, DB_NAME_CHANNELS string
var db *sql.DB

func init(){
    _ = godotenv.Load(".env")
    Token = os.Getenv("DG_TOKEN")
    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")
    DB_NAME_CHANNELS = os.Getenv("DB_NAME_CHANNELS")
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
    }else if m.Content == "!bind"{ // send codes in channel that this command has been typed in
        go handleBind(s, m)
    }else if m.Content == "!unbind"{
        go handleUnbind(s, m)
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

func handleBind(s *discordgo.Session, m *discordgo.MessageCreate){
    var existsStm string = fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE channel_id='%s');`, DB_NAME_CHANNELS, m.ChannelID)
    var exists bool

    _ = db.QueryRow(existsStm).Scan(&exists)

    if !exists{
        var insertStm string = fmt.Sprintf(`INSERT INTO %s("channel_id") VALUES('%s');`, DB_NAME_CHANNELS, m.ChannelID)

        _, err := db.Exec(insertStm)

        if err != nil{
            log.Printf("Error encountered during insertion of new channel: %s", err)
            _, err = s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
            if err != nil{
                log.Printf("Error sending message internal server error/insertion of new channel: %s", err)
                return
            }
        }

        _, err = s.ChannelMessageSend(m.ChannelID, "Channel has been bound sucessfully!")
        if err != nil{
            log.Printf("Error sending message channel bound: %s", err)
        }

    }else if exists{
        _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("This channel is already been bound!"))
        if err != nil{
            log.Printf("Error sending message channel already bound: %s", err)
        }
    }else{
        _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Internal server error, please try again!"))
        if err != nil{
            log.Printf("Error sending message internal serve error/bind insert else: %s", err)
        }
    }

}

func handleUnbind(s *discordgo.Session, m *discordgo.MessageCreate){
    var existsStm string = fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE channel_id='%s');`, DB_NAME_CHANNELS, m.ChannelID)
    var exists bool

    _ = db.QueryRow(existsStm).Scan(&exists)
    if exists{
        var deleteStm string = fmt.Sprintf(`DELETE FROM %s WHERE channel_id='%s';`, DB_NAME_CHANNELS, m.ChannelID)
        _, err := db.Exec(deleteStm)

        if err != nil{
            log.Printf("Error encountered deletion of channel: %s", err)
            _, err = s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
            if err != nil{
                log.Printf("Error sending message internal server error/unbind/deltestm: %s", err)
                return
            }
        }

        _, err = s.ChannelMessageSend(m.ChannelID, "Channel unbound successfully!")
        if err != nil{
            log.Printf("Error sending message unbound successfully: %s", err)
        }

    }else if !exists{
        _, err := s.ChannelMessageSend(m.ChannelID, "Bot is not bound to this channel!")
        if err !=nil {
            log.Printf("Error sending message bot not bound: %s", err)
        }
    }else{
        _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Internal server error, please try again!"))
        if err != nil{
            log.Printf("Error sending message unbind/else: %s", err)
        }
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
