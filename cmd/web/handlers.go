package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lCanSay/avatarApi/internal/validator"
	models "github.com/lCanSay/avatarApi/pkg/models"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome!")
}

func (app *application) CreateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		Age            int    `json:"age"`
		Gender         string `json:"gender"`
		Abilities      string `json:"abilities"`
		Image          string `json:"image"`
		Affiliation_id int    `json:"affiliation_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	character := &models.Character{
		Name:           input.Name,
		Age:            input.Age,
		Gender:         input.Gender,
		Abilities:      input.Abilities,
		Image:          input.Image,
		Affiliation_id: input.Affiliation_id,
	}

	err = app.models.Characters.Insert(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	// Check if the permission already exists for the user
	exists, err := app.models.Permissions.CheckForUser(user.ID, "characters:write")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		// Add the permission only if it doesn't already exist
		err = app.models.Permissions.AddForUser(user.ID, "characters:write")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.writeJSON(w, http.StatusCreated, envelope{"character": character}, nil)
}

func (app *application) GetCharactersList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		AgeFrom int
		AgeTo   int
		Gender  string
		models.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Extract query parameters for name, age range, page, page size, and sort.
	input.Name = app.readStrings(qs, "name", "")
	input.AgeFrom = app.readInt(qs, "ageFrom", 0, v)
	input.AgeTo = app.readInt(qs, "ageTo", 0, v)
	input.Gender = app.readStrings(qs, "gender", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Define the sort safe list for characters.
	input.Filters.SortSafeList = []string{
		// Ascending sort values
		"id", "name", "age",
		// Descending sort values
		"-id", "-name", "-age",
	}

	// Validate the input filters.
	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve characters from the database using the provided filters.
	characters, metadata, err := app.models.Characters.GetAll(input.Name, input.AgeFrom, input.AgeTo, input.Gender, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the response with characters and metadata.
	app.writeJSON(w, http.StatusOK, envelope{"characters": characters, "metadata": metadata}, nil)
}

func (app *application) GetCharacterByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
}

func (app *application) DeleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Characters.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}

func (app *application) UpdateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name          *string `json:"name"`
		Age           *int    `json:"age"`
		Gender        *string `json:"gender"`
		Abilities     *string `json:"abilities"`
		Image         *string `json:"image"`
		AffiliationID *int    `json:"affiliation_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		character.Name = *input.Name
	}

	if input.Age != nil {
		character.Age = *input.Age
	}

	if input.Gender != nil {
		character.Gender = *input.Gender
	}

	if input.Abilities != nil {
		character.Abilities = *input.Abilities
	}

	if input.Image != nil {
		character.Image = *input.Image
	}

	if input.AffiliationID != nil {
		character.Affiliation_id = *input.AffiliationID
	}

	err = app.models.Characters.Update(character)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
}

// Affiliation Handlers-----------------------------------------------------------

func PostAffiliation(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var affiliation models.Affiliation
	err := json.NewDecoder(r.Body).Decode(&affiliation)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	// Insert the affiliation into the database
	err = models.InsertAffiliation(db, affiliation)
	if err != nil {
		http.Error(w, "Failed to insert affiliation into database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Affiliation created successfully")
}

func GetAffiliations(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /affiliations")

	db := r.Context().Value("db").(*sql.DB)
	affiliations, err := models.GetAllAffiliations(db)
	if err != nil {
		http.Error(w, "Failed to get affiliations from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(affiliations)
}

func GetAffiliationById(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /affiliations/{id}")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid affiliation ID", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	affiliation, err := models.GetAffiliationByID(db, id)
	if err != nil {
		http.Error(w, "Failed to get affiliation from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(affiliation)
}

func DeleteAffiliation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid affiliation ID", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*sql.DB)
	err = models.DeleteAffiliation(db, id)
	if err != nil {
		http.Error(w, "Failed to delete affiliation from database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Affiliation deleted successfully")
}

func UpdateAffiliation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid affiliation ID", http.StatusBadRequest)
		return
	}

	var updatedAffiliation models.Affiliation
	err = json.NewDecoder(r.Body).Decode(&updatedAffiliation)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedAffiliation.Id = id

	db := r.Context().Value("db").(*sql.DB)
	err = models.UpdateAffiliation(db, updatedAffiliation)
	if err != nil {
		http.Error(w, "Failed to update affiliation in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Affiliation updated successfully")
}
