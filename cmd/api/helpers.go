package main

import (
	"encoding/json"
	"net/http"
)

// helper for sending JSON responses. This takes:
// destination http.ResponseWriter,
// HTTP status code to send,
// data to encode to JSON,
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) WriteJson(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Loop through the header map argument and add each header
	// to the http.ResponseWriter header map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
