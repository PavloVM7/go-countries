package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
)

type CountryContinentRecord struct {
	CountryId   uint16
	ContinentId uint32
}

func rowsToCountryContinentRecord(scn scannable) (CountryContinentRecord, error) {
	var res CountryContinentRecord
	err := scn.Scan(&res.CountryId, &res.ContinentId)
	return res, err
}

type countryContinentsDB struct {
	prepStmt prepStatementI
}

func (db *countryContinentsDB) createCountryContinents(countryId uint16, continents ...uint32) ([]CountryContinentRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_continents (country_id, continent_id) VALUES ($1,$2)")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	result := lists.NewLinkedList[CountryContinentRecord]()
	for _, continent := range continents {
		_, err = stmt.Exec(countryId, continent)
		if err != nil {
			return nil, fmt.Errorf("%w (country: %d, continent: %d)", err, countryId, continent)
		}
		result.AddLast(CountryContinentRecord{CountryId: countryId, ContinentId: continent})
	}
	return result.ToArray(), nil
}

func (db *countryContinentsDB) readCountryContinents(countryId uint16) ([]CountryContinentRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT * FROM country_continents WHERE country_id=$1")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	rows, errQ := stmt.Query(countryId)
	if errQ != nil {
		return nil, fmt.Errorf("%w (country:%d)", errQ, countryId)
	}
	result := lists.NewLinkedList[CountryContinentRecord]()
	for rows.Next() {
		record, er := rowsToCountryContinentRecord(rows)
		if er != nil {
			return nil, fmt.Errorf("%w (country:%d)", er, countryId)
		}
		result.AddLast(record)
	}
	return result.ToArray(), nil
}
