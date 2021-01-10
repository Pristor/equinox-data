package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	createQuery = "INSERT INTO stations(device_id, location, location_change, created_at) VALUES ($1, $2, $3, $4)"
	listQuery   = "SELECT device_id, location, location_change, created_at FROM stations"
	getQuery    = "SELECT device_id, location, location_change, created_at FROM stations WHERE device_id=$1"
	updateQuery = "UPDATE stations SET location=$1, location_change=$2 where device_id=$3"
)

type Station struct {
	ID              string            `json:"ID"`
	Location        string            `json:"location"`
	LocationChanged []LocationHistory `json:"location_change"`
	CreatedAt       time.Time         `json:"created_at"`
}

type LocationHistory struct {
	OldLocation string    `json:"old_location"`
	NewLocation string    `json:"new_location"`
	ChangedAt   time.Time `json:"changed_at"`
}

func create(station Station) error {
	_, err := config.db.Exec(createQuery, station.ID, station.Location, nil, station.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func list() (*[]Station, error) {
	rows, err := config.db.Query(listQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stations := make([]Station, 0)

	for rows.Next() {
		station, err := scan(rows)

		if err != nil {
			return nil, err
		}
		stations = append(stations, *station)
	}

	return &stations, nil
}

func scan(rows *sql.Rows) (*Station, error) {
	var station Station
	var js []byte
	var err error

	if err = rows.Scan(&station.ID, &station.Location, &js, &station.CreatedAt); err != nil {
		return nil, err
	}

	station.LocationChanged, err = fromDB(js)

	if err != nil {
		return nil, err
	}

	return &station, nil
}

func get(id string) (*Station, error) {
	rows, err := config.db.Query(getQuery, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("not found")
	}

	station, err := scan(rows)

	if err != nil {
		return nil, err
	}

	return station, nil

}

func update(station Station) (*Station, error) {
	oldStation, err := get(station.ID)

	if err != nil {
		return nil, err
	}

	if station.Location == oldStation.Location {
		return oldStation, nil
	}

	station.LocationChanged = oldStation.LocationChanged
	station.CreatedAt = oldStation.CreatedAt

	station.LocationChanged = append(station.LocationChanged, LocationHistory{
		OldLocation: oldStation.Location,
		NewLocation: station.Location,
		ChangedAt:   time.Now(),
	})

	js, err := toDb(station.LocationChanged)
	if err != nil {
		return nil, err
	}

	if _, err := config.db.Exec(updateQuery, station.Location, js, station.ID); err != nil {
		return nil, err
	}

	return &station, nil
}

func toDb(v []LocationHistory) (string, error) {
	var tempStr []string

	for _, str := range v {
		b, err := json.Marshal(str)
		if err != nil {
			return "", err
		}
		tempStr = append(tempStr, strconv.Quote(string(b)))
	}
	finalStr := "{" + strings.Join(tempStr[:], ",") + "}"
	return finalStr, nil
}

func fromDB(js []byte) ([]LocationHistory, error) {
	var locationChange []LocationHistory

	s := string(js)
	newStr := postgreArrayToGoArray(s)
	b := []byte(newStr)

	err := json.Unmarshal(b, &locationChange)
	if err != nil {
		return nil, err
	}

	return locationChange, nil
}

func postgreArrayToGoArray(data string) string {
	if len(data) == 0 {
		return "[]"
	}
	goArrayString := fmt.Sprintf("[%s]", data[1:len(data)-1])
	goArrayString = strings.ReplaceAll(goArrayString, "\"{", "{")
	goArrayString = strings.ReplaceAll(goArrayString, "}\"", "}")
	goArrayString = strings.ReplaceAll(goArrayString, "\\\"", "\"")
	return goArrayString
}
