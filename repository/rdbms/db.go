package rdbms

import (
	"database/sql"
	"fmt"
	"os"
)

type scannable interface {
	Scan(dest ...any) error
}
type Database struct {
	db *sql.DB
	regionDb
	translationDb
	bordersDb
	tldDb
	countryContinentsDB
}

func (db *Database) CreateCountry(record *CountryRecord) error {
	stmt, err := db.db.Prepare(`INSERT INTO countries (country_id, alpha2_code, alpha3_code, olympic_code, 
                       fifa_code, flag, population, area, independent, landlocked, un_member, latitude, longitude, 
                       region_id, subregion_id, official_name, common_name, start_of_week, status) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.CountryId, record.Alpha2Code, record.Alpha3Code, record.OlympicCode, record.FifaCode,
		record.Flag, record.Population, record.Area, record.Independent, record.Landlocked, record.UnMember,
		record.Latitude, record.Longitude, record.RegionId, record.SubregionId, record.OfficialName, record.CommonName,
		record.StartOfWeek, record.Status)
	return err
}

func (db *Database) GetCountry(countryId uint16) (CountryRecord, error) {
	var result CountryRecord
	stmt, err := db.db.Prepare("SELECT * FROM countries WHERE country_id=$1")
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	rows, errQ := stmt.Query(countryId)
	if errQ != nil {
		return result, errQ
	}
	defer func(rows *sql.Rows) {
		showError(rows.Close())
	}(rows)
	if rows.Next() {
		return db.rowsToRecord(rows)
	}
	return CountryRecord{}, fmt.Errorf("country not found by id=%d", countryId)
}

func (db *Database) rowsToRecord(scnbl scannable) (CountryRecord, error) {
	var resul CountryRecord
	err := scnbl.Scan(&resul.CountryId, &resul.Alpha2Code, &resul.Alpha3Code, &resul.OlympicCode, &resul.FifaCode,
		&resul.Flag, &resul.Population, &resul.Area, &resul.Independent, &resul.Landlocked, &resul.UnMember,
		&resul.Latitude, &resul.Longitude, &resul.RegionId, &resul.SubregionId, &resul.OfficialName, &resul.CommonName,
		&resul.StartOfWeek, &resul.Status)
	return resul, err
}

func (db *Database) CreateLanguage(language string) (LanguageRecord, error) {
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
func (db *Database) GetLanguage(language string) (LanguageRecord, error) {
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
func (db *Database) Prepare() {

}
func (db *Database) Close() error {
	return db.db.Close()
}
func NewDatabase(db *sql.DB) *Database {
	var result Database
	result.db = db
	result.regionDb.db = db
	result.translationDb.db = db
	result.bordersDb.db = db
	result.tldDb.db = db
	result.countryContinentsDB.db = db
	return &result
}
func showError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "statement close error:", err)
	}
}
