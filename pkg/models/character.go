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
