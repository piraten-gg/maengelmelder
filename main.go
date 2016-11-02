package main

import (
	"fmt"
	"github.com/codedust/go-httpserve"
	"github.com/piraten-gg/maengelmelder/persistence"
	"log"
	"net/http"
)

// the global connection to the database
var storage *persistence.StorageConn

func main() {
	var err error

	// load the database
	if storage, err = persistence.Open("./data/storage.sqlite3"); err != nil {
		log.Fatal("Creating storage failed", err.Error())
	}

	// setup web server
	mux := http.NewServeMux()
	mux.Handle("/events", http.Handler(handleWS))
	mux.Handle("/api/", http.Handler(handleAPI))
	mux.Handle("/", http.FileServer(http.Dir(CFG_HTML_DIR)))

	// create or load certificate
	CreateCertificateIfNotExist(
		CFG_DATA_DIR+CFG_CERT_PREFIX+"cert.pem",
		CFG_DATA_DIR+CFG_CERT_PREFIX+"key.pem",
		CFG_SERVER_FQDN,
		3072)

	fmt.Println("Starting webserver at https://0.0.0.0:" + CFG_SERVER_PORT)

	// start the web server
	err = httpserve.ListenAndUpgradeTLS(
		"0.0.0.0:"+CFG_SERVER_PORT,
		CFG_DATA_DIR+CFG_CERT_PREFIX+"cert.pem",
		CFG_DATA_DIR+CFG_CERT_PREFIX+"key.pem",
		mux)

	log.Fatal(err)
}
