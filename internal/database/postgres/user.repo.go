package postgres

import (
	"context"
	"database/sql"

	"github.com/LaQuannT/astronaut-api/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func newUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *model.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO "user" (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at;`

	err = tx.QueryRowContext(ctx, stmt, u.FirstName, u.LastName, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO api_key (user_id) VALUES ($1) RETURNING key;`

	err = tx.QueryRowContext(ctx, stmt, u.ID).Scan(&u.APIKey)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT u.id, u.first_name, u.last_name, u.email, u.password, a.key, u.created_at, u.updated_at FROM "user" AS U 
         INNER JOIN api_key AS a ON a.user_id = u.id WHERE u.id = $1;`

	u := new(model.User)

	err = tx.QueryRowContext(ctx, stmt, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return u, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT u.id, u.first_name, u.last_name, u.email, u.password, a.key, u.created_at, u.updated_at FROM "user" AS U 
         INNER JOIN api_key AS a ON a.user_id = u.id WHERE u.email = $1;`

	u := new(model.User)

	err = tx.QueryRowContext(ctx, stmt, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return u, nil
}

func (r *UserRepository) FindAllUsers(ctx context.Context) ([]*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var users []*model.User

	stmt := `SELECT id, first_name, last_name, email, created_at, updated_at FROM "user";`
	rows, err := tx.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(model.User)
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	tx.Commit()

	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *model.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE "user" SET first_name=$1, last_name=$2, email=$3 WHERE id = $4;`

	_, err = tx.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, u.ID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *UserRepository) RestUserPassword(ctx context.Context, hash string, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE "user" SET password=$1 WHERE id = $2;`

	result, err := tx.ExecContext(ctx, stmt, hash, id)
	if err != nil {
		return err
	}

	changes, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if changes != 1 {
		return model.ErrNoChange
	}
	tx.Commit()
	return nil
}

func (r *UserRepository) GenerateNewUserAPIKey(ctx context.Context, id int) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var key string

	stmt := `UPDATE api_key SET key=uuid_generate_v4() WHERE user_id = $1 RETURNING key;`
	err = tx.QueryRowContext(ctx, stmt, id).Scan(&key)
	if err != nil {
		return "", err
	}
	tx.Commit()
	return key, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM api_key WHERE user_id = $1;`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM "user" WHERE id = $1;`
	result, err := tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	changes, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if changes != 1 {
		return model.ErrNoChange
	}

	stmt = `DELETE FROM admin WHERE EXISTS(SELECT * FROM admin WHERE id = $1);`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (r *UserRepository) GiveAdminPrivileges(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO admin (user_id) VALUES ($1);`
	_, err = tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (r *UserRepository) RevokeAdminPrivileges(ctx context.Context, id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM admin WHERE user_id = $1;`
	result, err := tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	changes, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if changes != 1 {
		return model.ErrNoChange
	}
	tx.Commit()

	return nil
}

func (r *UserRepository) FindUserByAPIKey(ctx context.Context, key string) (*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := `SELECT  u.id, u.first_name, u.last_name, u.email, u.password, a.key, u.created_at, u.updated_at  FROM api_key AS a 
    INNER JOIN "user" AS u ON u.id = a.user_id
    WHERE a.key = $1;`

	u := new(model.User)

	err = tx.QueryRowContext(ctx, stmt, key).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.APIKey, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return u, nil
}

func (r *UserRepository) IsAdmin(ctx context.Context, userID int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt := `SELECT COUNT(DISTINCT user_id) FROM admin WHERE user_id=$1;`

	var userCount int

	err = tx.QueryRowContext(ctx, stmt, userID).Scan(&userCount)
	if err != nil {
		return 0, err
	}
	tx.Commit()

	return userCount, nil
}
