package main

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"
// 	"sort"
// 	"strconv"
// 	"strings"

// 	//"github.com/gorilla/mux"
// 	models "github.com/lCanSay/avatarApi/pkg/models"
// )

// func ListCharacters(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Name        string
// 		Affiliation string
// 		Gender      string
// 		Abilities   string
// 		Age         int
// 		Id          int
// 		Page        int
// 		PageSize    int
// 		SortOrder   string
// 		SortBy      string
// 	}
// 	qs := r.URL.Query()

// 	input.Name = readString(qs, "name", "")
// 	input.Affiliation = readString(qs, "affiliation", "")
// 	input.Gender = readString(qs, "gender", "")
// 	input.Page = readInt(qs, "page", 1)
// 	input.PageSize = readInt(qs, "page_size", 20)
// 	input.Id = readInt(qs, "id", 0)
// 	input.Age = readInt(qs, "age", 0)
// 	input.Abilities = readString(qs, "abilities", "")
// 	input.SortBy = readString(qs, "sort_by", "id")
// 	input.SortOrder = readString(qs, "sort_order", "asc")

// 	// Filter characters based on input parameters
// 	filteredCharacters, err := filterCharacters(r.Context(), input)
// 	if err != nil {
// 		http.Error(w, "Failed to get characters from database", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send the filtered characters as JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(filteredCharacters)
// }

// func filterCharacters(ctx context.Context, input struct {
// 	Name        string
// 	Affiliation string
// 	Gender      string
// 	Abilities   string
// 	Age         int
// 	Id          int
// 	Page        int
// 	PageSize    int
// 	SortOrder   string
// 	SortBy      string
// }) ([]models.Character, error) {

// 	db := ctx.Value("db").(*sql.DB)
// 	characters, err := models.GetAllCharacters(db)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var filteredCharacters []models.Character

// 	// Filter characters based on input parameters
// 	for _, character := range characters {
// 		if strings.Contains(strings.ToLower(character.Name), strings.ToLower(input.Name)) &&
// 			(strings.ToLower(character.Gender) == strings.ToLower(input.Gender) || input.Gender == "") &&
// 			(input.Age == 0 || character.Age == input.Age) &&
// 			(input.Id == 0 || character.Id == input.Id) &&
// 			(strings.Contains(strings.ToLower(character.Abilities), strings.ToLower(input.Abilities)) || input.Abilities == "") {
// 			filteredCharacters = append(filteredCharacters, character)
// 		}
// 	}

// 	// Sort filteredCharacters based on the specified field and order
// 	switch strings.ToLower(input.SortOrder) {
// 	case "asc":
// 		switch strings.ToLower(input.SortBy) {
// 		case "id":
// 			sort.Slice(filteredCharacters, func(i, j int) bool {
// 				return filteredCharacters[i].Id < filteredCharacters[j].Id
// 			})
// 		case "name":
// 			sort.Slice(filteredCharacters, func(i, j int) bool {
// 				return strings.ToLower(filteredCharacters[i].Name) < strings.ToLower(filteredCharacters[j].Name)
// 			})
// 		}

// 	case "desc":
// 		switch strings.ToLower(input.SortBy) {
// 		case "id":
// 			sort.Slice(filteredCharacters, func(i, j int) bool {
// 				return filteredCharacters[i].Id > filteredCharacters[j].Id
// 			})
// 		case "name":
// 			sort.Slice(filteredCharacters, func(i, j int) bool {
// 				return strings.ToLower(filteredCharacters[i].Name) > strings.ToLower(filteredCharacters[j].Name)
// 			})
// 		}
// 	}

// 	// Pagination
// 	start := (input.Page - 1) * input.PageSize
// 	end := start + input.PageSize
// 	if start >= len(filteredCharacters) {
// 		start = len(filteredCharacters)
// 	}
// 	if end > len(filteredCharacters) {
// 		end = len(filteredCharacters)
// 	}
// 	filteredCharacters = filteredCharacters[start:end]

// 	return filteredCharacters, nil
// }

// func readString(qs map[string][]string, key, defaultValue string) string {
// 	if val, ok := qs[key]; ok {
// 		return val[0]
// 	}
// 	return defaultValue
// }

// func readInt(qs map[string][]string, key string, defaultValue int) int {
// 	if val, ok := qs[key]; ok {
// 		if i, err := strconv.Atoi(val[0]); err == nil {
// 			return i
// 		}
// 	}
// 	return defaultValue
// }
