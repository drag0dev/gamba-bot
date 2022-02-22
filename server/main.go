package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB
var DB_URL string

//TODO: rate limiting

// check if supplied id is in database
func exists(w http.ResponseWriter, r *http.Request){

}

// add a new user id to the database
func add(w http.ResponseWriter, r* http.Request){

}

func init(){
    err := godotenv.Load(".env")

    if err != nil {
        log.Print(err)
        return
    }
    DB_URL = os.Getenv("DB_URL")
}

func main(){
    database, err := sql.Open("postgres", DB_URL)
    if err != nil {
        log.Print(err)
        return
    }
    db = database

    http.HandleFunc("/exists", exists)
    http.HandleFunc("/add", add)

    log.Print("Starting server on port 8080!")
    log.Print(http.ListenAndServe(":8080", nil))

}
