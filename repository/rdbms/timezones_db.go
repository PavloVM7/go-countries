package rdbms

import (
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
)

type timezoneRecord struct {
	id        int32
	countryId uint16
	tz        string
}

func toTimezoneRecord(scn scannable, record *timezoneRecord) error {
	err := scn.Scan(&record.id, &record.countryId, &record.tz)
	return err
}

type timezonesDB struct {
	prepStmt prepStatementI
}

func (db *timezonesDB) createTimezones(countryId uint16, timezones ...string) ([]timezoneRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_timezone (country_id, timezone)VALUES ($1,$2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	result := make([]timezoneRecord, 0, len(timezones))
	for _, tz := range timezones {
		tz := timezoneRecord{countryId: countryId, tz: tz}
		err = stmt.QueryRow(tz.countryId, tz.tz).Scan(&tz.id)
		if err != nil {
			return nil, err
		}
		result = append(result, tz)
	}
	return result, nil
}
func (db *timezonesDB) readTimezones(countryId uint16) ([]timezoneRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT * FROM country_timezone WHERE country_id = $1")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, err := stmt.Query(countryId)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[timezoneRecord]()
	for rows.Next() {
		var record timezoneRecord
		if er := toTimezoneRecord(rows, &record); er != nil {
			return nil, er
		}
		result.AddLast(record)
	}
	return result.ToArray(), nil
}
