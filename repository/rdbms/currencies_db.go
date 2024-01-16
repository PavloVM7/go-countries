package rdbms

import (
	"database/sql"
	"errors"
)

type currencyRecord struct {
	currencyId uint32
	short      string
	name       string
	symbol     string
}

func toCurrencyRecord(scn scannable, result *currencyRecord) error {
	return scn.Scan(&result.currencyId, &result.short, &result.name, &result.symbol)
}

type currenciesDB struct {
	prepStmt prepStatementI
}

func (db *currenciesDB) readOrCreateCurrency(short, name, symbol string) (currencyRecord, error) {
	result, err := db.getCurrency(short)
	if err == nil {
		return result, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return result, err
	}
	return db.createCurrency(short, name, symbol)
}
func (db *currenciesDB) createCurrency(short, name, symbol string) (currencyRecord, error) {
	result := currencyRecord{short: short, name: name, symbol: symbol}
	stmt, err := db.prepStmt.Prepare("INSERT INTO currencies (currency, currency_name, symbol) VALUES ($1, $2, $3) RETURNING currency_id")
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	err = stmt.QueryRow(short, name, symbol).Scan(&result.currencyId)
	return result, err
}

func (db *currenciesDB) getCurrency(short string) (currencyRecord, error) {
	result := currencyRecord{short: short}
	stmt, err := db.prepStmt.Prepare("SELECT * FROM currencies WHERE currency=$1")
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	err = toCurrencyRecord(stmt.QueryRow(short), &result)
	return result, nil
}
