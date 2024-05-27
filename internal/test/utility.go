package test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

func migration(filepath, connStr, direction string) error {
	m, err := migrate.New(
		"file://../../migration",
		connStr)
	if err != nil {
		return err
	}

	switch direction {
	case "up":
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	case "down":
		err = m.Down()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}
	return nil
}

func clearTables(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `DELETE FROM military_history;`
	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM astronaut_mission;`

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM astronaut_log;`

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM astronaut;
	ALTER SEQUENCE astronaut_id_seq RESTART WITH 1;`

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM mission;
	ALTER SEQUENCE 	mission_id_seq RESTART WITH 1;`

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
