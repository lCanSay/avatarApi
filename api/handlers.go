package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/lCanSay/avatarApi/pkg/models"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome!")
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

func GetCharacters(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sql.DB)
	characters, err := models.GetAllCharacters(db)
	if err != nil {
		http.Error(w, "Failed to get characters from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

// GetCharacterById retrieves a character by its ID and sends it as a JSON response.
func GetCharacterById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	character, err := models.GetCharacterByID(db, id)
	if err != nil {
		http.Error(w, "Failed to get character from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	err = models.DeleteCharacter(db, id)
	if err != nil {
		http.Error(w, "Failed to delete character from database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Character deleted successfully")
}

func UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	var updatedCharacter models.Character
	err = json.NewDecoder(r.Body).Decode(&updatedCharacter)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedCharacter.Id = id

	db := r.Context().Value("db").(*sql.DB)
	err = models.UpdateCharacter(db, updatedCharacter)
	if err != nil {
		http.Error(w, "Failed to update character in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Character updated successfully")
}
