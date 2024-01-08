package rdbms

import "database/sql"

type LanguageRecord struct {
	LanguageId   uint16
	Language     string
	LanguageName string
}

func toLanguageRecord(scn scannable, result *LanguageRecord) error {
	return scn.Scan(&result.LanguageId, &result.Language, &result.LanguageName)
}

type languagesDb struct {
	db *sql.DB
}

func (db *languagesDb) CreateLanguage(language, languageName string) (LanguageRecord, error) {
	result := LanguageRecord{LanguageId: 0, Language: language, LanguageName: languageName}
	stmt, err := db.db.Prepare("INSERT INTO languages (language, language_name) VALUES ($1, $2) RETURNING language_id")
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	err = stmt.QueryRow(result.Language, result.LanguageName).Scan(&result.LanguageId)
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
	err = toLanguageRecord(stmt.QueryRow(result.Language), &result)
	return result, err
}
