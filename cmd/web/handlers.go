package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/lCanSay/avatarApi/pkg/models"
	//"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome!")
}

func GetCharacters(w http.ResponseWriter, r *http.Request) {
	//characters := api.Characters

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(models.Characters)
	// jsonResponse, err := json.Marshal(characters)
	// if err != nil {
	// 	return
	// }

	//w.Write(jsonResponse)
}
