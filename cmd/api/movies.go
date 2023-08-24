package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AustinMusiku/Greenlight/internal/data"
	"github.com/AustinMusiku/Greenlight/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// This struct will be our *target decode destination*.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new Validator instance.
	v := validator.New()

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Include a Location header to let the client know
	// which URL they can find the newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := app.readIDParam(ps)
	if err != nil { // invalid id parameter
		app.logger.Println(err)
		http.NotFound(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	// Use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := app.readIDParam(ps)
	if err != nil { // invalid id parameter
		app.logger.Println(err)
		http.NotFound(w, r)
		return
	}

	// Fetch the existing movie record from the database, sending a 404 Not Found
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// This struct will be our *target decode destination*.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Read the JSON request body with updated data into the input struct.
	app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// overwrite the new values over the existing values in the movie record.
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	// Initialize a new Validator instance.
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the movie record in the database, checking for any edit conflicts
	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return a 200 OK status code along with the movie data in the response body.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
