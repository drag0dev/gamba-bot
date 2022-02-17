package csgocases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var BASE_URL string
var BEARER_TOKEN string
var CSGOCASES_ID string = "943452686820761600"
var DB_NAME_WEBSITES string

type meta struct{
    Oldest_id string `json:"oldest_id"`
    Newest_id string `json:"newest_id"`
    Result_count int `result_count:"result_count"`
    Next_token string `next_token:"next_token"`
}

type individualTweet struct{
    Id string `json:"id"`
    Text string `json:"text"`
}

type resBody struct{
    Meta meta `json:"meta"`
    Data []individualTweet `json:"data"`
}

func getlastOldestId(db *sql.DB) (error, string) {
    var selectStm string = fmt.Sprintf("SELECT lastid FROM %s WHERE name='csgocases'", DB_NAME_WEBSITES)

    var oldId string

    err := db.QueryRow(selectStm).Scan(&oldId)
    if err != nil{
        return err, ""
    }else{
        return nil, oldId
    }
}

func Scrape(db *sql.DB) error {
    err := godotenv.Load(".env")
    if err != nil{
        log.Print("Cannot load .env for scraper!")
        return err
    }
    BASE_URL = os.Getenv("TWITTER_BASE_URL")
    BEARER_TOKEN = os.Getenv("BEARER_TOKEN")
    DB_NAME_WEBSITES = os.Getenv("DB_NAME_WEBSITES")
    err, oldOldestId  := getlastOldestId(db)

    if err != nil{
        return err
    }

    // fetching tweets
    client := http.Client{}

    req, err := http.NewRequest("get", BASE_URL + "users/" + CSGOCASES_ID + "/tweets?max_results=20", nil)
    if err != nil{
        log.Print("Problem with making a new request!")
        return err
    }

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))

    res, err := client.Do(req)
    if err != nil{
        log.Printf("Failed to fetch csgocases, error: %s\n", err)
        return err
    }

    body, err := ioutil.ReadAll(res.Body)

    if err != nil{
        log.Printf("Failed to parse body, error: %s\n", err)
        return err
    }

    var resJSON resBody

    err = json.Unmarshal([]byte(body), &resJSON)

    if err!=nil{
        return err
    }

    // no new tweets
    if resJSON.Meta.Oldest_id <= oldOldestId{
        return nil
    }

    // going through new tweets
    var tweets []individualTweet = resJSON.Data

    for _, tweet := range tweets{
        if tweet.Id <= oldOldestId{
            break
        }

        if strings.Contains(tweet.Text, "promocode will get free $"){
            log.Print(tweet.Text)
        }
    }

    return nil
}
