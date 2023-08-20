package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new router.
	router := httprouter.New()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	router.GET("/v1/healthcheck", app.healthcheckHandler)

	router.POST("/v1/movies", app.createMovieHandler)
	// TODO: Add a new route for /v1/movies/:id endpoint.

	// Return the httprouter instance.
	return router
}
