package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

    "database/sql"
    _ "github.com/lib/pq"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var Token string
var HOST, USER, PASSWORD, DB_NAME_USERS string
var db *sql.DB

func init(){
    var err error
    err = godotenv.Load(".env")

    if err != nil {
        log.Printf("Error opening .env file, %s\n", err)
        return
    }

    Token = os.Getenv("DG_TOKEN")
    HOST = os.Getenv("HOST")
    USER = os.Getenv("USER")
    PASSWORD = os.Getenv("PASSWORD")
    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")
}

func messageReply(s *discordgo.Session,m *discordgo.MessageCreate){
    if m.Author.ID == s.State.User.ID{
        return
    }

    if m.Content == "!subscribe"{
        // mock response until db is made for storing subscribed users
        var insertStm = fmt.Sprintf(`INESRT INTO users (user_id) VALUES ('%s');`, m.Author.Username)
        log.Print(insertStm)
        _, err := db.Exec(insertStm)

        if err != nil{
            s.ChannelMessageSend(m.ChannelID ,"There was a problem subscribing user: " + m.Author.Username)
            log.Printf("Problem inserting new user, %s\n", err)
            return
        }else{
            s.ChannelMessageSend(m.ChannelID, m.Author.Username + " sucessfully subscribed!")
        }
    }

    if m.Content == "!unsubscribe"{
        s.ChannelMessageSend(m.ChannelID, m.Author.Username + " successfully unsubscribed!")
    }
}

func main (){
    var conninfo string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DB_NAME_USERS)
    var err error
    db, err = sql.Open("postgres", conninfo)

    if err != nil{
        log.Printf("Problem opening connection do the db, %s\n", err)
        return
    }

    dgSession, err := discordgo.New("Bot " + Token)
    if err != nil{
        log.Printf("Error making a new session, %s\n", err)
        return
    }

    dgSession.AddHandler(messageReply)
    dgSession.Identify.Intents = discordgo.IntentsGuildMessages

    err = dgSession.Open()

    if err != nil {
        log.Printf("Error opening connection to Discord. %s\n", err)
        os.Exit(1)
    }

    // don't know if i need this, will stay for now
    log.Printf(`Now running. Press CTRL-C to exit.`)
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

   dgSession.Close()
}
