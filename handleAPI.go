package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"log"
	"math/rand"
	"time"
)

func ApiPreMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	w.Header().Set("Cache-Control", "no-cache")

	if storage == nil {
		log.Print("[handleAPI] ERROR: storage is nil.")
		rejectWithDefaultErrorJSON(w)
		return
	}
	next(w, r)
}

func getMarkers(w web.ResponseWriter, r *web.Request) {
	markers := storage.GetMarkers()
	if markers == nil {
		fmt.Println("Error fetching markers.")
	}
	jsonMarkers, _ := json.Marshal(markers)
	fmt.Fprintf(w, string(jsonMarkers))
}

func getNewMarker(w web.ResponseWriter, r *web.Request) {
	// TODO: remove this (for testing purposes only)
	type Issue struct {
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
	}

	var issue Issue

	random := rand.New(rand.NewSource(time.Now().Unix()))

	issue.Lat = 49.9008 - 0.2 + random.Float32()*0.36
	issue.Lon = 8.3500 - 0.05 + random.Float32()*0.3

	osmmarker, err := fetchGeoInformation(issue.Lat, issue.Lon,
		16)

	if err != nil {
		fmt.Println("Fetching geo information failed.")
	}

	if osmmarker.Address.County == CFG_OWN_COUNTY {
		if err = storage.StoreMarker(issue.Lat, issue.Lon, osmmarker.DisplayName); err != nil {
			fmt.Println("Error storing marker", err.Error())
		}
	} else {
		rejectWithErrorJSON(w, "invalidmarker", "Marker outside of boundaries.")
		return
	}

	a, _ := json.Marshal(issue)
	fmt.Fprintf(w, string(a))
}

func newMarker(w web.ResponseWriter, r *web.Request) {
	requestContent := make([]byte, r.ContentLength)
	bytesRead, err := r.Body.Read(requestContent)
	if err != nil && bytesRead == 0 {
		rejectWithDefaultErrorJSON(w)
		return
	}

	type marker struct {
		Lat          float32 `json:"lat"`
		Lon          float32 `json:"lon"`
		Category     int     `json:"category"`
		Desc         string  `json:"desc"`
		Confidential bool    `json:"confidential"`
		UserName     string  `json:"user_name"`
		UserMail     string  `json:"user_mail"`
	}

	var data marker
	if err = json.Unmarshal(requestContent, &data); err != nil {
		rejectWithErrorJSON(w, "couldNotParse", "Could not parse data.")
		return
	}

	fmt.Println("New marker: ", data.Lat, data.Lon)
}
