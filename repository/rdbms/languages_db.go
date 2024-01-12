package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
)

type languageRecord struct {
	languageId   uint16
	language     string
	languageName string
}

func toLanguageRecord(scn scannable, result *languageRecord) error {
	return scn.Scan(&result.languageId, &result.language, &result.languageName)
}

type languagesDb struct {
	db *sql.DB
}

func (db *languagesDb) readOrCrateLanguage(language, languageName string) (languageRecord, error) {
	record := languageRecord{languageId: 0, language: language, languageName: languageName}
	stmtSelect, err := db.prepareSelectLanguage()
	if err != nil {
		return record, wrapLanguageError(err, record)
	}
	defer closeAndShowError(stmtSelect)
	err = db.readAndUpdateLanguageRecord(stmtSelect, &record)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		stmtInsert, er := db.prepareInsertLanguage()
		if er != nil {
			return record, wrapLanguageError(er, record)
		}
		defer closeAndShowError(stmtInsert)
		err = insertRecord(stmtInsert, &record)
		if err != nil {
			if pqErr := toPqError(err); pqErr != nil {
				if pqErr.Code == UniqueViolationCode {
					err = db.readAndUpdateLanguageRecord(stmtSelect, &record)
				}
			}
		}
	}
	return record, err
}
func (db *languagesDb) updateLanguageRecord(record *languageRecord) error {
	stmt, err := db.prepareUpdateLanguage()
	if err != nil {
		return err
	}
	defer closeAndShowError(stmt)
	return updateRecord(stmt, record)
}
func (db *languagesDb) createLanguage(language, languageName string) (languageRecord, error) {
	result := languageRecord{languageId: 0, language: language, languageName: languageName}
	stmt, err := db.prepareInsertLanguage()
	if err != nil {
		return result, err
	}
	defer closeAndShowError(stmt)
	err = insertRecord(stmt, &result)
	return result, err
}
func (db *languagesDb) readeLanguage(language string) (languageRecord, error) {
	result := languageRecord{languageId: 0, language: language}
	stmt, err := db.prepareSelectLanguage()
	if err != nil {
		return result, err
	}
	defer closeAndShowError(stmt)
	err = selectRecord(stmt, &result)
	return result, err
}
func (db *languagesDb) readAndUpdateLanguage(stmSelect *sql.Stmt, language string, languageName string) (languageRecord, error) {
	result := languageRecord{languageId: 0, language: language, languageName: languageName}
	err := db.readAndUpdateLanguageRecord(stmSelect, &result)
	return result, err
}
func (db *languagesDb) readAndUpdateLanguageRecord(stmSelect *sql.Stmt, record *languageRecord) error {
	oldName := record.languageName
	err := selectRecord(stmSelect, record)
	if err == nil {
		if record.languageName == "" && oldName != "" {
			record.languageName = oldName
			err = db.updateLanguageRecord(record)
			if err == nil {
				err = selectRecord(stmSelect, record)
			}
		}
	}
	if err != nil {
		return wrapLanguageError(err, *record)
	}
	return nil
}
func (db *languagesDb) prepareUpdateLanguage() (*sql.Stmt, error) {
	return db.db.Prepare("UPDATE languages SET language_name = $2 WHERE language = $1")
}

func (db *languagesDb) prepareInsertLanguage() (*sql.Stmt, error) {
	return db.db.Prepare("INSERT INTO languages (language, language_name) VALUES ($1, $2) RETURNING language_id")
}
func (db *languagesDb) prepareSelectLanguage() (*sql.Stmt, error) {
	return db.db.Prepare("SELECT * FROM languages WHERE language = $1")
}

func insertRecord(stmt *sql.Stmt, record *languageRecord) error {
	return stmt.QueryRow(record.language, record.languageName).Scan(&record.languageId)
}
func selectRecord(stmt *sql.Stmt, record *languageRecord) error {
	err := toLanguageRecord(stmt.QueryRow(record.language), record)
	return err
}
func updateRecord(stmt *sql.Stmt, record *languageRecord) error {
	_, err := stmt.Exec(record.language, record.languageName)
	return err
}

func wrapLanguageError(err error, record languageRecord) error {
	return fmt.Errorf("%w, language: %d:%s:'%s'", err, record.languageId, record.language, record.languageName)
}
