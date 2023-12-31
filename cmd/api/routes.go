package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new router.
	router := httprouter.New()

	// Convert the error helpers to http.Handler(s) using the http.HandlerFunc() adapter,
	// and then set them as the custom error handlers for 404 and 405 responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.

	// healthcheck route
	router.GET("/v1/healthcheck", app.healthcheckHandler)

	// movie routes
	router.GET("/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.POST("/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))

	router.GET("/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.PATCH("/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.DELETE("/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	// user routes
	router.POST("/v1/users", app.registerUserHandler)
	router.PUT("/v1/users/activated", app.activateUserHandler)

	// token routes
	router.POST("/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// metrics
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
