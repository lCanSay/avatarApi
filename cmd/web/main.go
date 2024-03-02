package main

import (
	//"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/characters", GetCharacters).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
