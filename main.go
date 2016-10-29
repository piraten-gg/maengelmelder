package main

import (
	"./persistence"
	"github.com/codedust/go-httpserve"
	"fmt"
	"log"
	"net/http"
)

// the global connection to the database
var storage *persistence.StorageConn

func main() {
	fmt.Println("Starting webserver at https://localhost:"+CFG_SERVER_PORT)

	// Start the server
	serveGUI()
}

func serveGUI() {
	var err error
	if storage, err = persistence.Open("./data/storage.sqlite3"); err != nil {
		log.Fatal("Creating storage failed", err.Error())
	}

	mux := http.NewServeMux()

	// paths that require authentication
	mux.Handle("/events", http.Handler(handleWS))
	mux.Handle("/api/", http.Handler(handleAPI))
	mux.Handle("/", http.FileServer(http.Dir(CFG_HTML_DIR)))

	CreateCertificateIfNotExist(
		CFG_DATA_DIR+CFG_CERT_PREFIX+"cert.pem",
		CFG_DATA_DIR+CFG_CERT_PREFIX+"key.pem",
		CFG_SERVER_FQDN,
		3072)

	err = httpserve.ListenAndUpgradeTLS(
		"0.0.0.0:"+CFG_SERVER_PORT,
		CFG_DATA_DIR+CFG_CERT_PREFIX+"cert.pem",
		CFG_DATA_DIR+CFG_CERT_PREFIX+"key.pem",
		mux)

	log.Fatal(err)
}
