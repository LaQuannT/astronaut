package postgres

import (
	"context"
	"database/sql"
	"github.com/LaQuannT/astronaut-api/internal/model"
)

type AcademicLogRepository struct {
	db *sql.DB
}

func newAcademicRepo(db *sql.DB) *AcademicLogRepository {
	return &AcademicLogRepository{
		db: db,
	}
}

func (r *AcademicLogRepository) CreateMajor(ctx context.Context, m *model.Major) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO major (course) VALUES ($1) RETURNING id;`

	err = tx.QueryRowContext(ctx, stmt, m.Course).Scan(&m.ID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) CreateAlmaMater(ctx context.Context, a *model.AlmaMater) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO alma_mater (school) VALUES ($1) RETURNING id;`
	err = tx.QueryRowContext(ctx, stmt, a.School).Scan(&a.ID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) AddUnderGradMajor(ctx context.Context, astronautID, majorID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO astronaut_undergrad_major (astronaut_id, major_id) VALUES ($1, $2);`
	_, err = tx.ExecContext(ctx, stmt, astronautID, majorID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) AddGradMajor(ctx context.Context, astronautID, majorID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO astronaut_grad_major (astronaut_id, major_id) VALUES ($1, $2);`
	_, err = tx.ExecContext(ctx, stmt, astronautID, majorID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) AddAstronautAlmaMater(ctx context.Context, astronautID, almaMaterID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO astronaut_alma_mater (astronaut_id, alma_mater_id) VALUES ($1, $2);`
	_, err = tx.ExecContext(ctx, stmt, astronautID, almaMaterID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) UpdateMajor(ctx context.Context, m *model.Major) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE major SET course=$1  WHERE id=$2;`

	_, err = tx.ExecContext(ctx, stmt, m.Course, m.ID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) UpdateAlmaMater(ctx context.Context, a *model.AlmaMater) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE alma_mater SET school=$1 WHERE id=$2;`
	_, err = tx.ExecContext(ctx, stmt, a.School, a.ID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) FindMajorByID(ctx context.Context, id int) (*model.Major, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	m := new(model.Major)

	stmt := `SELECT * FROM major WHERE id=$1;`

	err = tx.QueryRowContext(ctx, stmt, id).Scan(&m.ID, &m.Course)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return m, nil
}

func (r *AcademicLogRepository) FindAlmaMaterByID(ctx context.Context, id int) (*model.AlmaMater, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	m := new(model.AlmaMater)

	stmt := `SELECT * FROM alma_mater WHERE id=$1;`

	err = tx.QueryRowContext(ctx, stmt, id).Scan(&m.ID, &m.School)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return m, nil
}

func (r *AcademicLogRepository) FindAstronautUnderGradMajors(ctx context.Context, astronautID int) ([]*model.Major, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT m.id, m.course FROM astronaut_undergrad_major AS u
	INNER JOIN major AS m ON u.major_id = m.id
	ORDER BY m.course;`

	rows, err := tx.QueryContext(ctx, stmt, astronautID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var majors []*model.Major

	for rows.Next() {
		m := new(model.Major)
		if err := rows.Scan(&m.ID, &m.Course); err != nil {
			return nil, err
		}
		majors = append(majors, m)
	}
	tx.Commit()

	return majors, nil
}

func (r *AcademicLogRepository) FindAstronautGradMajors(ctx context.Context, astronautID int) ([]*model.Major, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT m.id, m.course FROM astronaut_grad_major AS g
	INNER JOIN major AS m ON g.major_id = m.id
	ORDER BY m.course;`

	rows, err := tx.QueryContext(ctx, stmt, astronautID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var majors []*model.Major

	for rows.Next() {
		m := new(model.Major)
		if err := rows.Scan(&m.ID, &m.Course); err != nil {
			return nil, err
		}
		majors = append(majors, m)
	}
	tx.Commit()

	return majors, nil
}

func (r *AcademicLogRepository) FindAstronautAlmaMaters(ctx context.Context, astronautID int) ([]*model.AlmaMater, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT am.id, am.school FROM astronaut_alma_mater AS aa
	INNER JOIN alma_mater AS am ON aa.alma_mater_id = am.id
	ORDER BY am.school;`

	rows, err := tx.QueryContext(ctx, stmt, astronautID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var almaMaters []*model.AlmaMater

	for rows.Next() {
		m := new(model.AlmaMater)
		if err := rows.Scan(&m.ID, &m.School); err != nil {
			return nil, err
		}
		almaMaters = append(almaMaters, m)
	}
	tx.Commit()

	return almaMaters, nil
}

func (r *AcademicLogRepository) DeleteMajor(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM major WHERE id=$1;`

	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM astronaut_undergrad_major WHERE major_id=$1;`

	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM astronaut_grad_major WHERE major_id=$1;`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) DeleteAstronautUnderGradMajor(ctx context.Context, astronautID, majorID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM astronaut_undergrad_major WHERE astronaut_id=$1 AND major_id=$2;`

	_, err = tx.ExecContext(ctx, stmt, astronautID, majorID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) DeleteAstronautGradMajor(ctx context.Context, astronautID, majorID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM astronaut_grad_major WHERE astronaut_id=$1 AND major_id=$2;`

	_, err = tx.ExecContext(ctx, stmt, astronautID, majorID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) DeleteAlmaMater(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM alma_mater WHERE id=$1;`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM astronaut_alma_mater WHERE id=$1;`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *AcademicLogRepository) DeleteAstronautAlmaMater(ctx context.Context, astronautID, majorID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM astronaut_alma_mater WHERE astronaut_id=$1 AND major_id=$2;`
	_, err = tx.ExecContext(ctx, stmt, astronautID, majorID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
