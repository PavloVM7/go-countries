package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"strings"
)

type bordersDb struct {
	db *sql.DB
}

func (db *bordersDb) GetBorders(countryId uint16) ([]BorderRecord, error) {
	stmt, err := db.db.Prepare("SELECT * FROM borders WHERE country_id=$1")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	rows, errQ := stmt.Query(countryId)
	if errQ != nil {
		return nil, errQ
	}
	defer func(rows *sql.Rows) {
		showError(rows.Close())
	}(rows)
	result := lists.NewLinkedList[BorderRecord]()
	for rows.Next() {
		if record, er := db.toBorderRecord(rows); er == nil {
			result.AddLast(record)
		} else {
			return nil, fmt.Errorf("%w (country: %d)", er, countryId)
		}
	}
	return result.ToArray(), nil
}
func (db *bordersDb) toBorderRecord(scn scannable) (BorderRecord, error) {
	var res BorderRecord
	err := scn.Scan(&res.Id, &res.CountryId, &res.Alpha3Code)
	return res, err
}
func (db *bordersDb) CreteBorders(countryId uint16, borders ...string) ([]BorderRecord, error) {
	stmt, err := db.db.Prepare("INSERT INTO borders (country_id, alpha3_code) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	result := make([]BorderRecord, 0, len(borders))
	for _, border := range borders {
		border = strings.TrimSpace(border)
		if len(border) == 0 {
			continue
		}
		var id uint32
		if er := stmt.QueryRow(countryId, border).Scan(&id); er == nil {
			result = append(result, BorderRecord{Id: id, CountryId: countryId, Alpha3Code: border})
		} else {
			if err != nil {
				err = fmt.Errorf("%w; %w (country: %d, border: '%s')", err, er, countryId, border)
			} else {
				err = fmt.Errorf("%w (country: %d, border: '%s')", er, countryId, border)
			}
		}
	}
	return result, err
}
