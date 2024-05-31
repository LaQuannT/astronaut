package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/LaQuannT/astronaut-api/internal/database/postgres"
	"github.com/joho/godotenv"
)

var (
	dbConn       *sql.DB
	astroRepo    *postgres.AstronautRepository
	astroLogRepo *postgres.AstronautLogRepository
	academicRepo *postgres.AcademicLogRepository
	militaryRepo *postgres.MilitaryLogRepository
	missionRepo  *postgres.MissionRepository
	userRepo     *postgres.UserRepository
)

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load("../../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .env file: %v", err)
		os.Exit(1)
	}

	connStr := os.Getenv("TEST_DB_URL")
	dbConn, err = postgres.Connect(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	astroRepo, astroLogRepo, academicRepo, militaryRepo, missionRepo, userRepo = postgres.InitializeRepositories(dbConn)

	// ensures tables are built
	err = migration("file://../../migration", connStr, "up")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running up migration: %v\n", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}
