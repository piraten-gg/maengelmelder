package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var handleAPI = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")

	if storage == nil {
		log.Print("[handleAPI] ERROR: storage is nil.")
		rejectWithDefaultErrorJSON(w)
		return
	}

	request := r.URL.Path[len("/api"):]
	log.Println("[handleAPI]", request)

	switch {
	// GET REQUESTS
	case strings.HasPrefix(request, "/"):
		switch request {
		case "/markers":
			markers := storage.GetMarkers()
			if markers == nil {
				fmt.Println("Error fetching markers.")
			}
			jsonMarkers, _ := json.Marshal(markers)
			fmt.Fprintf(w, string(jsonMarkers))

		case "/newmarker":
			// TODO: remove this (for testing purposes only)
			type Issue struct {
				Lat float32 `json:"lat"`
				Lon float32 `json:"lon"`
			}

			var issue Issue

			r := rand.New(rand.NewSource(time.Now().Unix()))

			issue.Lat = 49.9008 - 0.2 + r.Float32()*0.36
			issue.Lon = 8.3500 - 0.05 + r.Float32()*0.3

			osmmarker, err := fetchGeoInformation(issue.Lat, issue.Lon,
				16)

			if err != nil {
				fmt.Println("Fetching geo information failed.")
			}

			if osmmarker.Address.County == "Kreis Groß-Gerau" {
				if err = storage.StoreMarker(issue.Lat, issue.Lon, osmmarker.DisplayName); err != nil {
					fmt.Println("Error storing marker", err.Error())
				}
			} else {
				rejectWithErrorJSON(w, "invalidmarker", "Marker outside of boundaries.")
				return
			}

			a, _ := json.Marshal(issue)
			fmt.Fprintf(w, string(a))

		default:
			// unknown GET request
			rejectWithDefaultErrorJSON(w)
			return
		}

	// POST REQUESTS
	case strings.HasPrefix(request, "/post/"):
		requestContent := make([]byte, r.ContentLength)
		bytesRead, err := r.Body.Read(requestContent)
		if err != nil && bytesRead == 0 {
			rejectWithDefaultErrorJSON(w)
			return
		}

		switch request {
		case "/post/marker":
			type marker struct {
				Lat float32 `json:"lat"`
				Lon float32 `json:"lon"`
			}

			var data marker
			if err = json.Unmarshal(requestContent, &data); err != nil {
				rejectWithDefaultErrorJSON(w)
				return
			}

			fmt.Println("New marker: ", data.Lat, data.Lon)

		default:
			// unknown POST request
			rejectWithDefaultErrorJSON(w)
			return
		}

	default:
		// unknown API request
		rejectWithDefaultErrorJSON(w)
		return
	}
})
