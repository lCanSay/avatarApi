package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Affiliation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type AffiliationModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m AffiliationModel) Insert(affiliation *Affiliation) error {
	query := `
		INSERT INTO affiliation (name, description, image) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	args := []interface{}{affiliation.Name, affiliation.Description, affiliation.Image}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&affiliation.Id)
}

func (m AffiliationModel) GetByID(id int) (*Affiliation, error) {
	query := "SELECT * FROM affiliation WHERE id = $1"

	var affiliation Affiliation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&affiliation.Id, &affiliation.Name, &affiliation.Description, &affiliation.Image)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve affiliation with id: %v, %w", id, err)
	}

	return &affiliation, nil
}

func (m AffiliationModel) Delete(id int) error {
	query := "DELETE FROM affiliation WHERE id = $1"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m AffiliationModel) Update(affiliation *Affiliation) error {
	query := `
		UPDATE affiliation
		SET name = $1, description = $2, image = $3
		WHERE id = $4
	`
	args := []interface{}{affiliation.Name, affiliation.Description, affiliation.Image, affiliation.Id}
	_, err := m.DB.ExecContext(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m AffiliationModel) GetAll(name string, filters Filters) ([]*Affiliation, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, name, description, image
		FROM affiliation
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

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

	var affiliations []*Affiliation
	for rows.Next() {
		var affiliation Affiliation
		err := rows.Scan(&totalRecords, &affiliation.Id, &affiliation.Name, &affiliation.Description, &affiliation.Image)
		if err != nil {
			return nil, Metadata{}, err
		}
		affiliations = append(affiliations, &affiliation)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return affiliations, metadata, nil
}
