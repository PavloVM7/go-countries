package rdbms

import (
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
	"pm.com/go-countries/domain"
)

type demonymRecord struct {
	id         uint32
	countryId  uint16
	languageId uint16
	female     string
	male       string
}

func toDemonymRecord(scn scannable, result *demonymRecord) error {
	return scn.Scan(&result.id, &result.countryId, &result.languageId, &result.female, &result.male)
}
func demonymToRecord(demonym domain.Demonym) demonymRecord {
	return demonymRecord{
		id:         0,
		countryId:  0,
		languageId: 0,
		female:     demonym.F,
		male:       demonym.M,
	}
}

type demonymsDB struct {
	prepStmt prepStatementI
}

func (db *demonymsDB) readDemonyms(countryId uint16) ([]demonymRecord, error) {
	stmt, err := db.prepStmt.Prepare(`SELECT id, country_id, language_id, female, male FROM country_demonyms WHERE country_id = $1`)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, err := stmt.Query(countryId)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[demonymRecord]()
	for rows.Next() {
		var record demonymRecord
		if err = toDemonymRecord(rows, &record); err != nil {
			return nil, err
		}
		result.AddLast(record)
	}
	return result.ToArray(), nil
}
func (db *demonymsDB) createDemonyms(records []*demonymRecord) error {
	stmt, err := db.prepStmt.Prepare(`INSERT INTO country_demonyms (country_id, language_id, female, male) 
VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		return err
	}
	defer closeWithShowError(stmt)
	for _, record := range records {
		if err = stmt.QueryRow(record.countryId, record.languageId, record.female, record.male).Scan(&record.id); err != nil {
			return err
		}
	}
	return nil
}
