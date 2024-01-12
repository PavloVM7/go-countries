package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"io"
	"os"
	"pm.com/go-countries/domain"
)

var (
	ErrDuplicateKey = errors.New("duplicate key")
)

const (
	UniqueViolationCode = "23505"
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
	tldDb
}

func (db *Database) CreateNewCountry(country *domain.Country) (err error) {
	tx, er := db.db.Begin()
	wrapErr := func(er error) {
		err = wrapCountryError(err, er, country.NumericCode(), country.Alpha3Code(), country.CommonName())
	}
	if er != nil {
		wrapErr(er)
		return
	}
	defer func(tx *sql.Tx) {
		erR := tx.Rollback()
		if erR != nil && !errors.Is(erR, sql.ErrTxDone) {
			wrapErr(fmt.Errorf("rollback error: %w", erR))
		}
	}(tx)

	countryRecord := newCountryRecord(country)

	continents, subregion, er := db.createRegions(tx, country.Region(), country.Subregion(), country.Continents()...)
	if er != nil {
		wrapErr(er)
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

	if er = db.createCountryContinents(tx, countryRecord.CountryId, continents...); er != nil {
		wrapErr(er)
		return
	}

	if er = db.createCountryBorders(tx, countryRecord.CountryId, country.Borders()...); er != nil {
		wrapErr(er)
		return
	}
	if er = tx.Commit(); er != nil {
		wrapErr(er)
	}
	return
}

func (db *Database) createCountryBorders(prepStmt prepStatementI, countryId uint16, borders ...string) error {
	bdb := bordersDb{prepStmt: prepStmt}
	if _, err := bdb.createBorders(countryId, borders...); err != nil {
		return fmt.Errorf("country-borders relations weren't created, %w", err)
	}
	return nil
}
func (db *Database) createCountryContinents(prepStmt prepStatementI, countryId uint16, continents ...uint32) error {
	tdb := countryContinentsDB{prepStmt: prepStmt}
	_, err := tdb.createCountryContinents(countryId, continents...)
	if err != nil {
		return fmt.Errorf("country-continents relations weren't created, %w", err)
	}
	return nil
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

func (db *Database) ReadCountry(countryId uint16) (country domain.Country, continents []uint32, regionId, subregionId uint32, err error) {
	wrapErr := func(er error) {
		err = wrapCountryError(err, er, countryId, country.Alpha3Code(), country.CommonName())
	}

	cdb := countriesDb{prepStmt: db.db}
	record, er := cdb.selectCountry(countryId)
	if er != nil {
		wrapErr(er)
		return
	}
	country, regionId, subregionId = newCountry(&record)
	contDb := countryContinentsDB{prepStmt: db.db}
	regRecords, erC := contDb.readCountryContinents(countryId)
	if erC != nil {
		wrapErr(erC)
		return
	}
	continents = make([]uint32, 0, len(regRecords))
	for _, r := range regRecords {
		continents = append(continents, r.ContinentId)
	}
	bdb := bordersDb{prepStmt: db.db}
	borderRecords, erB := bdb.readCountryBorders(countryId)
	if erB != nil {
		wrapErr(erB)
		return
	}
	borders := make([]string, 0, len(borderRecords))
	for _, b := range borderRecords {
		borders = append(borders, b.Alpha3Code)
	}
	country.SetBorders(borders...)
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
func wrapCountryError(baseErr, newErr error, countryId uint16, countryCode, countryName string) error {
	if baseErr == nil {
		return fmt.Errorf("%w (country: %d:%s:'%s')", newErr, countryId, countryCode, countryName)
	} else {
		return fmt.Errorf("%w; %w (country: %d:%s:'%s')", baseErr, newErr, countryId, countryCode, countryName)
	}
}
func showError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "close error:", err)
	}
}
func toPqError(err error) *pq.Error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr
	}
	return nil
}

func closeAndShowError(closable io.Closer) {
	showError(closable.Close())
}
