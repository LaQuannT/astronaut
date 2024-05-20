package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/model"
)

type MissionRepository struct {
	db *sql.DB
}

func newMissionRepo(db *sql.DB) *MissionRepository {
	return &MissionRepository{
		db: db,
	}
}

func (r *MissionRepository) CreateMission(ctx context.Context, m *model.Mission) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO mission (id, name, "alias", date_of_mission, successful) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	err = tx.QueryRowContext(ctx, stmt, m.ID, m.Name, m.Alias, m.Successful).Scan(&m.ID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
func (r *MissionRepository) FindMissionByID(ctx context.Context, id int) (*model.Mission, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	m := new(model.Mission)

	stmt := `SELECT * FROM mission WHERE id = $1;`
	err = tx.QueryRowContext(ctx, stmt, id).Scan(&m.ID, &m.Name, &m.Alias, &m.DateOfMission, &m.Successful)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return m, nil
}

func (r *MissionRepository) FindMissionByNameOrAlias(ctx context.Context, target string) ([]*model.Mission, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT * FROM mission WHERE name ILIKE $1 OR "alias" ILIKE $2 ORDER BY name;`
	target = fmt.Sprintf("%%%s%%", target)

	rows, err := tx.QueryContext(ctx, stmt, target, target)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission

	for rows.Next() {
		m := new(model.Mission)
		if err := rows.Scan(&m.ID, &m.Name, &m.Alias, &m.Successful); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	tx.Commit()

	return missions, nil
}

func (r *MissionRepository) FindAllMissions(ctx context.Context) ([]*model.Mission, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT * FROM mission;`
	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission

	for rows.Next() {
		m := new(model.Mission)
		if err := rows.Scan(&m.ID, &m.Name, &m.Alias, &m.Successful); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	tx.Commit()

	return missions, nil
}

func (r *MissionRepository) UpdateMission(ctx context.Context, m *model.Mission) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE mission SET name=$1, alias=$2, date_of_mission=$3, successful=$4 WHERE id=$5;`

	result, err := tx.ExecContext(ctx, stmt, m.Name, m.Alias, m.DateOfMission, m.Successful, m.ID)
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

func (r *MissionRepository) CreateAstronautMission(ctx context.Context, missionID, astronautID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO astronaut_mission (mission_id, astronaut_id) VALUES ($1, $2);`

	_, err = tx.ExecContext(ctx, stmt, missionID, astronautID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *MissionRepository) FindMissionsByAstronaut(ctx context.Context, astronautID int) ([]*model.Mission, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT id, name, "alias", date_of_mission, successful FROM mission AS M
        INNER JOIN  astronaut_mission AS am ON am.astronaut_id = m.astronaut_id 
        WHERE am.astronaut_id=$1;`

	rows, err := tx.QueryContext(ctx, stmt, astronautID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*model.Mission

	for rows.Next() {
		m := new(model.Mission)
		if err := rows.Scan(&m.ID, &m.Name, &m.Alias, &m.DateOfMission, &m.Successful); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	tx.Commit()

	return missions, nil
}

func (r *MissionRepository) DeleteAstronautMission(ctx context.Context, astronautID, missionID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM astronaut_mission WHERE astronaut_id=$1 AND mission_id=$2;`

	_, err = tx.ExecContext(ctx, stmt, astronautID, missionID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *MissionRepository) DeleteMission(ctx context.Context, missionID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM mission WHERE id=$1;`
	result, err := tx.ExecContext(ctx, stmt, missionID)
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

	stmt = `DELETE FROM astronaut_mission WHERE mission_id=$1;`
	_, err = tx.ExecContext(ctx, stmt, missionID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
