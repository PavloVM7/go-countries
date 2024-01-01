package rdbms

import (
	"database/sql"
	"fmt"
	"os"
)

type Database struct {
	db *sql.DB
}

func (db *Database) CreateLanguage(language string) (LanguageRecord, error) {
	result := LanguageRecord{LanguageId: -1, Language: language}
	stmt, err := db.db.Prepare("INSERT INTO languages (language) VALUES ($1) RETURNING language_id")
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "statement close error:", err)
		}
	}(stmt)
	err1 := stmt.QueryRow(result.Language).Scan(&result.LanguageId)
	if err1 != nil {
		return result, err1
	}
	return result, nil
}

func (db *Database) Prepare() {

}
func (db *Database) Close() error {
	return db.db.Close()
}
func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}
