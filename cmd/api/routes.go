package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new router.
	router := httprouter.New()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.
	router.GET("/v1/healthcheck", app.healthcheckHandler)

	router.POST("/v1/movies", app.createMovieHandler)
	router.GET("/v1/movies/:id", app.showMovieHandler)

	// Return the httprouter instance.
	return router
}
