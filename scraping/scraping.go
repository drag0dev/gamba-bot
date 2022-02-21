package scraping

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
var KEYDROP_ID string = "1271866668885643269"
var DB_NAME_siteS string
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

type ocrRes struct{
    ParsedResults []struct{
        TextOverlay struct{
            Lines []string                  `json:"Lines"`
            HasOverlay bool                 `json:"HasOverlay"`
            Message string                  `json:"Message"`
        }                                   `json:"TextOverlay"`
        TextOrientation string              `json:"TextOrientation"`
        FileParseExitCode int               `json:"FileParseExitCode"`
        ParsedText string                   `json:"ParsedText"`
        ErrorMessage string                 `json:"ErrorMessage"`
        ErrorDetails string                 `json:"ErrorDetails"`
    }                                       `json:"ParsedResults"`
    OCRExitCode int                         `json:"OCRExitCode"`
    IsErroredOnProcessing bool              `json:"IsErroredOnProcessing"`
    ProcessingTimeInMilliseconds string     `json:"ProcessingTimeInMilliseconds"`
    SearchablePDFURL string                 `json:"SearchablePDFURL"`
}

func getlastOldestId(db *sql.DB, website string) (error, string) {
    var selectStm string = fmt.Sprintf("SELECT lastid FROM %s WHERE name='%s'", DB_NAME_siteS, website)

    var oldId string

    err := db.QueryRow(selectStm).Scan(&oldId)
    if err != nil{
        return err, ""
    }else{
        return nil, oldId
    }
}

func getCodes(tweets []string) (error, [][]string){
    var imageURLs []string
    var codes [][]string

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

        var ocrResponse ocrRes

        // ocr sometimes fails and unmarshaling fails

        err = json.Unmarshal([]byte(body), &ocrResponse)
        if err != nil{
            return err, nil
        }

        var parsedCode string = ocrResponse.ParsedResults[0].ParsedText
        var code string = ""

        splitParsedCode := strings.Split(parsedCode, "\r\n")
        if len(splitParsedCode[2]) > 0{
            code = splitParsedCode[2]
        }

        temp := []string{code, url}
        codes = append(codes, temp)
    }

    return nil, codes
}

func addNewCodesToDB(db *sql.DB, codes [][]string, website string) error{
    for _, code := range codes{
        var insertStm string = fmt.Sprintf(`INSERT INTO %s("code", "url") VALUES('%s', '%s');`, website, code[0], code[1])
        _, err := db.Exec(insertStm)
        if err != nil{
            return err
        }
    }

    return nil
}

func updateNewestId(db *sql.DB, id string, website string) error{
    log.Print("Updating new id")
    var updateStm string = fmt.Sprintf(`UPDATE %s SET lastId = '%s' WHERE name = '%s';`, DB_NAME_siteS, id, website)

    _, err := db.Exec(updateStm)

    if err != nil{
        return err
    }
    return nil
}

func Scrape(db *sql.DB, errChan chan error, codesChan chan [][]string, done chan bool, site string) {
    err := godotenv.Load(".env")
    if err != nil{
        log.Print("Cannot load .env for scraper!")
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    BASE_URL_TWITTER = os.Getenv("TWITTER_BASE_URL")
    BEARER_TOKEN = os.Getenv("BEARER_TOKEN")
    DB_NAME_siteS = os.Getenv("DB_NAME_WEBSITES")
    BASE_URL_OCR = os.Getenv("OCR_BASE_URL")
    OCR_API_KEY = os.Getenv("OCR_API_KEY")

    err, oldNewestId  := getlastOldestId(db, site)

    if err != nil{
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    // fetching tweets
    client := http.Client{}

    var url string
    if site == "csgocases"{
        url = fmt.Sprintf(BASE_URL_TWITTER + "users/" + CSGOCASES_ID + "/tweets?max_results=20")
    }else if site == "keydrop"{
        url = fmt.Sprintf(BASE_URL_TWITTER + "users/" + KEYDROP_ID + "/tweets?max_results=20")
    }else{
        log.Print("Invalid website!")
        return
    }

    req, err := http.NewRequest("get", url, nil)
    if err != nil{
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))

    res, err := client.Do(req)
    if err != nil{
        log.Printf("Failed to fetch %s, error: %s\n", site, err)
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    body, err := ioutil.ReadAll(res.Body)

    if err != nil{
        log.Printf("Failed to parse body, error: %s\n", err)
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    var resJSON resBody

    err = json.Unmarshal([]byte(body), &resJSON)
    if err!=nil{
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    // no new tweets
    if resJSON.Meta.Newest_id <= oldNewestId{
        close(errChan)
        close(codesChan)
        done <- true
        return
    }

    // going through new tweets
    var tweets []individualTweet = resJSON.Data
    var tweetsIds []string

    for _, tweet := range tweets{
        if tweet.Id <= oldNewestId{
            break
        }

        if strings.Contains(tweet.Text, "promocode will get free $") && site == "csgocases"{
            tweetsIds = append(tweetsIds, tweet.Id)
        }else if strings.Contains(tweet.Text, "Golden Code") && site == "keydrop"{
            tweetsIds = append(tweetsIds, tweet.Id)
        }
    }

    if len(tweetsIds)==0{
        close(errChan)
        close(codesChan)
        done <- true
        return
    }

    err, codes := getCodes(tweetsIds)

    if err != nil {
        log.Print("Error getting codes!")
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    err = addNewCodesToDB(db, codes, site)

    if err!=nil{
        log.Print("Erorr adding new codes to the db!")
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    err = updateNewestId(db, resJSON.Meta.Newest_id, site)

    if err != nil{
        log.Print("Error updating newest id!")
        errChan <- err
        close(codesChan)
        done <- true
        return
    }

    close(errChan)
    done <- true
    codesChan <- codes
    return
}
