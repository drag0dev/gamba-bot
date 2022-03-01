package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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
    if strings.HasPrefix(m.Content, "$subscribe"){
        go handleSubscribe(s, m)
    }else if strings.HasPrefix(m.Content, "$unsubscribe"){
        go handleUnsubscribe(s, m)
    }else if strings.HasPrefix(m.Content, "$bind"){
        go handleBind(s, m)
    }else if strings.HasPrefix(m.Content, "$unbind"){
        go handleUnbind(s, m)
    }else if strings.HasPrefix(m.Content, "$help"){
        go handleHelp(s, m)
    }else if strings.HasPrefix(m.Content, "$grab"){
        go handleGrab(s, m)
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

    // role testing
    member, err := s.State.Member(m.GuildID, m.Author.ID)
    if err != nil{
        var err2 error
        member, err2 = s.GuildMember(m.GuildID, m.Author.ID)

        if err2 != nil{
            log.Printf("Error encountered while getting state of the member: %s", err2)
            _, err3:= s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
            if err3 != nil{
             log.Printf("Error sending message handleBind/state member: %s", err3)
            }
            return
        }
    }

    guild, err := s.Guild(m.GuildID)
    if err != nil{
        log.Printf("Error encountered while getting guild: %s", err)
        _, err:= s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
        if err != nil{
         log.Printf("Error sending message handleBind/state member: %s", err)
        }
        return

    }

    // s.State.Member(guildId, userID) doesn't work for some reason
    // check if the user role has permissions to manage server or is an admin

    // this is seems to be required
    s.RLock()

    var admin bool
    for _, roleID := range member.Roles{
        for _, guildRole := range guild.Roles{
            if guildRole.ID == roleID{
                if guildRole.Permissions&discordgo.PermissionManageServer != 0 || guildRole.Permissions&discordgo.PermissionAdministrator !=0{
                    admin = true
                    break
                }
            }
        }
    }

    s.RUnlock()

    // if the user is not admin
    if !admin{
        _, err := s.ChannelMessageSend(m.ChannelID, "Only admins have permission to bind a channel!")
        if err != nil{
            log.Printf("Error sending message handleBind/user not admin: %s", err)
        }
        return
    }

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

    // role testing
    member, err := s.State.Member(m.GuildID, m.Author.ID)
    if err != nil{
        var err2 error
        member, err2 = s.GuildMember(m.GuildID, m.Author.ID)

        if err2 != nil{
            log.Printf("Error encountered while getting state of the member: %s", err2)
            _, err3:= s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
            if err3 != nil{
             log.Printf("Error sending message handleBind/state member: %s", err3)
            }
            return
        }
    }

    guild, err := s.Guild(m.GuildID)
    if err != nil{
        log.Printf("Error encountered while getting guild: %s", err)
        _, err:= s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
        if err != nil{
         log.Printf("Error sending message handleBind/state member: %s", err)
        }
        return

    }

    // s.State.Member(guildId, userID) doesn't work for some reason
    // check if the user role has permissions to manage server or is an admin

    // this is seems to be required
    s.RLock()

    var admin bool
    for _, roleID := range member.Roles{
        for _, guildRole := range guild.Roles{
            if guildRole.ID == roleID{
                if guildRole.Permissions&discordgo.PermissionManageServer != 0 || guildRole.Permissions&discordgo.PermissionAdministrator !=0{
                    admin = true
                    break
                }
            }
        }
    }

    s.RUnlock()

    // if the user is not admin
    if !admin{
        _, err := s.ChannelMessageSend(m.ChannelID, "Only admins have permission to unbind a channel!")
        if err != nil{
            log.Printf("Error sending message handleBind/user not admin: %s", err)
        }
        return
    }

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

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate){
    var message string = "$subscribe - subscribe to gamba bot\n$unsubscribe - unsubscribe from gamba bot\n$bind - send future codes in a channel this was typed in\n$unbind - stop sending codes to the channel this was typed in\n$grab - csgocases/keydrop to get last 5 codes"
    _, err := s.ChannelMessageSend(m.ChannelID, message)

    if err != nil{
        log.Printf("Error sending help message: %s", err)
        return
    }
}

func handleGrab(s *discordgo.Session, m *discordgo.MessageCreate){
    userChannel, err := s.UserChannelCreate(m.Author.ID)
    if err != nil{
        log.Printf("Error encountered during channel creation handleGrab: %s", err)
        _, err := s.ChannelMessageSend(m.ChannelID, "Internal server error, please try again!")
        if err != nil{
            log.Printf("Error encountered sending message grabHandle/channelCreation: %s", err)
            return
        }
    }

    var messageSplit []string = strings.Split(m.Content, " ")

    if len(messageSplit) != 2 {
        _, err := s.ChannelMessageSend(m.ChannelID, "Invalid command, please type the command correctly!")
        if err != nil {
            log.Printf("Error encountered sending message, handleGrab/messageSplit/err: %s", err)
            return
        }
    }

    var storeName string = strings.ToLower(messageSplit[1])

    if storeName != "csgocases" && storeName != "keydrop"{
        _, err := s.ChannelMessageSend(m.ChannelID, "Invalid website, please try a different one!")
        if err != nil{
            log.Printf("Error encountered sending message, handleGrab/invalid website: %s", err)
        }
    }else{
        var selectStm string = fmt.Sprintf(`SELECT code, url FROM %s LIMIT 5`, storeName)

        rows, err := db.Query(selectStm)

        if err != nil{
            log.Printf("Error encountered handleGrab/db query: %s", err)
            _, err := s.ChannelMessageSend(m.ChannelID, "Internal serve error, please try again!")
            if err != nil{
                log.Printf("Error encountered sending message handleGrab/dq query: %s", err)
                return
            }
            return
        }

        defer rows.Close()

        var finalMessage strings.Builder

        for rows.Next(){
            var code string
            var url string
            err = rows.Scan(&code, &url)
            finalMessage.WriteString(strings.ToUpper(storeName) + " CODE: " + code + " (" + url + ")" )
        }

        err = rows.Err()
        if err != nil{
            log.Printf("Error encountered scanning row grabHandle/after finish: %s", err)
            _, err := s.ChannelMessageSend(m.ChannelID, "Internal server error please try again!")
            if err != nil{
                log.Printf("Error sending message grabHandle/after finish/row scan: %s", err)
            }
            return
        }

        if len(finalMessage.String()) == 0{
            _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Currently there is no codes for %s!", storeName))
            if err != nil{
                log.Printf("Error sending message no codes: %s", err)
            }
            return
        }

        _, err = s.ChannelMessageSend(userChannel.ID, finalMessage.String())
        if err != nil{
            log.Printf("Error sending message grabHandle/codes: %s", err)
        }else{
            _, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`%s check your DMs!`, m.Author.Username))
            if err != nil{
                log.Printf("Error sending message successfully grabbed: %s", err)
            }
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
