package rdbms

import "database/sql"

type TranslationRecord struct {
	Id           uint32
	CountryId    uint16
	LanguageId   uint16
	Native       bool
	OfficialName string
	CommonName   string
}

type translationDb struct {
	db *sql.DB
}

func (db *translationDb) CreateTransaction(countryId, languageId uint16, native bool, official, common string) (TranslationRecord, error) {
	result := TranslationRecord{Id: 0, CountryId: countryId, LanguageId: languageId, Native: native,
		OfficialName: official, CommonName: common}
	stmt, err := db.db.Prepare(`INSERT INTO translations (country_id, language_id, native, official_name, common_name) 
VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		return result, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	err = stmt.QueryRow(result.CountryId, result.LanguageId, result.Native, result.OfficialName, result.CommonName).
		Scan(&result.Id)
	return result, err
}
