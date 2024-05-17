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
    status, death_mission, death_date) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = tx.ExecContext(ctx, stmt, a.AstronautID, a.SpaceFlights, a.SpaceFlightHours, a.SpaceWalks, a.SpaceWalkHours,
		a.Status, a.DeathMissionID, a.DeathDate)
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

	stmt := `SELECT * FROM astronaut_log WHERE astronaut_id=$1;`
	err = tx.QueryRowContext(ctx, stmt, astronautID).Scan(&aLog.AstronautID, &aLog.SpaceFlights, &aLog.SpaceFlightHours,
		&aLog.SpaceWalks, &aLog.SpaceWalkHours, &aLog.Status, &aLog.DeathMissionID, &aLog.DeathDate)
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

	stmt := `SELECT * FROM astronaut_log;`
	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		aLog := new(model.AstronautLog)
		err := rows.Scan(&aLog.AstronautID, &aLog.SpaceFlights, &aLog.SpaceFlightHours,
			&aLog.SpaceWalks, &aLog.SpaceWalkHours, &aLog.Status, &aLog.DeathMissionID, &aLog.DeathDate)
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
    status=$5, death_mission=$6, death_date=$7 WHERE astronaut_id=$8;`

	_, err = tx.ExecContext(ctx, stmt, a.SpaceFlights, a.SpaceFlightHours, a.SpaceWalks, a.SpaceWalkHours,
		a.Status, a.DeathMissionID, a.DeathDate, a.AstronautID)
	if err != nil {
		return err
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

	_, err = tx.ExecContext(ctx, stmt, astronautID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
