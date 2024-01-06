package rdbms

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"pm.com/go-countries/domain"
)

type scannable interface {
	Scan(dest ...any) error
}
type prepStatementI interface {
	Prepare(query string) (*sql.Stmt, error)
}
type Database struct {
	db *sql.DB
	languagesDb
	translationDb
	bordersDb
	tldDb
}

func (db *Database) CreateNewCountry(country *domain.Country) (err error) {
	var tx *sql.Tx
	tx, err = db.db.Begin()
	if err != nil {
		return
	}
	wrapErr := func(er error) {
		if err == nil {
			err = fmt.Errorf("%w (country: %d:'%s')", er, country.NumericCode(), country.CommonName())
		} else {
			err = fmt.Errorf("%w; %w (country: %d:'%s')", err, er, country.NumericCode(), country.CommonName())
		}
	}
	defer func(tx *sql.Tx) {
		erR := tx.Rollback()
		if erR != nil {
			wrapErr(erR)
		}
	}(tx)

	countryRecord := newCountryRecord(country)

	continents, subregion, er := db.createRegions(tx, country.Region(), country.Subregion(), country.Continents()...)
	if er != nil {
		wrapErr(er)
		return
	}
	ccDb := countryContinentsDB{prepStmt: tx}
	_, er = ccDb.CreateCountryContinents(countryRecord.CountryId, continents...)
	if er != nil {
		wrapErr(fmt.Errorf("country-continents relations weren't created, %w", er))
		return
	}

	countryRecord.RegionId = subregion.ParentId
	countryRecord.SubregionId = subregion.RegionId

	cdb := countriesDb{prepStmt: tx}
	er = cdb.createCountry(&countryRecord)
	if er != nil {
		wrapErr(er)
		return
	}

	if er = tx.Commit(); er != nil {
		wrapErr(er)
	}
	return
}

func (db *Database) createRegions(prepStmt prepStatementI, region, subregion string, continents ...string) (contIds []uint32, sub RegionRecord, err error) {
	rdb := regionDb{prepStmt: prepStmt}
	var conts []RegionRecord
	conts, err = rdb.readOrCreateContinents(continents...)
	if err != nil {
		return
	}
	var cont RegionRecord
	contIds = make([]uint32, 0, len(conts))
	for _, c := range conts {
		contIds = append(contIds, c.RegionId)
		if c.RegionName == region {
			cont = c
		}
	}
	if len(conts) == 1 {
		cont = conts[0]
	}

	if cont.RegionId == 0 {
		err = fmt.Errorf("continent not found for the region '%s'", region)
		return
	}
	sub, err = rdb.readOrCreateSubregion(cont.RegionName, region, subregion)
	return
}

func (db *Database) Prepare() {

}
func (db *Database) Close() error {
	return db.db.Close()
}
func NewDatabase(db *sql.DB) *Database {
	var result Database
	result.db = db
	return &result
}
func showError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "close error:", err)
	}
}
func closeAndShowError(closable io.Closer) {
	showError(closable.Close())
}
