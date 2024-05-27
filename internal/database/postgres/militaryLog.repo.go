package postgres

import (
	"context"
	"database/sql"

	"github.com/LaQuannT/astronaut-api/internal/model"
)

type MilitaryLogRepository struct {
	db *sql.DB
}

func newMilitaryLogRepo(db *sql.DB) *MilitaryLogRepository {
	return &MilitaryLogRepository{
		db: db,
	}
}

func (r *MilitaryLogRepository) CreateMilitaryLog(ctx context.Context, m *model.MilitaryLog) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO military_history (astronaut_id, branch, rank, retired) VALUES ($1, $2, $3, $4);`

	_, err = tx.ExecContext(ctx, stmt, m.AstronautID, m.Branch, m.Rank, m.Retired)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *MilitaryLogRepository) FindMilitaryLog(ctx context.Context, astronautID int) (*model.MilitaryLog, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	m := new(model.MilitaryLog)

	stmt := `SELECT * FROM military_history WHERE astronaut_id = $1;`

	err = tx.QueryRowContext(ctx, stmt, astronautID).Scan(&m.AstronautID, &m.Branch, &m.Rank, &m.Retired)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return m, nil
}

func (r *MilitaryLogRepository) FindAllMilitaryLogs(ctx context.Context) ([]*model.MilitaryLog, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var mLogs []*model.MilitaryLog

	stmt := `SELECT * FROM military_history;`

	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		m := new(model.MilitaryLog)
		if err := rows.Scan(&m.AstronautID, &m.Branch, &m.Rank, &m.Retired); err != nil {
			return nil, err
		}
		mLogs = append(mLogs, m)
	}
	tx.Commit()

	return mLogs, nil
}

func (r *MilitaryLogRepository) UpdateMilitaryLog(ctx context.Context, m *model.MilitaryLog) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE military_history SET branch=$1, rank=$2, retired=$3 WHERE astronaut_id=$4`

	result, err := tx.ExecContext(ctx, stmt, m.Branch, m.Rank, m.Retired, m.AstronautID)
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

func (r *MilitaryLogRepository) DeleteMilitaryLog(ctx context.Context, astronautID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM military_history WHERE astronaut_id=$1`
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
