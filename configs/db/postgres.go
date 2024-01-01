package db

import (
	"errors"
	"os"
)

type PostgresConfig struct {
	DB       string
	User     string
	Password string
}

func (pc *PostgresConfig) Read() error {
	db, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		return errors.New("db is not defined")
	}
	pc.DB = db
	user, uExists := os.LookupEnv("POSTGRES_USER")
	if !uExists {
		return errors.New("user is not defined")
	}
	pc.User = user
	pass, pExists := os.LookupEnv("POSTGRES_PASSWORD")
	if !pExists {
		return errors.New("user password is not defined")
	}
	pc.Password = pass
	return nil
}
