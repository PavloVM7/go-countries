package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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
	prepStmt prepStatementI
}

func (db *languagesDb) readLanguages(ids ...uint16) ([]languageRecord, error) {
	stmtSelect, err := db.prepStmt.Prepare("SELECT * FROM languages WHERE language_id = ANY ($1)")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmtSelect)
	rows, er := stmtSelect.Query(pq.Array(ids))
	if er != nil {
		return nil, er
	}
	defer closeWithShowError(rows)
	languages := make([]languageRecord, 0, len(ids))
	for rows.Next() {
		language := languageRecord{}
		err = rows.Scan(&language.languageId, &language.language, &language.languageName)
		if err != nil {
			return languages, err
		}
		languages = append(languages, language)
	}
	return languages, nil
}
func (db *languagesDb) readOrCrateLanguage(language, languageName string) (languageRecord, error) {
	record := languageRecord{languageId: 0, language: language, languageName: languageName}
	err := db.readOrCreateLanguageRecord(&record)
	return record, err
}
func (db *languagesDb) readOrCreateLanguageRecords(records []*languageRecord) error {
	for _, record := range records {
		err := db.readOrCreateLanguageRecord(record)
		if err != nil {
			return err
		}
	}
	return nil
}
func (db *languagesDb) readOrCreateLanguageRecord(record *languageRecord) error {
	stmtSelect, err := db.prepareSelectLanguage()
	if err != nil {
		return wrapLanguageError(err, *record)
	}
	defer closeWithShowError(stmtSelect)
	err = db.readAndUpdateLanguageRecord(stmtSelect, record)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		stmtInsert, er := db.prepareInsertLanguage()
		if er != nil {
			return wrapLanguageError(er, *record)
		}
		defer closeWithShowError(stmtInsert)
		err = insertRecord(stmtInsert, record)
		if err != nil && isErrorUniqueViolation(err) {
			err = db.readAndUpdateLanguageRecord(stmtSelect, record)
		}
	}
	return err
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
func insertRecord(stmt *sql.Stmt, record *languageRecord) error {
	return stmt.QueryRow(record.language, record.languageName).Scan(&record.languageId)
}
func (db *languagesDb) updateLanguageRecord(record *languageRecord) error {
	stmt, err := db.prepareUpdateLanguage()
	if err != nil {
		return err
	}
	defer closeWithShowError(stmt)
	return updateRecord(stmt, record)
}
func updateRecord(stmt *sql.Stmt, record *languageRecord) error {
	_, err := stmt.Exec(record.language, record.languageName)
	return err
}

func (db *languagesDb) createLanguage(language, languageName string) (languageRecord, error) {
	result := languageRecord{languageId: 0, language: language, languageName: languageName}
	stmt, err := db.prepareInsertLanguage()
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	err = insertRecord(stmt, &result)
	return result, err
}
func (db *languagesDb) readLanguage(language string) (languageRecord, error) {
	result := languageRecord{languageId: 0, language: language}
	stmt, err := db.prepareSelectLanguage()
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	err = selectRecord(stmt, &result)
	return result, err
}

func (db *languagesDb) prepareUpdateLanguage() (*sql.Stmt, error) {
	return db.prepStmt.Prepare("UPDATE languages SET language_name = $2 WHERE language = $1")
}

func (db *languagesDb) prepareInsertLanguage() (*sql.Stmt, error) {
	return db.prepStmt.Prepare("INSERT INTO languages (language, language_name) VALUES ($1, $2) RETURNING language_id")
}
func (db *languagesDb) prepareSelectLanguage() (*sql.Stmt, error) {
	return db.prepStmt.Prepare("SELECT * FROM languages WHERE language = $1")
}
func selectRecord(stmt *sql.Stmt, record *languageRecord) error {
	err := toLanguageRecord(stmt.QueryRow(record.language), record)
	return err
}

func wrapLanguageError(err error, record languageRecord) error {
	return fmt.Errorf("%w, language: %d:%s:'%s'", err, record.languageId, record.language, record.languageName)
}
