package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// fileExists returns true if the given file or directory exists, otherwise
// false
// path   the given file or directory
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type Address struct {
	Road     string `json:"road"`
	County   string `json:"county"`
	Postcode string `json:"postcode"`
}

type OSMmarker struct {
	Lat         float32 `json:"lat,string"`
	Lon         float32 `json:"lon,string"`
	DisplayName string  `json:"display_name"`
	Address     Address `json:"address"`
}

func fetchGeoInformation(lat float64, lon float64, zoom int) (OSMmarker, error) {
	baseURL := "https://nominatim.openstreetmap.org/reverse?format=json"

	latStr := strconv.FormatFloat(lat, 'f', -1, 64)
	lonStr := strconv.FormatFloat(lon, 'f', -1, 64)
	zoomStr := strconv.Itoa(zoom)

	fmt.Println("Fetch " + baseURL + "&lat=" + latStr + "&lon=" + lonStr +
		"&zoom=" + zoomStr + "&addressdetails=1")
	r, err := http.Get(baseURL + "&lat=" + latStr + "&lon=" + lonStr + "&zoom=" +
		zoomStr + "&addressdetails=1")
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Error fetching data from OSM", err.Error())
		return OSMmarker{}, err
	}

	var response OSMmarker

	if err = json.NewDecoder(r.Body).Decode(&response); err != nil {
		fmt.Println(err.Error())
		return OSMmarker{}, err
	}

	return response, nil
}
