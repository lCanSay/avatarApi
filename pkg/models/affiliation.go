package models

import (
	"database/sql"
	"log"
)

type Affiliation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

// just for testing
var Affiliations = []Affiliation{
	{Name: "Air Nomads", Description: "The Air Nomads are a peaceful, nomadic society known for their spiritual connection to the elements and their mastery of airbending.", Image: "https://example.com/air_nomads.jpg"},
	{Name: "Water Tribe", Description: "The Water Tribe is a group of people living in the polar regions, known for their strong sense of community and their waterbending abilities.", Image: "https://example.com/water_tribe.jpg"},
	{Name: "Earth Kingdom", Description: "The Earth Kingdom is a large, diverse nation known for its powerful earthbending citizens and its vast, rugged landscapes.", Image: "https://example.com/earth_kingdom.jpg"},
	{Name: "Fire Nation", Description: "The Fire Nation is a militaristic nation known for its advanced technology and its firebending warriors.", Image: "https://example.com/fire_nation.jpg"},
}

func InsertAffiliation(db *sql.DB, affiliation Affiliation) error {
	query := `INSERT INTO affiliation (id, name, description, image)
        VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, affiliation.Id, affiliation.Name, affiliation.Description, affiliation.Image)
	if err != nil {
		log.Fatal("query error")
		return err
	}

	return nil
}

func GetAllAffiliations(db *sql.DB) ([]Affiliation, error) {
	query := "SELECT * FROM affiliation"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var affiliations []Affiliation
	for rows.Next() {
		var affiliation Affiliation
		err := rows.Scan(&affiliation.Id, &affiliation.Name, &affiliation.Description, &affiliation.Image)
		if err != nil {
			return nil, err
		}
		affiliations = append(affiliations, affiliation)
	}

	return affiliations, nil
}

func GetAffiliationByID(db *sql.DB, id int) (*Affiliation, error) {
	query := "SELECT * FROM affiliation WHERE id = $1"
	row := db.QueryRow(query, id)

	var affiliation Affiliation
	err := row.Scan(&affiliation.Id, &affiliation.Name, &affiliation.Description, &affiliation.Image)
	if err != nil {
		return nil, err
	}

	return &affiliation, nil
}

func DeleteAffiliation(db *sql.DB, id int) error {
	query := `
        DELETE FROM affiliation
        WHERE id = $1
    `
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAffiliation(db *sql.DB, affiliation Affiliation) error {
	query := `
        UPDATE affiliation
        SET name = $2, description = $3, image = $4
        WHERE id = $1
    `
	_, err := db.Exec(query, affiliation.Id, affiliation.Name, affiliation.Description, affiliation.Image)
	if err != nil {
		return err
	}

	return nil
}
