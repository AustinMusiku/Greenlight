package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new router.
	router := httprouter.New()

	// Convert the error helpers to http.Handler(s) using the http.HandlerFunc() adapter,
	// and then set them as the custom error handlers for 404 and 405 responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.
	router.GET("/v1/healthcheck", app.healthcheckHandler)

	router.POST("/v1/movies", app.createMovieHandler)
	router.GET("/v1/movies/:id", app.showMovieHandler)

	// Return the httprouter instance.
	return router
}
