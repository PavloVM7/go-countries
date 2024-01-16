package rdbms

import "github.com/PavloVM7/go-collections/pkg/collections/lists"

type countryLanguageRecord struct {
	id         uint32
	countryId  uint16
	languageId uint16
}

func toCountryLanguageRecord(scn scannable, record *countryLanguageRecord) error {
	return scn.Scan(&record.id, &record.countryId, &record.languageId)
}

type countryLanguagesDB struct {
	prepStmt prepStatementI
}

func (db *countryLanguagesDB) createCountryLanguages(countryId uint16, languages ...uint16) ([]countryLanguageRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_languages (country_id, language_id) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)

	records := make([]countryLanguageRecord, 0, len(languages))
	for _, languageId := range languages {
		record := countryLanguageRecord{countryId: countryId, languageId: languageId}
		if err = stmt.QueryRow(countryId, languageId).Scan(&record.id); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
func (db *countryLanguagesDB) readCountryLanguages(countryId uint16) ([]countryLanguageRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT * FROM country_languages WHERE country_id = $1")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, err := stmt.Query(countryId)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[countryLanguageRecord]()
	for rows.Next() {
		record := countryLanguageRecord{}
		if err = toCountryLanguageRecord(rows, &record); err != nil {
			return nil, err
		}
		result.AddLast(record)
	}
	return result.ToArray(), nil
}
