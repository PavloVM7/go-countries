package rdbms

import (
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"pm.com/go-countries/domain"
)

type translationRecord struct {
	id           uint32
	countryId    uint16
	languageId   uint16
	native       bool
	officialName string
	commonName   string
}

func toTranslationRecord(scn scannable, record *translationRecord) error {
	return scn.Scan(&record.id, &record.countryId, &record.languageId, &record.native, &record.officialName, &record.commonName)
}

func translationToRecord(translation domain.Translation) translationRecord {
	return translationRecord{
		countryId:    0,
		languageId:   0,
		native:       translation.Native,
		officialName: translation.Official,
		commonName:   translation.Common,
	}
}

type translationDb struct {
	prepStmt prepStatementI
}

func (db *translationDb) createTranslations(translations ...*translationRecord) error {
	stmt, err := db.prepStmt.Prepare(`INSERT INTO translations (country_id, language_id, native, official_name, common_name) 
VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		return err
	}
	defer closeWithShowError(stmt)
	for i := 0; i < len(translations); i++ {
		translation := translations[i]
		err = stmt.QueryRow(translation.countryId, translation.languageId, translation.native, translation.officialName, translation.commonName).
			Scan(&translation.id)
	}
	return err
}
func (db *translationDb) readTranslations(countryId uint16) ([]translationRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT id, country_id, language_id, native, official_name, common_name FROM translations WHERE country_id = $1")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, er := stmt.Query(countryId)
	if er != nil {
		return nil, er
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[translationRecord]()
	for rows.Next() {
		var translation translationRecord
		err = toTranslationRecord(rows, &translation)
		if err != nil {
			return nil, err
		}
		result.AddLast(translation)
	}
	return result.ToArray(), err
}
func (db *translationDb) CreateTransaction(countryId, languageId uint16, native bool, official, common string) (translationRecord, error) {
	result := translationRecord{id: 0, countryId: countryId, languageId: languageId, native: native,
		officialName: official, commonName: common}
	stmt, err := db.prepStmt.Prepare(`INSERT INTO translations (country_id, language_id, native, official_name, common_name) 
VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)

	err = stmt.QueryRow(result.countryId, result.languageId, result.native, result.officialName, result.commonName).
		Scan(&result.id)
	return result, err
}
