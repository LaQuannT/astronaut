package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/model"
)

type AstronautRepo struct {
	db *sql.DB
}

func NewAstronautRepo(db *sql.DB) *AstronautRepo {
	return &AstronautRepo{
		db: db,
	}
}

func (r *AstronautRepo) CreateAstronaut(ctx context.Context, a *model.Astronaut) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt := `INSERT INTO astronaut (first_name, last_name, gender, birth_date, birth_place) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	err = tx.QueryRowContext(ctx, stmt, a.FirstName, a.LastName, a.Gender, a.BirthDate, a.BirthPlace).Scan(&a.ID)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (r *AstronautRepo) FindAstronautByID(ctx context.Context, id int) (*model.Astronaut, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed begin tx: %w", err)
	}
	defer tx.Rollback()

	a := new(model.Astronaut)

	stmt := `SELECT * FROM astronaut WHERE id = $1;`
	err = tx.QueryRowContext(ctx, stmt, id).Scan(&a.ID, &a.FirstName, &a.LastName, &a.Gender, &a.BirthPlace)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return a, nil
}

func (r *AstronautRepo) UpdateAstronaut(ctx context.Context, a *model.Astronaut) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt := `UPDATE astronaut SET first_name=$1, last_name=$2, gender=$3, birth_date=$4, birth_place=$5 WHERE id = $6;`

	_, err = tx.ExecContext(ctx, stmt, a.FirstName, a.LastName, a.Gender, a.BirthDate, a.BirthPlace, a.ID)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (r *AstronautRepo) DeleteAstronaut(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt := `DELETE FROM astronaut WHERE id = $1;`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
