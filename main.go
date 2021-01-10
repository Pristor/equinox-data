package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var (
	psqlInfo = "host=localhost port=5432 user=postgres password=mysupercoolpsqlpassword database=equinox sslmode=disable"

	config struct{
		db *sql.DB
	}
)

func init() {
	var err error
	config.db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	println("success")
}

func main()  {
	log.Fatal(http.ListenAndServe(":8000", routes()))
}

func routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/stations", handleNewStation).Methods(http.MethodPost)
	r.HandleFunc("/stations", handleListStations).Methods(http.MethodGet)
	r.HandleFunc("/stations/{device_id}", handleUpdateStation).Methods(http.MethodPut)

	return r
}
