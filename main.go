package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var Token string

func init(){
    var err error
    err = godotenv.Load(".env")

    if err != nil {
        log.Printf("Error opening .env file, %s\n", err)
        return
    }

    Token = os.Getenv("DG_TOKEN")
}

func messageReply(s *discordgo.Session,m *discordgo.MessageCreate){
    if m.Author.ID == s.State.User.ID{
        return
    }

    if m.Content == "!subscribe"{
        // mock response until db is made for storing subscribed users
        s.ChannelMessageSend(m.ChannelID, m.Author.Username + " sucessfully subscribed!")
    }

    if m.Content == "!unsubscribe"{
        s.ChannelMessageSend(m.ChannelID, m.Author.Username + " successfully unsubscribed!")
    }
}

func main (){
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
