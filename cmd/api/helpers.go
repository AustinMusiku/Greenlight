package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// helper for sending JSON responses. This takes:
// destination http.ResponseWriter,
// HTTP status code to send,
// data to encode to JSON,
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) WriteJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
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

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(ps httprouter.Params) (int64, error) {
	idParam := ps.ByName("id")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
