package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func InitializeRepositories(db *sql.DB) (
	*AstronautRepository,
	*AstronautLogRepository,
	*AcademicLogRepository,
	*MilitaryLogRepository,
	*MissionRepository,
	*UserRepository,
) {
	return newAstronautRepo(db),
		newAstronautLogRepo(db),
		newAcademicRepo(db),
		newMilitaryLogRepo(db),
		newMissionRepo(db),
		newUserRepo(db)
}
