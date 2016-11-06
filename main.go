package main

import (
	"fmt"
	"github.com/codedust/go-httpserve"
	"github.com/gocraft/web"
	"github.com/piraten-gg/maengelmelder/persistence"
	"log"
)

// the global connection to the database
var storage *persistence.StorageConn

type Context struct{}

func main() {
	var err error

	// load the database
	if storage, err = persistence.Open("./data/storage.sqlite3"); err != nil {
		log.Fatal("Creating storage failed", err.Error())
	}

	// setup web server
	router := web.New(Context{})
	router.Middleware(web.LoggerMiddleware)
	router.Middleware(web.ShowErrorsMiddleware) // TODO: only in dev environment
	router.Middleware(ApiPreMiddleware)
	router.Middleware(web.StaticMiddleware(
		CFG_HTML_DIR, web.StaticOption{IndexFile: "index.html"}))
	router.Get("/ws", func(rw web.ResponseWriter, req *web.Request) {
		handleWS.ServeHTTP(rw, req.Request)
	})

	router.Get("/api/v1/markers/show", getMarkers)
	router.Post("/api/v1/markers/new", newMarker)

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
		router)

	log.Fatal(err)
}
