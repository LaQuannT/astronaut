package main

import (
	"os"

	"github.com/LaQuannT/astronaut-api/internal/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	a := new(app.App)

	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	a.Run(os.Getenv("APP_PORT"))
}
