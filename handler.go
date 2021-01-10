package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func handleNewStation(w http.ResponseWriter, req *http.Request) {
	var station Station

	if err := json.NewDecoder(req.Body).Decode(&station); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	station.ID = uuid.New().String()
	station.CreatedAt = time.Now()
	station.LocationChanged = []LocationHistory{}

	if err := create(station); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, http.StatusOK, station)
}

func handleListStations(w http.ResponseWriter, req *http.Request) {
	stations, err := list()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, http.StatusOK, stations)
}

func handleUpdateStation(w http.ResponseWriter, req *http.Request)  {
	var station Station

	if err := json.NewDecoder(req.Body).Decode(&station); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	station.ID = mux.Vars(req)["device_id"]

	updatedStation, err := update(station)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, http.StatusOK, updatedStation)
}

func render(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&v); err != nil {
		log.Println(err)
	}
}


