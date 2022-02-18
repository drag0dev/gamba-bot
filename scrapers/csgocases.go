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

var BASE_URL_TWITTER string
var BEARER_TOKEN string
var CSGOCASES_ID string = "943452686820761600"
var DB_NAME_WEBSITES string
var BASE_URL_OCR string
var OCR_API_KEY string

type meta struct{
    Oldest_id string `json:"oldest_id"`
    Newest_id string `json:"newest_id"`
    Result_count int `json:"result_count"`
    Next_token string `json:"next_token"`
}

type individualTweet struct{
    Id string `json:"id"`
    Text string `json:"text"`
}

type resBody struct{
    Meta meta `json:"meta"`
    Data []individualTweet `json:"data"`
}

type tweetInfo struct{
    Data struct{
        Attachments struct{
            Media_keys []string `json:"media_keys"`
        }                       `json:"attachments"`
        Id string               `json:"id"`
        Text string             `json:"text"`
    }                           `json:"data"`

    Includes struct{
        Media []struct{
            Media_key string    `json:"media_key"`
            Type string         `json:"type"`
            Url string          `json:"url"`
        }                       `json:"media"`
    }                           `json:"includes"`
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

func getCodes(tweets []string) (error, []string){
    var imageURLs []string
    var codes []string
    client := http.Client{}

    // grabbing urls
    for _, tweet := range tweets{
        req, err := http.NewRequest("get", BASE_URL_TWITTER + "tweets/" + tweet + "?media.fields=url&expansions=attachments.media_keys", nil)
        if err != nil{
            return err, nil
        }

        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))
        res, err := client.Do(req)
        if err != nil{
            return err, nil
        }

        body, err := ioutil.ReadAll(res.Body)

        if err != nil{
            return err, nil
        }

        var resJSON tweetInfo
        err = json.Unmarshal([]byte(body), &resJSON)

        if err != nil{
            return err, nil
        }
        imageURLs = append(imageURLs, resJSON.Includes.Media[0].Url)
    }

    // passing images through ocr
    for _, url := range imageURLs{
        req, err := http.NewRequest("get", BASE_URL_OCR + fmt.Sprintf("url?apiKey=%s&url=%s",OCR_API_KEY, url), nil)
        if err != nil{
            return err, nil
        }

        res, err := client.Do(req)
        if err != nil{
            return err, nil
        }

        body, err := ioutil.ReadAll(res.Body)
        if err != nil{
            return err, nil
        }

        log.Print(string(body))
        log.Print(url)
        log.Println("")

    }

    return nil, codes
}

func Scrape(db *sql.DB) error {
    err := godotenv.Load(".env")
    if err != nil{
        log.Print("Cannot load .env for scraper!")
        return err
    }
    BASE_URL_TWITTER = os.Getenv("TWITTER_BASE_URL")
    BEARER_TOKEN = os.Getenv("BEARER_TOKEN")
    DB_NAME_WEBSITES = os.Getenv("DB_NAME_WEBSITES")
    BASE_URL_OCR = os.Getenv("OCR_BASE_URL")
    OCR_API_KEY = os.Getenv("OCR_API_KEY")

    err, oldOldestId  := getlastOldestId(db)

    if err != nil{
        return err
    }

    // fetching tweets
    client := http.Client{}

    req, err := http.NewRequest("get", BASE_URL_TWITTER + "users/" + CSGOCASES_ID + "/tweets?max_results=20", nil)
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
    var tweetsIds []string

    for _, tweet := range tweets{
        if tweet.Id <= oldOldestId{
            break
        }

        if strings.Contains(tweet.Text, "promocode will get free $"){
            tweetsIds = append(tweetsIds, tweet.Id)
        }
    }

    err, _ = getCodes(tweetsIds)

    if err != nil {
        return err
    }

    return nil
}
