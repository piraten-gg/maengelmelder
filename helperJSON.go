package main

import (
	"encoding/json"
	"net/http"
)

// rejectWithErrorJSON writes an error encoded as JSON to a http.ResponseWriter
// w        the http.ResponseWriter
// code     an error code that identifies the error
// message  a message explaining what went wrong (should be human readable)
func rejectWithErrorJSON(w http.ResponseWriter, code string, message string) {
	type Err struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	e := Err{Code: code, Message: message}
	jsonErr, _ := json.Marshal(e)
	http.Error(w, string(jsonErr), 422)
}

// rejectWithDefaultErrorJSON writes a default error encoded as JSON to a
// http.ResponseWriter. rejectWithDefaultErrorJSON(w) is equivalent to
// rejectWithErrorJSON(w, "unknown", "An unknown error occoured."))
// w  the http.ResponseWriter
func rejectWithDefaultErrorJSON(w http.ResponseWriter) {
	type Err struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	e := Err{Code: "unknown", Message: "An unknown error occoured."}
	jsonErr, _ := json.Marshal(e)
	http.Error(w, string(jsonErr), 422)
}

// createSimpleJSONEvent creates a simple JSON event used in a WS connection
// name  the name of the type of the event
func createSimpleJSONEvent(name string) string {
	type jsonEvent struct {
		Type string `json:"type"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type: name,
	})

	return string(e)
}
