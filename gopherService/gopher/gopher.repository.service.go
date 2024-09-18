package gopher

import (
	"database/sql"
	"fmt"
	"gopherService/customErrors"
)

type RepositoryService interface {
	Create(IncomingGopher) (OutgoingGopher, error)
	Read(int) (OutgoingGopher, error)
}

type repositoryServiceImpl struct {
	db *sql.DB
}

func createGopher(tx *sql.Tx, gopher IncomingGopher) (OutgoingGopher, error) {
	var newGopher OutgoingGopher
	err := tx.QueryRow(
		"INSERT INTO gophers (name, age, color) VALUES ($1, $2, $3) RETURNING id, name, age, color",
		gopher.Name, gopher.Age, gopher.Color,
	).Scan(&newGopher.Id, &newGopher.Name, &newGopher.Age, &newGopher.Color)

	if err != nil {
		return OutgoingGopher{}, &customErrors.DatabaseError{Action: "inserting into gophers table", ErrorString: err.Error()}
	}

	return newGopher, nil
}

func (rs *repositoryServiceImpl) Create(gopher IncomingGopher) (OutgoingGopher, error) {
	tx, err := rs.db.Begin()

	if err != nil {
		return OutgoingGopher{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// if tx.Commit() is successful, this will be a no-op
	defer tx.Rollback()

	outgoingGopher, err := createGopher(tx, gopher)

	if err != nil {
		return OutgoingGopher{}, fmt.Errorf("failed to create gopher: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return OutgoingGopher{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return outgoingGopher, nil
}

func (rs *repositoryServiceImpl) Read(id int) (OutgoingGopher, error) {

	var fetchedGopher OutgoingGopher

	err := rs.db.QueryRow(
		"SELECT id, name, age, color FROM gophers WHERE id = $1",
		id,
	).Scan(&fetchedGopher.Id, &fetchedGopher.Name, &fetchedGopher.Age, &fetchedGopher.Color)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return OutgoingGopher{}, &customErrors.NoRowsError{}
		default:
			return OutgoingGopher{}, err
		}
	}

	return fetchedGopher, nil
}

func NewGopherRepositoryService(db *sql.DB) RepositoryService {
	return &repositoryServiceImpl{db: db}
}
