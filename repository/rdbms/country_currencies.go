package rdbms

import (
	"fmt"
	"github.com/PavloVM7/go-collections/pkg/collections/lists"
)

type countryCurrenciesRecord struct {
	id         uint32
	countryId  uint16
	currencyId uint32
}

func toCountryCurrenciesRecord(scn scannable, result *countryCurrenciesRecord) error {
	return scn.Scan(&result.id, &result.countryId, &result.currencyId)
}

type countryCurrenciesDB struct {
	prepStmt prepStatementI
}

func (db *countryCurrenciesDB) createCountryCurrencies(countryId uint16, currencies ...uint32) ([]countryCurrenciesRecord, error) {
	stmt, err := db.prepStmt.Prepare("INSERT INTO country_currencies (country_id, currency_id) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	result := make([]countryCurrenciesRecord, 0, len(currencies))
	for _, currency := range currencies {
		res := countryCurrenciesRecord{countryId: countryId, currencyId: currency}
		if err = stmt.QueryRow(countryId, currency).Scan(&res.id); err != nil {
			return nil, fmt.Errorf("currency %v wasn't created, %w", currency, err)
		}
		result = append(result, res)
	}
	return result, nil
}
func (db *countryCurrenciesDB) readCountryCurrencies(countryId uint16) ([]countryCurrenciesRecord, error) {
	stmt, err := db.prepStmt.Prepare("SELECT * FROM country_currencies WHERE country_id=$1")
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(stmt)
	rows, err := stmt.Query(countryId)
	if err != nil {
		return nil, err
	}
	defer closeWithShowError(rows)
	result := lists.NewLinkedList[countryCurrenciesRecord]()
	for rows.Next() {
		res := countryCurrenciesRecord{}
		if err = toCountryCurrenciesRecord(rows, &res); err != nil {
			return nil, err
		}
		result.AddLast(res)
	}
	return result.ToArray(), nil
}
