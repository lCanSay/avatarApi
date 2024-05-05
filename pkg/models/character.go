package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
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

func (m CharacterModel) Insert(character *Character) error {
	query := `
		INSERT INTO character (name, age, gender, abilities, image, affiliation_id) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`
	args := []interface{}{character.Name, character.Age, character.Gender, character.Abilities, character.Image, character.Affiliation_id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&character.Id)
}

func (m CharacterModel) GetByID(id int) (*Character, error) {
	query := "SELECT * FROM character WHERE id = $1"

	var character Character
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&character.Id, &character.Name, &character.Age, &character.Gender, &character.Abilities, &character.Image, &character.Affiliation_id)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive character with id: %v, %w", id, err)
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

func (m CharacterModel) Update(character *Character) error {
	query := `
		UPDATE character
		SET name = $1, age = $2, gender = $3, abilities = $4, image = $5, affiliation_id = $6
		WHERE id = $7
	`
	args := []interface{}{character.Name, character.Age, character.Gender, character.Abilities, character.Image, character.Affiliation_id, character.Id}
	_, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) GetAll(name string, ageFrom, ageTo int, gender string, filters Filters) ([]*Character, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, name, age, gender, abilities, image, affiliation_id
		FROM character
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (age >= $2 OR $2 = 0)
		AND (age <= $3 OR $3 = 0)
		AND (LOWER(gender) = LOWER($4) OR $4 = '')
		ORDER BY %s %s, id ASC
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
		err := rows.Scan(&totalRecords, &character.Id, &character.Name, &character.Age, &character.Gender, &character.Abilities, &character.Image, &character.Affiliation_id)
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
