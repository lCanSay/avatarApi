package main

import (
	"errors"
	"fmt"
	"net/http"

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
		Abilities      int    `json:"abilities"`
		Image          string `json:"image"`
		Affiliation_id int    `json:"affiliation_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v := validator.New()
	ability, err := app.models.Abilities.GetByID(input.Abilities) // Retrieve ability by ID
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	character := &models.Character{
		Name:           input.Name,
		Age:            input.Age,
		Gender:         input.Gender,
		Abilities:      ability.Name,
		Image:          input.Image,
		Affiliation_id: input.Affiliation_id,
	}

	models.ValidateCharacter(v, character)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Characters.Insert(character, input.Abilities)
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
		Abilities     *int    `json:"abilities"`
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
		// Retrieve the ability name by its ID
		ability, err := app.models.Abilities.GetByID(*input.Abilities)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		character.Abilities = ability.Name
	}

	if input.Image != nil {
		character.Image = *input.Image
	}

	if input.AffiliationID != nil {
		character.Affiliation_id = *input.AffiliationID
	}

	err = app.models.Characters.Update(character, *input.Abilities)
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

func (app *application) CreateAffiliationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Image       string `json:"image"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v := validator.New()

	affiliation := &models.Affiliation{
		Name:        input.Name,
		Image:       input.Image,
		Description: input.Description,
	}

	models.ValidateAffiliation(v, affiliation)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Affiliations.Insert(affiliation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	// Check if the permission already exists for the user
	exists, err := app.models.Permissions.CheckForUser(user.ID, "affiliations:write")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		// Add the permission only if it doesn't already exist
		err = app.models.Permissions.AddForUser(user.ID, "affiliations:write")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.writeJSON(w, http.StatusCreated, envelope{"affiliation": affiliation}, nil)
}

func (app *application) GetAffiliationsListHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		models.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Extract query parameters for name, page, page size, and sort.
	input.Name = app.readStrings(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Define the sort safe list for affiliations.
	input.Filters.SortSafeList = []string{
		// Ascending sort values
		"id", "name",
		// Descending sort values
		"-id", "-name",
	}

	// Validate the input filters.
	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve affiliations from the database using the provided filters.
	affiliations, metadata, err := app.models.Affiliations.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the response with affiliations and metadata.
	app.writeJSON(w, http.StatusOK, envelope{"affiliations": affiliations, "metadata": metadata}, nil)
}

func (app *application) GetAffiliationByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	affiliation, err := app.models.Affiliations.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"affiliation": affiliation}, nil)
}

func (app *application) DeleteAffiliationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Affiliations.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "affiliation deleted successfully"}, nil)
}

func (app *application) UpdateAffiliationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	affiliation, err := app.models.Affiliations.GetByID(id)
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
		Name        *string `json:"name"`
		Image       *string `json:"image"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		affiliation.Name = *input.Name
	}

	if input.Image != nil {
		affiliation.Image = *input.Image
	}

	if input.Description != nil {
		affiliation.Description = *input.Description
	}

	err = app.models.Affiliations.Update(affiliation)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"affiliation": affiliation}, nil)
}

func (app *application) GetCharactersByAffiliationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	characters, err := app.models.Characters.GetByAffiliationID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
}

// ability handlers

func (app *application) CreateAbilityHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Element     string `json:"element"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v := validator.New()

	ability := &models.Ability{
		Name:        input.Name,
		Element:     input.Element,
		Description: input.Description,
		Image:       input.Image,
	}

	models.ValidateAbility(v, ability)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Abilities.Insert(ability)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	// Check if the permission already exists for the user
	exists, err := app.models.Permissions.CheckForUser(user.ID, "abilities:write")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		// Add the permission only if it doesn't already exist
		err = app.models.Permissions.AddForUser(user.ID, "abilities:write")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.writeJSON(w, http.StatusCreated, envelope{"ability": ability}, nil)
}

// GetAbilitiesListHandler handles listing abilities with optional filters.
func (app *application) GetAbilitiesListHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Element string
		models.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.Element = app.readStrings(qs, "element", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "name", "element",
		"-id", "-name", "-element",
	}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	abilities, metadata, err := app.models.Abilities.GetAll(input.Name, input.Element, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"abilities": abilities, "metadata": metadata}, nil)
}

// GetAbilityByIdHandler handles retrieving an ability by its ID.
func (app *application) GetAbilityByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	ability, err := app.models.Abilities.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"ability": ability}, nil)
}

// DeleteAbilityHandler handles the deletion of an ability by its ID.
func (app *application) DeleteAbilityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Abilities.Delete(id)
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

// UpdateAbilityHandler handles updating an ability by its ID.
func (app *application) UpdateAbilityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	ability, err := app.models.Abilities.GetByID(id)
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
		Name        *string `json:"name"`
		Element     *string `json:"element"`
		Description *string `json:"description"`
		Image       *string `json:"image"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		ability.Name = *input.Name
	}

	if input.Element != nil {
		ability.Element = *input.Element
	}

	if input.Description != nil {
		ability.Description = *input.Description
	}

	if input.Image != nil {
		ability.Image = *input.Image
	}

	err = app.models.Abilities.Update(ability)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"ability": ability}, nil)
}

func (app *application) GetCharactersByAbilityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	characters, err := app.models.Characters.GetByAbilityID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
}
