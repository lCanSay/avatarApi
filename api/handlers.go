package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "github.com/lCanSay/avatarApi/pkg/models"
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

func PostCharacter(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var character models.Character
	err := json.NewDecoder(r.Body).Decode(&character)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	// Insert the character into the database
	err = models.InsertCharacter(db, character)
	if err != nil {
		http.Error(w, "Failed to insert character into database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Character created successfully")
}
