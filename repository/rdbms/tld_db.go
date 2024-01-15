package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"strings"
)

type tldRecord struct {
	id        uint32
	countryId uint16
	tld       string
}

type tldDb struct {
	prepStmt prepStatementI
}

func (db *tldDb) createTopLevelDomains(countryId uint16, tlDomains ...string) ([]tldRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO top_level_domains (country_id, tld)VALUES ($1,$2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		showError(stmt.Close())
	}(stmt)
	result := make([]tldRecord, 0, len(tlDomains))
	for _, tld := range tlDomains {
		tld = strings.TrimSpace(tld)
		if len(tld) == 0 {
			continue
		}
		var id uint32
		if er := stmt.QueryRow(countryId, tld).Scan(&id); er == nil {
			result = append(result, tldRecord{id: id, countryId: countryId, tld: tld})
		}
	}
	return result, nil
}
func (db *tldDb) readTopLevelDomains(countryId uint16) ([]tldRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT * FROM top_level_domains WHERE country_id=$1")
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
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[tldRecord]()
	for rows.Next() {
		rec, er := db.toTopLevelDomain(rows)
		if er != nil {
			return nil, fmt.Errorf("%w (country: %d)", er, countryId)
		}
		result.AddLast(rec)
	}
	return result.ToArray(), nil
}
func (db *tldDb) toTopLevelDomain(scn scannable) (tldRecord, error) {
	var res tldRecord
	err := scn.Scan(&res.id, &res.countryId, &res.tld)
	return res, err
}
