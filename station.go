package main

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

const (
	createQuery = "INSERT INTO stations(device_id, location, location_change, created_at) VALUES ($1, $2, $3, $4)"
	listQuery = "SELECT device_id, location, location_change, created_at FROM stations"
)

type Station struct {
	ID              string                   `json:"ID"`
	Location        string                   `json:"location"`
	LocationChanged []map[string]interface{} `json:"location_changed"`
	CreatedAt       time.Time                `json:"created_at"`
}

func handleNewStation(w http.ResponseWriter, req *http.Request) {
	var station Station

	if err := json.NewDecoder(req.Body).Decode(&station); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	station.ID = uuid.New().String()
	station.CreatedAt = time.Now()
	station.LocationChanged = []map[string]interface{}{}

	if _, err := config.db.Exec(createQuery, station.ID, station.Location,
		nil, station.CreatedAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, http.StatusOK, station)
}

func handleListStations(w http.ResponseWriter, req *http.Request) {
	rows, err := config.db.Query(listQuery)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stations := make([]Station, 0)

	for rows.Next() {
		station, err := scan(rows)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stations = append(stations, *station)
	}

	render(w, http.StatusOK, stations)
}

func scan(rows *sql.Rows) (*Station, error) {
	var station Station
	var js []byte

	if err := rows.Scan(&station.ID, &station.Location, &js, &station.CreatedAt); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(js, &station.LocationChanged); err != nil {
		station.LocationChanged = []map[string]interface{}{}
	}

	return &station, nil
}

func render(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&v); err != nil {
		log.Println(err)
	}
}
