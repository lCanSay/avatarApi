package models

import (
	"database/sql"
	"log"
)

type Character struct {
	Id          int    `json:"id"`
	Name        string `json:"fname"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Affiliation string `json:"affiliation"`
	Abilities   string `json:"abilities"` // elements or technics
	Image       string `json:"image"`
}

var Characters = []Character{
	{Id: 1, Name: "Aang", Age: 112, Gender: "Male", Affiliation: "Air Nomads", Abilities: "Airbending, Energybending", Image: "https://example.com/aang.jpg"},
	{Id: 2, Name: "Katara", Age: 14, Gender: "Female", Affiliation: "Water Tribe", Abilities: "Waterbending, Healing", Image: "https://example.com/katara.jpg"},
	{Id: 3, Name: "Zuko", Age: 16, Gender: "Male", Affiliation: "Fire Nation", Abilities: "Firebending", Image: "https://example.com/zuko.jpg"},
	// will add more later
}

func InsertCharacter(db *sql.DB, character Character) error {
	query := `INSERT INTO character (id, name, age, gender, affiliation, abilities, image)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(query, character.Id, character.Name, character.Age, character.Gender, character.Affiliation, character.Abilities, character.Image)
	if err != nil {
		log.Fatal("query error")
		return err
	}

	return nil
}

func GetAllCharacters(db *sql.DB) ([]Character, error) {
	query := "SELECT * FROM character"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []Character
	for rows.Next() {
		var character Character
		err := rows.Scan(&character.Id, &character.Name, &character.Age, &character.Gender, &character.Affiliation, &character.Abilities, &character.Image)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

func GetCharacterByID(db *sql.DB, id int) (*Character, error) {
	query := "SELECT * FROM character WHERE id = $1"
	row := db.QueryRow(query, id)

	var character Character
	err := row.Scan(&character.Id, &character.Name, &character.Age, &character.Gender, &character.Affiliation, &character.Abilities, &character.Image)
	if err != nil {
		return nil, err
	}

	return &character, nil
}

func DeleteCharacter(db *sql.DB, id int) error {
	query := `
        DELETE FROM character
        WHERE id = $1
    `
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCharacter(db *sql.DB, character Character) error {
	query := `
        UPDATE character
        SET name = $2, age = $3, gender = $4, affiliation = $5, abilities = $6, image = $7
        WHERE id = $1
    `
	_, err := db.Exec(query, character.Id, character.Name, character.Age, character.Gender, character.Affiliation, character.Abilities, character.Image)
	if err != nil {
		return err
	}

	return nil
}
