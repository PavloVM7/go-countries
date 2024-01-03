package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"strings"
)

type tldDb struct {
	db *sql.DB
}

func (db *tldDb) CreateTopLevelDomains(countryId uint16, tlDomains ...string) ([]TldRecord, error) {
	stmt, err := db.db.Prepare("INSERT INTO top_level_domains (country_id, tld)VALUES ($1,$2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	result := make([]TldRecord, 0, len(tlDomains))
	for _, tld := range tlDomains {
		tld = strings.TrimSpace(tld)
		if len(tld) == 0 {
			continue
		}
		var id uint32
		if er := stmt.QueryRow(countryId, tld).Scan(&id); er == nil {
			result = append(result, TldRecord{Id: id, CountryId: countryId, Tld: tld})
		}
	}
	return result, nil
}
func (db *tldDb) GetTopLevelDomains(countryId uint16) ([]TldRecord, error) {
	stmt, err := db.db.Prepare("SELECT * FROM top_level_domains WHERE country_id=$1")
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
		showError(stmt.Close())
	}(rows)
	result := lists.NewLinkedList[TldRecord]()
	for rows.Next() {
		rec, er := db.toTopLevelDomain(rows)
		if er != nil {
			return nil, fmt.Errorf("%w (country: %d)", er, countryId)
		}
		result.AddLast(rec)
	}
	return result.ToArray(), nil
}
func (db *tldDb) toTopLevelDomain(scn scannable) (TldRecord, error) {
	var res TldRecord
	err := scn.Scan(&res.Id, &res.CountryId, &res.Tld)
	return res, err
}
