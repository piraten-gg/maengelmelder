package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"log"
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

func newMarker(w web.ResponseWriter, r *web.Request) {
	requestContent := make([]byte, r.ContentLength)
	bytesRead, err := r.Body.Read(requestContent)
	if err != nil && bytesRead == 0 {
		rejectWithDefaultErrorJSON(w)
		return
	}

	type marker struct {
		Lat          float64 `json:"lat"`
		Lon          float64 `json:"lon"`
		Zoom         int     `json:"zoom"`
		Category     int     `json:"category"`
		Descrption   string  `json:"descrption"`
		Confidential bool    `json:"confidential"`
		UserName     string  `json:"user_name"`
		UserMail     string  `json:"user_mail"`
	}

	var data marker
	if err = json.Unmarshal(requestContent, &data); err != nil {
		rejectWithErrorJSON(w, "could_not_parse", err.Error())
		return
	}

	fmt.Println("New marker: ", data.Lat, data.Lon)

	osmmarker, err := fetchGeoInformation(data.Lat, data.Lon, data.Zoom)
	if err != nil {
		log.Println("Fetching geo information failed.")
		rejectWithErrorJSON(w, "temporary_unavailable", "Service is temporary unavailable.")
		return
	}

	if osmmarker.Address.County == CFG_OWN_COUNTY {
		if err = storage.StoreMarker(data.Lat, data.Lon, osmmarker.DisplayName); err != nil {
			fmt.Println("Error storing marker", err.Error())
			rejectWithErrorJSON(w, "could_not_store", "Could not store marker.")
			return
		}
	} else {
		rejectWithErrorJSON(w, "out_of_boundaries", "Marker outside of boundaries.")
		return
	}

	respondJSON(w, stringMap{"code": "success"})
}
