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
