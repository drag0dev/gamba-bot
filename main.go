package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)


var Session, _ = discordgo.New()

func init(){
    var err error
    err = godotenv.Load(".env")

    if err != nil {
        log.Printf("Error opening .env file, %s\n", err)
        return
    }

    Session.Token = os.Getenv("DG_TOKEN")
}

func main (){
    var err error
    err = Session.Open()

    if err != nil {
        log.Printf("Error opening connection to Discord. %s\n", err)
        os.Exit(1)
    }

    // don't know if i need this, will stay for now
    log.Printf(`Now running. Press CTRL-C to exit.`)
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    Session.Close()
}
