package rdbms

import (
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"strings"
)

type altSpellingRecord struct {
	id          uint32
	countryId   uint16
	altSpelling string
}

func toAltSpellingRecord(scn scannable, record *altSpellingRecord) error {
	err := scn.Scan(&record.id, &record.countryId, &record.altSpelling)
	return err
}

type altSpellingDB struct {
	prepStmt prepStatementI
}

func (db *altSpellingDB) createAltSpellings(countryId uint16, altSpellings ...string) ([]altSpellingRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_alt_spellings (country_id, spelling)VALUES ($1,$2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)

	records := make([]altSpellingRecord, 0, len(altSpellings))
	for _, altSpelling := range altSpellings {
		altSpelling = strings.TrimSpace(altSpelling)
		if altSpelling == "" {
			continue
		}
		var id uint32
		if err = stmt.QueryRow(countryId, altSpelling).Scan(&id); err != nil {
			return nil, err
		}
		records = append(records, altSpellingRecord{id: id, countryId: countryId, altSpelling: altSpelling})
	}
	return records, nil
}

func (db *altSpellingDB) readAltSpellings(countryId uint16) ([]altSpellingRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT * FROM country_alt_spellings WHERE country_id = $1")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, er := stmt.Query(countryId)
	if er != nil {
		return nil, er
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[altSpellingRecord]()
	for rows.Next() {
		var record altSpellingRecord
		if err = toAltSpellingRecord(rows, &record); err != nil {
			return nil, err
		}
		result.AddLast(record)
	}
	return result.ToArray(), nil
}
