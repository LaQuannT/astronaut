package postgres

import (
	"context"
	"database/sql"
	"github.com/LaQuannT/astronaut-api/internal/model"
)

type AstronautLogRepository struct {
	db *sql.DB
}

func newAstronautLogRepo(db *sql.DB) *AstronautLogRepository {
	return &AstronautLogRepository{
		db: db,
	}
}

func (r *AstronautLogRepository) CreateAstronautLog(ctx context.Context, a *model.AstronautLog) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO astronaut_log (astronaut_id, space_flights, space_flight_hrs, space_walks, space_walk_hrs,
    status, death_date) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = tx.ExecContext(ctx, stmt, &a.AstronautID, &a.SpaceFlights, &a.SpaceFlightHours, &a.SpaceWalks, &a.SpaceWalkHours,
		&a.Status, newNullString(a.DeathDate))
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AstronautLogRepository) FindAstronautLogById(ctx context.Context, astronautID int) (*model.AstronautLog, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	aLog := new(model.AstronautLog)

	// returns all fields converting death_date to string type turing null values to empty string
	stmt := `SELECT astronaut_id, space_flights, space_flight_hrs, space_walks, space_walk_hrs,
    status, COALESCE(death_date::VARCHAR(255), '') AS death_date FROM astronaut_log WHERE astronaut_id=$1;`
	err = tx.QueryRowContext(ctx, stmt, astronautID).Scan(&aLog.AstronautID, &aLog.SpaceFlights, &aLog.SpaceFlightHours,
		&aLog.SpaceWalks, &aLog.SpaceWalkHours, &aLog.Status, &aLog.DeathDate)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return aLog, nil
}

func (r *AstronautLogRepository) FindAstronautLogs(ctx context.Context) ([]*model.AstronautLog, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var aLogs []*model.AstronautLog

	stmt := `SELECT astronaut_id, space_flights, space_flight_hrs, space_walks, space_walk_hrs, 
       status, COALESCE(death_date::VARCHAR(255), '') AS death_date FROM astronaut_log;`
	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		aLog := new(model.AstronautLog)
		err := rows.Scan(&aLog.AstronautID, &aLog.SpaceFlights, &aLog.SpaceFlightHours,
			&aLog.SpaceWalks, &aLog.SpaceWalkHours, &aLog.Status, &aLog.DeathDate)
		if err != nil {
			return nil, err
		}
		aLogs = append(aLogs, aLog)
	}
	tx.Commit()

	return aLogs, nil
}

func (r *AstronautLogRepository) UpdateAstronautLog(ctx context.Context, a *model.AstronautLog) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE astronaut_log SET space_flights=$1, space_flight_hrs=$2, space_walks=$3, space_walk_hrs=$4,
    status=$5, death_date=$6 WHERE astronaut_id=$7;`

	result, err := tx.ExecContext(ctx, stmt, a.SpaceFlights, a.SpaceFlightHours, a.SpaceWalks, a.SpaceWalkHours,
		a.Status, newNullString(a.DeathDate), a.AstronautID)
	if err != nil {
		return err
	}

	changes, err := result.RowsAffected()
	switch {
	case err != nil:
		return err
	case changes != 1:
		return model.ErrNoChange
	}
	tx.Commit()

	return nil
}

func (r *AstronautLogRepository) DeleteAstronautLog(ctx context.Context, astronautID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM astronaut_log WHERE astronaut_id=$1;`

	result, err := tx.ExecContext(ctx, stmt, astronautID)
	if err != nil {
		return err
	}
	changes, err := result.RowsAffected()
	switch {
	case err != nil:
		return err
	case changes != 1:
		return model.ErrNoChange
	}
	tx.Commit()

	return nil
}
