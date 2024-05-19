package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/model"
)

type AstronautRepository struct {
	db *sql.DB
}

func NewAstronautRepo(db *sql.DB) *AstronautRepository {
	return &AstronautRepository{
		db: db,
	}
}

func (r *AstronautRepository) CreateAstronaut(ctx context.Context, a *model.Astronaut) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
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

func (r *AstronautRepository) FindAstronautByID(ctx context.Context, id int) (*model.Astronaut, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	a := new(model.Astronaut)

	stmt := `SELECT * FROM astronaut WHERE id = $1;`
	err = tx.QueryRowContext(ctx, stmt, id).Scan(&a.ID, &a.FirstName, &a.LastName, &a.Gender, &a.BirthDate, &a.BirthPlace)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return a, nil
}

func (r *AstronautRepository) UpdateAstronaut(ctx context.Context, a *model.Astronaut) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
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

func (r *AstronautRepository) DeleteAstronaut(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
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

func (r *AstronautRepository) FindAstronauts(ctx context.Context) ([]*model.Astronaut, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var astronauts []*model.Astronaut

	stmt := `SELECT * FROM astronaut ORDER BY last_name;`

	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Astronaut
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Gender, &a.BirthDate, &a.BirthPlace)
		if err != nil {
			return nil, err
		}
		astronauts = append(astronauts, &a)
	}
	tx.Commit()

	return astronauts, nil
}

func (r *AstronautRepository) FindAstronautByName(ctx context.Context, name string) ([]*model.Astronaut, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var astronauts []*model.Astronaut

	stmt := `SELECT * FROM astronaut WHERE CONCAT(first_name, ' ', last_name) ILIKE $1 ORDER BY last_name;`
	name = fmt.Sprintf("%%%s%%", name)

	rows, err := tx.QueryContext(ctx, stmt, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Astronaut
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Gender, &a.BirthDate, &a.BirthPlace)
		if err != nil {
			return nil, err
		}
		astronauts = append(astronauts, &a)
	}
	tx.Commit()

	return astronauts, nil
}
