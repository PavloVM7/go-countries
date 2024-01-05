package rdbms

import "database/sql"

type LanguageRecord struct {
	LanguageId uint16
	Language   string
}

type languagesDb struct {
	db *sql.DB
}

func (db *languagesDb) CreateLanguage(language string) (LanguageRecord, error) {
	result := LanguageRecord{LanguageId: 0, Language: language}
	stmt, err := db.db.Prepare("INSERT INTO languages (language) VALUES ($1) RETURNING language_id")
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	err = stmt.QueryRow(result.Language).Scan(&result.LanguageId)
	return result, err
}
func (db *languagesDb) GetLanguage(language string) (LanguageRecord, error) {
	result := LanguageRecord{LanguageId: 0, Language: language}
	stmt, err := db.db.Prepare("SELECT * FROM languages WHERE language=$1")
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	err = stmt.QueryRow(result.Language).Scan(&result.LanguageId, &result.Language)
	return result, err
}
