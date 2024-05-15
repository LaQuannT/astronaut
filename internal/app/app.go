package app

import (
	"database/sql"
	"fmt"
	"github.com/LaQuannT/astronaut-api/internal/database"
	"log"
	"net/http"
)

type App struct {
	Router *http.ServeMux
	DB     *sql.DB
}

func (a *App) Initialize(username, password, host, port, dbname, sslmode string) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		username, password, host, port, dbname, sslmode)

	var err error
	a.DB, err = database.Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = http.NewServeMux()
}

func (a *App) Run(port string) {
	port = fmt.Sprintf(":%s", port)

	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(port, a.Router)
	if err != nil {
		log.Fatal(err)
	}
}
