package rdbms

import (
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
)

type countryCapitalRecord struct {
	id        uint32
	countryId uint16
	capital   string
}

func toCountryCapitalRecord(scn scannable, result *countryCapitalRecord) error {
	return scn.Scan(&result.id, &result.countryId, &result.capital)
}

type countryCapitalsDB struct {
	prepStmt prepStatementI
}

func (db *countryCapitalsDB) createCapitals(countryId uint16, capitals ...string) ([]countryCapitalRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_capitals (country_id, capital) VALUES ($1,$2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	result := make([]countryCapitalRecord, 0, len(capitals))
	for _, capital := range capitals {
		var id uint32
		if er := stmt.QueryRow(countryId, capital).Scan(&id); er == nil {
			result = append(result, countryCapitalRecord{countryId: countryId, capital: capital, id: id})
		} else {
			return nil, er
		}
	}
	return result, err
}
func (db *countryCapitalsDB) readCountryCapitals(countryId uint16) ([]countryCapitalRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT id, country_id, capital FROM country_capitals WHERE country_id = $1")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, err := stmt.Query(countryId)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[countryCapitalRecord]()
	for rows.Next() {
		var record countryCapitalRecord
		if er := toCountryCapitalRecord(rows, &record); er != nil {
			return nil, er
		}
		result.AddLast(record)
	}
	return result.ToArray(), err
}
