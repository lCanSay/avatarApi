package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")

	//api := r.PathPrefix("/api").Subrouter()

	// localhost:8080/api/characters
	r.HandleFunc("/characters", app.GetCharactersList).Methods("GET")
	//r.HandleFunc("/characters", app.CreateCharacterHandler).Methods("POST")
	r.HandleFunc("/characters", app.requirePermissions("characters:read", app.CreateCharacterHandler)).Methods("POST")

	//r.HandleFunc("/characters/{id:[0-9]+}", app.requirePermissions("characters:read", app.GetCharacterByIdHandler)).Methods("GET")
	r.HandleFunc("/characters/{id:[0-9]+}", app.GetCharacterByIdHandler).Methods("GET")

	//r.HandleFunc("/characters/{id:[0-9]+}", app.UpdateCharacterHandler).Methods("PUT")
	r.HandleFunc("/characters/{id:[0-9]+}", app.requirePermissions("characters:write", app.UpdateCharacterHandler)).Methods("PUT")

	r.HandleFunc("/characters/{id:[0-9]+}", app.requirePermissions("characters:write", app.DeleteCharacterHandler)).Methods("DELETE")
	//r.HandleFunc("/characters/{id:[0-9]+}", app.DeleteCharacterHandler).Methods("DELETE")

	// Affiliation routes
	r.HandleFunc("/affiliations", app.GetAffiliationsListHandler).Methods("GET")
	r.HandleFunc("/affiliations/{id:[0-9]+}", app.GetAffiliationByIdHandler).Methods("GET")
	r.HandleFunc("/affiliations", app.CreateAffiliationHandler).Methods("POST")
	r.HandleFunc("/affiliations/{id:[0-9]+}", app.UpdateAffiliationHandler).Methods("PUT")
	r.HandleFunc("/affiliations/{id:[0-9]+}", app.DeleteAffiliationHandler).Methods("DELETE")

	// User routes
	users1 := r.PathPrefix("").Subrouter()
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
