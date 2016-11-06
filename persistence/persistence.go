package persistence

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"time"
)

var (
	KeyNotFound = errors.New("Key does not exist")
)

type StorageConn struct {
	db  *sql.DB
	mux sync.Mutex
}

// Open creates a connection to the database
// always close the connection with `defer storageConn.Close()`
func Open(filename string) (*StorageConn, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
		return &StorageConn{}, err
	}

	// create database tables
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS keyValueStorage (
		key TEXT PRIMARY KEY,
		value TEXT
	);
	CREATE TABLE IF NOT EXISTS markers (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		lat	REAL,
		lon	REAL,
		displayName	TEXT,
		time	INTEGER
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Panicf("%q: %s\n", err, sqlStmt)
		return &StorageConn{}, err
	}

	s := &StorageConn{db: db}
	return s, nil
}

// Close safely closes the connection to the database
func (s *StorageConn) Close() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.db.Close()
}

// StoreKeyValue stores a key-value pair of strings
// key    the key (case sensitive)
// value  the value
func (s *StorageConn) StoreKeyValue(key string, value string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, err := s.db.Exec(`INSERT OR REPLACE INTO keyValueStorage(key, value) VALUES(?, ?)`, key, value)
	if err != nil {
		log.Print("[persistence StoreKeyValue] INSERT statement failed")
		return err
	}

	return nil
}

// GetKeyValue returns the value corresponding to the given key or an empty
// string if the value could not be determined
// key  the key (case sensitive)
func (s *StorageConn) GetKeyValue(key string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	rows, err := s.db.Query("SELECT value FROM keyValueStorage WHERE key = ?", key)
	if err != nil {
		log.Print("[persistence GetKeyValue] SELECT statement failed")
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var value string
		rows.Scan(&value)
		return value, nil
	}

	return "", KeyNotFound
}


// StoreMarker stores a marker
func (s *StorageConn) StoreMarker(lon float64, lat float64, displayName string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, err := s.db.Exec(`INSERT INTO markers(lat, lon, displayName, time) VALUES(?, ?, ?, ?)`, lon, lat, displayName, time.Now().Unix()*1000)
	if err != nil {
		log.Print("[persistence StoreKeyValue] INSERT statement failed")
		return err
	}

	return nil
}


type Marker struct {
	Lat         float32 `json:"lat,float"`
	Lon         float32 `json:"lon,float"`
	DisplayName string  `json:"display_name"`
}


// GetMarkers returns previously stored markers
func (s *StorageConn) GetMarkers() []Marker {
  s.mux.Lock()
  defer s.mux.Unlock()

	limit := -1 // no limit

  rows, err := s.db.Query("SELECT lat, lon, displayName, time FROM markers LIMIT ?", limit)
  if err != nil {
    log.Print("[persistence GetMessages] SELECT statement failed")
    return nil
  }
  defer rows.Close()

  var markers []Marker

  for rows.Next() {
		var lat float32
    var lon float32
    var displayName string
    var time int64
    rows.Scan(&lat, &lon, &displayName, &time)
    markers = append(markers, Marker{Lat: lat, Lon: lon, DisplayName: displayName})
  }

  return markers
}
