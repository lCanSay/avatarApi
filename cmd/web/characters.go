package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	//"github.com/gorilla/mux"
	models "github.com/lCanSay/avatarApi/pkg/models"
)

func ListCharacters(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string
		Affiliation string
		Gender      string
		MinAge      int
		MaxAge      int
		Id          int
		Page        int
		PageSize    int
	}
	qs := r.URL.Query()

	input.Name = readString(qs, "name", "")
	input.Affiliation = readString(qs, "affiliation", "")
	input.Gender = readString(qs, "gender", "")
	input.Page = readInt(qs, "page", 1)
	input.PageSize = readInt(qs, "page_size", 20)
	input.Id = readInt(qs, "id", 1)
	input.MinAge = readInt(qs, "min_age", 0)
	input.MaxAge = readInt(qs, "max_age", 1000)

	// Filter characters based on input parameters
	filteredCharacters, err := filterCharacters(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to get characters from database", http.StatusInternalServerError)
		return
	}

	// Send the filtered characters as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredCharacters)
}

func filterCharacters(ctx context.Context, input struct {
	Name        string
	Affiliation string
	Gender      string
	MinAge      int
	MaxAge      int
	Id          int
	Page        int
	PageSize    int
}) ([]models.Character, error) {

	db := ctx.Value("db").(*sql.DB)
	characters, err := models.GetAllCharacters(db)
	if err != nil {
		return nil, err
	}

	var filteredCharacters []models.Character

	// Filter characters based on input parameters
	for _, character := range characters {
		if strings.Contains(strings.ToLower(character.Name), strings.ToLower(input.Name)) &&
			strings.Contains(strings.ToLower(character.Affiliation), strings.ToLower(input.Affiliation)) &&
			(strings.ToLower(character.Gender) == strings.ToLower(input.Gender) || input.Gender == "") &&
			(character.Age <= input.MinAge && (input.MaxAge == 0 || character.Age >= input.MaxAge)) &&
			(character.Id == input.Id || input.Id == 0) {
			filteredCharacters = append(filteredCharacters, character)
		}
	}

	// Pagination
	start := (input.Page - 1) * input.PageSize
	end := start + input.PageSize
	if start >= len(filteredCharacters) {
		start = len(filteredCharacters)
	}
	if end > len(filteredCharacters) {
		end = len(filteredCharacters)
	}
	filteredCharacters = filteredCharacters[start:end]

	return filteredCharacters, nil
}

func readString(qs map[string][]string, key, defaultValue string) string {
	if val, ok := qs[key]; ok {
		return val[0]
	}
	return defaultValue
}

func readInt(qs map[string][]string, key string, defaultValue int) int {
	if val, ok := qs[key]; ok {
		if i, err := strconv.Atoi(val[0]); err == nil {
			return i
		}
	}
	return defaultValue
}
