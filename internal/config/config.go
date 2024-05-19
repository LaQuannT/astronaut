package config

import (
	"errors"
	"os"
)

type Config struct {
	DBUsername string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBSSLMode  string
	Port       string
	Host       string
}

func New() (*Config, error) {
	username, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		return nil, errors.New("DB_USERNAME environment variable not set")
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, errors.New("DB_PASSWORD environment variable not set")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, errors.New("DB_NAME environment variable not set")
	}

	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return nil, errors.New("DB_HOST environment variable not set")
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return nil, errors.New("DB_PORT environment variable not set")
	}

	sslMode, ok := os.LookupEnv("DB_SSL_MODE")
	if !ok {
		return nil, errors.New("DB_SSL_MODE environment variable not set")
	}

	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		return nil, errors.New("APP_PORT environment variable not set")
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		return nil, errors.New("APP_HOST environment variable not set")
	}

	return &Config{
		DBUsername: username,
		DBPassword: password,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBSSLMode:  sslMode,
		Port:       port,
		Host:       host,
	}, nil
}
