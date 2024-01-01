package rdbms

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(host string, port int, user, password, dbname string) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		if er := db.Close(); er != nil {
			err = fmt.Errorf("db error: %w, close error: %w", err, er)
		}
		return nil, err
	}
	return db, nil
}
