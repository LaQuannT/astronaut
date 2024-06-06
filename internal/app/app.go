package app

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/LaQuannT/astronaut-api/internal/config"
	"github.com/LaQuannT/astronaut-api/internal/database/postgres"
	"github.com/LaQuannT/astronaut-api/internal/transport"
)

func initialize(c *config.Config) http.Handler {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.DBUsername, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)

	dbConn, err := postgres.Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}

	astronautRepository, _, _, _, _, usrRepository := postgres.InitializeRepositories(dbConn)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler := transport.NewServer(
		logger,
		usrRepository,
		astronautRepository,
	)
	return handler
}

func Run() {
	c, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	handler := initialize(c)

	srv := &http.Server{
		Addr:    net.JoinHostPort(c.Host, c.Port),
		Handler: handler,
	}

	log.Printf("Server listening on %q", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
