package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lCanSay/avatarApi/internal/validator"
)

type Ability struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Element     string `json:"element"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type AbilityModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m AbilityModel) Insert(ability *Ability) error {
	query := `
        INSERT INTO ability (name, element, description, image) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `
	args := []interface{}{ability.Name, ability.Element, ability.Description, ability.Image}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&ability.Id)
}

func (m AbilityModel) GetByID(id int) (*Ability, error) {
	query := "SELECT * FROM ability WHERE id = $1"

	var ability Ability
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&ability.Id, &ability.Name, &ability.Element, &ability.Description, &ability.Image)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve ability with id: %v, %w", id, err)
	}

	return &ability, nil
}

func (m AbilityModel) Delete(id int) error {
	query := "DELETE FROM ability WHERE id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m AbilityModel) Update(ability *Ability) error {
	query := `
        UPDATE ability
        SET name = $1, element = $2, description = $3, image = $4
        WHERE id = $5
    `
	args := []interface{}{ability.Name, ability.Element, ability.Description, ability.Image, ability.Id}
	_, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m AbilityModel) GetAll(name string, element string, filters Filters) ([]*Ability, Metadata, error) {
	query := fmt.Sprintf(
		`
        SELECT count(*) OVER(), id, name, element, description, image
        FROM ability
        WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (LOWER(element) = LOWER($2) OR $2 = '')
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4
        `,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, element, filters.limit(), filters.offset()}

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

	var abilities []*Ability
	for rows.Next() {
		var ability Ability
		err := rows.Scan(&totalRecords, &ability.Id, &ability.Name, &ability.Element, &ability.Description, &ability.Image)
		if err != nil {
			return nil, Metadata{}, err
		}
		abilities = append(abilities, &ability)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return abilities, metadata, nil
}

func ValidateAbility(v *validator.Validator, ability *Ability) {
	// Validate ability.Name
	v.Check(ability.Name != "", "name", "must be provided")
	v.Check(len(ability.Name) <= 100, "name", "must not be more than 100 characters long")

	// Validate ability.Element
	v.Check(ability.Element != "", "element", "must be provided")
	v.Check(len(ability.Element) <= 50, "element", "must not be more than 50 characters long")

	// Validate ability.Image
	v.Check(ability.Image != "", "image", "must be provided")
	v.Check(len(ability.Image) <= 200, "image", "must not be more than 200 characters long")

	// Validate ability.Description
	v.Check(ability.Description != "", "description", "must be provided")
	v.Check(len(ability.Description) <= 500, "description", "must not be more than 500 characters long")
}
