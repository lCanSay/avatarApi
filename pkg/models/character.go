package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lCanSay/avatarApi/internal/validator"
)

type Character struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Abilities      string `json:"abilities"` // elements or technics
	Image          string `json:"image"`
	Affiliation_id int    `json:"affiliation"`
}

type CharacterModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m CharacterModel) Insert(character *Character, abilityID int) error {
	query := `
		INSERT INTO character (name, age, gender, image, affiliation_id) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	args := []interface{}{character.Name, character.Age, character.Gender, character.Image, character.Affiliation_id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&character.Id)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO character_ability (character_id, ability_id)
		VALUES ($1, $2)
	`
	_, err = m.DB.ExecContext(ctx, query, character.Id, abilityID)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) GetByID(id int) (*Character, error) {
	query := `
		SELECT c.id, c.name, c.age, c.gender, c.image, c.affiliation_id, 
		       COALESCE(a.name, '') AS ability
		FROM character c
		LEFT JOIN character_ability ca ON c.id = ca.character_id
		LEFT JOIN ability a ON ca.ability_id = a.id
		WHERE c.id = $1
	`

	var character Character
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&character.Id, &character.Name, &character.Age, &character.Gender, &character.Image, &character.Affiliation_id, &character.Abilities)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve character with id: %v, %w", id, err)
	}

	return &character, nil
}

func (m CharacterModel) Delete(id int) error {
	query := "DELETE FROM character WHERE id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m CharacterModel) Update(character *Character, abilityID int) error {
	query := `
		UPDATE character
		SET name = $1, age = $2, gender = $3, image = $4, affiliation_id = $5
		WHERE id = $6
	`
	args := []interface{}{character.Name, character.Age, character.Gender, character.Image, character.Affiliation_id, character.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query = `DELETE FROM character_ability WHERE character_id = $1`
	_, err = m.DB.ExecContext(ctx, query, character.Id)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO character_ability (character_id, ability_id)
		VALUES ($1, $2)
	`
	_, err = m.DB.ExecContext(ctx, query, character.Id, abilityID)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) GetAll(name string, ageFrom, ageTo int, gender string, filters Filters) ([]*Character, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), c.id, c.name, c.age, c.gender, c.image, c.affiliation_id, 
		       COALESCE(a.name, '') AS ability
		FROM character c
		LEFT JOIN character_ability ca ON c.id = ca.character_id
		LEFT JOIN ability a ON ca.ability_id = a.id
		WHERE (LOWER(c.name) = LOWER($1) OR $1 = '')
		AND (c.age >= $2 OR $2 = 0)
		AND (c.age <= $3 OR $3 = 0)
		AND (LOWER(c.gender) = LOWER($4) OR $4 = '')
		GROUP BY c.id, a.name
		ORDER BY %s %s, c.id ASC
		LIMIT $5 OFFSET $6
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, ageFrom, ageTo, gender, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var characters []*Character
	for rows.Next() {
		var character Character
		err := rows.Scan(&totalRecords, &character.Id, &character.Name, &character.Age, &character.Gender, &character.Image, &character.Affiliation_id, &character.Abilities)
		if err != nil {
			return nil, Metadata{}, err
		}
		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return characters, metadata, nil
}

func ValidateCharacter(v *validator.Validator, character *Character) {
	// Validate character.Name
	v.Check(character.Name != "", "name", "must be provided")
	v.Check(len(character.Name) <= 100, "name", "must not be more than 100 characters long")

	// Validate character.Age
	v.Check(character.Age >= 0, "age", "must be a positive integer")
	v.Check(character.Age <= 150, "age", "must not exceed 150")

	// Validate character.Gender
	v.Check(character.Gender != "", "gender", "must be provided")
	v.Check(validator.In(character.Gender, "male", "female", "other"), "gender", "must be 'male', 'female', or 'other'")

	// Validate character.Abilities
	v.Check(character.Abilities != "", "abilities", "must be provided")
	v.Check(len(character.Abilities) <= 500, "abilities", "must not be more than 500 characters long")

	// Validate character.Image
	v.Check(character.Image != "", "image", "must be provided")
	v.Check(len(character.Image) <= 200, "image", "must not be more than 200 characters long")

	// Validate character.AffiliationID
	v.Check(character.Affiliation_id > 0, "affiliation_id", "must be a positive integer")
}
