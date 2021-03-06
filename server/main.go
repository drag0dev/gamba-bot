package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/didip/tollbooth/v6"
)

var db *sql.DB
var DB_URL string
var DB_NAME_USERS string

//TODO: rate limiting

type reqId struct{
    Id string `json:"id"`
}

type resFound struct{
    Found string `json:"found"`
}

type resError struct{
    Error string `json:"error"`
}

func setCors(w *http.ResponseWriter){
    (*w).Header().Set("Content-Type", "application/json")
    (*w).Header().Set("Access-Control-Allow-Origin", "https://affectionate-lovelace-900e2f.netlify.app")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Accept-Encoding, Content-Length")
}

// check if user is subscribed
// POST @ /exists
func exists(w http.ResponseWriter, r *http.Request){
    setCors(&w)
    if r.Method == "OPTIONS"{
        return
    }

    if r.Method != "POST"{
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    reqBody, err := ioutil.ReadAll(r.Body)

    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    r.Body.Close()

    var suppliedId reqId

    err = json.Unmarshal([]byte(reqBody), &suppliedId)

    if err != nil{
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if len(suppliedId.Id)==0{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(resError{Error: "id is required"})
        return
    }

    var selectStm string = fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE user_id='%s');", DB_NAME_USERS, suppliedId.Id)
    var found bool

    err = db.QueryRow(selectStm).Scan(&found)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("Error encountered during the check of existance: %s", err)
        return
    }

    if !found{
       json.NewEncoder(w).Encode(resFound{Found: "false"})
       return

   }else if found{
        json.NewEncoder(w).Encode(resFound{Found: "true"})
        return

   }else{
        w.WriteHeader(http.StatusInternalServerError)
        return
   }
}

// add a new user id to the database
// POST @ /exists
func subscribe(w http.ResponseWriter, r *http.Request){
    setCors(&w)
    if r.Method == "OPTIONS"{
        return
    }

    if r.Method != "POST"{
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    reqBody, err := ioutil.ReadAll(r.Body)

    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    r.Body.Close()

    var suppliedIdJSON reqId
    err = json.Unmarshal([]byte(reqBody), &suppliedIdJSON)

    if err != nil{
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // check if the id empty
    if len(suppliedIdJSON.Id)==0{
        w.WriteHeader(http.StatusBadRequest)
        err = json.NewEncoder(w).Encode(resError{Error: "id is required"})
        return
    }

    // check if user is already subscribed
    var existsStm string = fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE user_id='%s');", DB_NAME_USERS, suppliedIdJSON.Id)

    var found bool
    err = db.QueryRow(existsStm).Scan(&found)

    if err != nil{
        log.Printf("Error encountered during existance check: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if !found{
        // if user is not subscribed
        var insertStm string = fmt.Sprintf(`INSERT INTO %s("user_id") VALUES('%s');`, DB_NAME_USERS, suppliedIdJSON.Id)

        _, err = db.Exec(insertStm)
        if err != nil{
            log.Printf("Error encountered on insert new user: %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        return

    }else if found{
        w.WriteHeader(http.StatusBadRequest)
        err = json.NewEncoder(w).Encode(resError{Error: "user is already subscribed"})
        return

    }else{
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}

// unsubscribe user
// POST @ /unsubscribe
func unsubscribe(w http.ResponseWriter, r *http.Request){
    setCors(&w)
    if r.Method == "OPTIONS"{
        return
    }

    if r.Method != "POST"{
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    reqBody, err := ioutil.ReadAll(r.Body)

    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    r.Body.Close()

    var suppliedIdJSON reqId
    err = json.Unmarshal([]byte(reqBody), &suppliedIdJSON)

    if err != nil{
        w.WriteHeader(http.StatusBadGateway)
        return
    }

    if len(suppliedIdJSON.Id)==0{
        w.WriteHeader(http.StatusBadRequest)
        err = json.NewEncoder(w).Encode(resError{Error: "id is required"})
        return
    }

    // check if user is subscribed
    var existsStm string = fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE user_id='%s');", DB_NAME_USERS, suppliedIdJSON.Id)
    var found bool

    err = db.QueryRow(existsStm).Scan(&found)
    if err != nil{
        log.Printf("Error encountered during existance check: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if !found{
        w.WriteHeader(http.StatusBadRequest)
        err = json.NewEncoder(w).Encode(resError{Error: "user is not subscribed"})
        return

    }else if found{
        var delStm string = fmt.Sprintf(`DELETE FROM %s WHERE user_id='%s';`, DB_NAME_USERS, suppliedIdJSON.Id)

        _, err = db.Exec(delStm)

        if err != nil{
            log.Printf("Error encountered during deletion of users: %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        return

    }else{
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

}

func init(){
    err := godotenv.Load(".env")

    if err != nil {
        log.Print(err)
        return
    }
    DB_URL = os.Getenv("DB_URL")
    DB_NAME_USERS = os.Getenv("DB_NAME_USERS")
}

func main(){
    database, err := sql.Open("postgres", DB_URL)
    if err != nil {
        log.Print(err)
        return
    }
    db = database

    // 2 request every second for each endpoint for each user
    lmt := tollbooth.NewLimiter(2, nil)
    lmt.SetMessage("You have reached maximum request limit.")
    lmt.SetMessageContentType("text/plain; charset=utf-8")

    http.Handle("/exists", tollbooth.LimitFuncHandler(lmt, exists))
    http.Handle("/subscribe", tollbooth.LimitFuncHandler(lmt, subscribe))
    http.Handle("/unsubscribe", tollbooth.LimitFuncHandler(lmt, unsubscribe))

    port := os.Getenv("PORT")
    if port == ""{
        port = "8080"
    }

    log.Printf("Starting server on port %s!", port)
    log.Print(http.ListenAndServe(":" + port, nil))

}
