package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
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

	countryRecord.RegionId = subregion.parentId
	countryRecord.SubregionId = subregion.regionId

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
	if er = db.createCountyCapitals(tx, countryRecord.CountryId, country.CapitalInfo(),
		country.Capital()...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createCountryBorders(tx, countryRecord.CountryId, country.Borders()...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createTopLevelDomains(tx, countryRecord.CountryId, country.TopLevelDomains()...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createAltSpellings(tx, countryRecord.CountryId, country.AltSpellings()...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createTimezones(tx, countryRecord.CountryId, country.Timezones()...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createCurrencies(tx, countryRecord.CountryId, country.Currencies()...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createLanguages(tx, countryRecord.CountryId, country.Languages()...); er != nil {
		wrapErr(er)
		return
	}
	translations := append([]domain.Translation(nil), country.NativeNames()...)
	translations = append(translations, country.Translations()...)
	if er = db.createTranslations(tx, countryRecord.CountryId, translations...); er != nil {
		wrapErr(er)
		return
	}
	if er = db.createDemonyms(tx, countryRecord.CountryId, country.Demonyms()...); er != nil {
		wrapErr(er)
		return
	}
	if er = tx.Commit(); er != nil {
		wrapErr(er)
	}
	return
}
func (db *Database) createDemonyms(prepStmt prepStatementI, countryId uint16, demonyms ...domain.Demonym) error {
	ldb := languagesDb{prepStmt: prepStmt}
	records := make([]*demonymRecord, 0, len(demonyms))
	for _, demonym := range demonyms {
		lang, err := ldb.readOrCrateLanguage(demonym.Language, "")
		if err != nil {
			return fmt.Errorf("demonym language %v wasn't created, %w", demonym, err)
		}
		rec := demonymToRecord(demonym)
		rec.languageId = lang.languageId
		rec.countryId = countryId
		records = append(records, &rec)
	}
	ddb := demonymsDB{prepStmt: prepStmt}
	err := ddb.createDemonyms(records)
	if err != nil {
		return fmt.Errorf("demonyms weren't created, %w", err)
	}
	return nil
}
func (db *Database) createTranslations(prepStmt prepStatementI, countryId uint16, translations ...domain.Translation) error {
	ldb := languagesDb{prepStmt: prepStmt}
	records := make([]*translationRecord, 0, len(translations))
	for _, translation := range translations {
		lang, err := ldb.readOrCrateLanguage(translation.Language, "")
		if err != nil {
			return fmt.Errorf("translation language %v wasn't created, %w", translation, err)
		}
		rec := translationToRecord(translation)
		rec.languageId = lang.languageId
		rec.countryId = countryId
		records = append(records, &rec)
	}
	tdb := translationDb{prepStmt: prepStmt}
	err := tdb.createTranslations(records...)
	if err != nil {
		return fmt.Errorf("translations weren't created, %w", err)
	}
	return nil
}
func (db *Database) createLanguages(prepStmt prepStatementI, countryId uint16, languages ...domain.Language) error {
	ldb := languagesDb{prepStmt: prepStmt}
	ids := make([]uint16, 0, len(languages))
	for _, language := range languages {
		if rec, err := ldb.readOrCrateLanguage(language.Short, language.Name); err != nil {
			return fmt.Errorf("language %v wasn't created, %w", language, err)
		} else {
			ids = append(ids, rec.languageId)
		}
	}
	cldb := countryLanguagesDB{prepStmt: prepStmt}
	if _, err := cldb.createCountryLanguages(countryId, ids...); err != nil {
		return fmt.Errorf("country-languages relations weren't created, %w", err)
	}
	return nil
}
func (db *Database) createCurrencies(prepStmt prepStatementI, countryId uint16, currencies ...domain.Currency) error {
	cdb := currenciesDB{prepStmt: prepStmt}
	ids := make([]uint32, 0, len(currencies))
	for _, currency := range currencies {
		if rec, err := cdb.readOrCreateCurrency(currency.Short, currency.Name, currency.Symbol); err != nil {
			return fmt.Errorf("currency %v wasn't created, %w", currency, err)
		} else {
			ids = append(ids, rec.currencyId)
		}
	}
	ccdb := countryCurrenciesDB{prepStmt: prepStmt}
	if _, err := ccdb.createCountryCurrencies(countryId, ids...); err != nil {
		return fmt.Errorf("country-currencies relations weren't created, %w", err)
	}
	return nil
}
func (db *Database) createTimezones(prepStmt prepStatementI, countryId uint16, timezones ...string) error {
	tzdb := timezonesDB{prepStmt: prepStmt}
	if _, err := tzdb.createTimezones(countryId, timezones...); err != nil {
		return fmt.Errorf("timezones weren't created, %w", err)
	}
	return nil
}
func (db *Database) createAltSpellings(prepStmt prepStatementI, countryId uint16, altSpellings ...string) error {
	asdb := altSpellingDB{prepStmt: prepStmt}
	if _, err := asdb.createAltSpellings(countryId, altSpellings...); err != nil {
		return fmt.Errorf("alt-spellings weren't created, %w", err)
	}
	return nil
}
func (db *Database) createTopLevelDomains(prepStmt prepStatementI, countryId uint16, topLevelDomains ...string) error {
	tldb := tldDb{prepStmt: prepStmt}
	if _, err := tldb.createTopLevelDomains(countryId, topLevelDomains...); err != nil {
		return fmt.Errorf("top-level-domains weren't created, %w", err)
	}
	return nil
}
func (db *Database) createCountryBorders(prepStmt prepStatementI, countryId uint16, borders ...string) error {
	bdb := bordersDb{prepStmt: prepStmt}
	if _, err := bdb.createBorders(countryId, borders...); err != nil {
		return fmt.Errorf("country-borders relations weren't created, %w", err)
	}
	return nil
}
func (db *Database) createCountyCapitals(prepStmt prepStatementI, countryId uint16, capitalsInfo []domain.LatLng, capitals ...string) error {
	cdb := countryCapitalsDB{prepStmt: prepStmt}
	caps, err := cdb.createCapitals(countryId, capitals...)
	if err != nil {
		return fmt.Errorf("country-capitals relations weren't created, %w", err)
	}
	ids := make([]uint32, 0, len(caps))
	for _, c := range caps {
		ids = append(ids, c.id)
	}
	capInfoDb := capitalInfoDb{prepStmt: prepStmt}
	err = capInfoDb.createCapitalsInfo(ids, capitalsInfo)
	if err != nil {
		return fmt.Errorf("country-capitals-info relations weren't created, %w", err)
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
func (db *Database) createRegions(prepStmt prepStatementI, region, subregion string, continents ...string) (contIds []uint32, sub regionRecord, err error) {
	rdb := regionDb{prepStmt: prepStmt}
	var conts []regionRecord
	conts, err = rdb.readOrCreateContinents(continents...)
	if err != nil {
		return
	}
	var cont regionRecord
	contIds = make([]uint32, 0, len(conts))
	for _, c := range conts {
		contIds = append(contIds, c.regionId)
		if c.regionName == region {
			cont = c
		}
	}
	if len(conts) == 1 {
		cont = conts[0]
	}

	if cont.regionId == 0 {
		err = fmt.Errorf("continent not found for the region '%s'", region)
		return
	}
	sub, err = rdb.readOrCreateSubregion(cont.regionName, region, subregion)
	return
}

func (db *Database) ReadCountry(countryId uint16) (country domain.Country, regionId, subregionId uint32, err error) {
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

	if er = db.readCountryContinents(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryRegionSubregion(&country, regionId, subregionId); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryCapitals(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryBorders(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readTopLevelDomains(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryAltSpellings(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryTimezones(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryCurrencies(&country); er != nil {
		wrapErr(er)
		return
	}
	if er = db.readCountryLanguages(&country); er != nil {
		wrapErr(er)
		return
	}
	if err = db.readTranslations(&country); err != nil {
		wrapErr(err)
		return
	}
	return
}
func (db *Database) readTranslations(country *domain.Country) error {
	tdb := translationDb{prepStmt: db.db}
	records, err := tdb.readTranslations(country.NumericCode())
	if err != nil {
		return err
	}
	ids := make([]uint16, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.languageId)
	}
	ldb := languagesDb{prepStmt: db.db}
	langs, err := ldb.readLanguages(ids...)
	if err != nil {
		return err
	}
	lmp := make(map[uint16]string, len(langs))
	for _, l := range langs {
		lmp[l.languageId] = l.language
	}
	for _, record := range records {
		lng := lmp[record.languageId]
		if lng == "" {
			lng = "?"
		}
		if record.native {
			country.AddNativeName(lng, record.commonName, record.officialName)
		} else {
			country.AddTranslation(lng, record.commonName, record.officialName)
		}
	}
	return nil
}
func (db *Database) readCountryLanguages(country *domain.Country) error {
	cdb := countryLanguagesDB{prepStmt: db.db}
	languages, err := cdb.readCountryLanguages(country.NumericCode())
	if err != nil {
		return err
	}
	ids := make([]uint16, 0, len(languages))
	for _, record := range languages {
		ids = append(ids, record.languageId)
	}

	ldb := languagesDb{prepStmt: db.db}
	records, er := ldb.readLanguages(ids...)
	if er != nil {
		return er
	}
	for _, l := range records {
		country.AddLanguage(l.language, l.languageName)
	}
	return nil
}
func (db *Database) readCountryCurrencies(country *domain.Country) error {
	ccdb := countryCurrenciesDB{prepStmt: db.db}
	countryCurrencies, err := ccdb.readCountryCurrencies(country.NumericCode())
	if err != nil {
		return err
	}
	ids := make([]uint32, 0, len(countryCurrencies))
	for _, record := range countryCurrencies {
		ids = append(ids, record.currencyId)
	}
	cdb := currenciesDB{prepStmt: db.db}
	currencies, er := cdb.readCurrencies(ids...)
	if er != nil {
		return er
	}
	for _, c := range currencies {
		country.AddCurrency(c.short, c.name, c.symbol)
	}
	return nil
}
func (db *Database) readCountryTimezones(country *domain.Country) error {
	tdb := timezonesDB{prepStmt: db.db}
	timezones, err := tdb.readTimezones(country.NumericCode())
	if err != nil {
		return err
	}
	tzs := make([]string, 0, len(timezones))
	for _, record := range timezones {
		tzs = append(tzs, record.tz)
	}
	country.SetTimezones(tzs...)
	return nil
}
func (db *Database) readCountryAltSpellings(country *domain.Country) error {
	cdb := altSpellingDB{prepStmt: db.db}
	altSpellings, err := cdb.readAltSpellings(country.NumericCode())
	if err != nil {
		return err
	}
	altS := make([]string, 0, len(altSpellings))
	for _, record := range altSpellings {
		altS = append(altS, record.altSpelling)
	}
	country.SetAltSpellings(altS...)
	return nil
}
func (db *Database) readTopLevelDomains(country *domain.Country) error {
	tdb := tldDb{prepStmt: db.db}
	topLevelDomains, err := tdb.readTopLevelDomains(country.NumericCode())
	if err != nil {
		return err
	}
	tlds := make([]string, 0, len(topLevelDomains))
	for _, tld := range topLevelDomains {
		tlds = append(tlds, tld.tld)
	}
	country.SetTopLevelDomains(tlds...)
	return nil
}
func (db *Database) readCountryCapitals(country *domain.Country) error {
	cdb := countryCapitalsDB{prepStmt: db.db}
	capitalRecords, err := cdb.readCountryCapitals(country.NumericCode())
	if err != nil {
		return err
	}
	capitals := make([]string, 0, len(capitalRecords))
	ids := make([]uint32, 0, len(capitalRecords))
	for _, c := range capitalRecords {
		capitals = append(capitals, c.capital)
		ids = append(ids, c.id)
	}
	country.SetCapital(capitals...)

	capInfoDb := capitalInfoDb{prepStmt: db.db}
	capInfoRecords, err := capInfoDb.readCapitalInfo(ids...)
	if err != nil {
		return err
	}
	capitalsInfo := make([]domain.LatLng, 0, len(capInfoRecords))
	for _, c := range capInfoRecords {
		capitalsInfo = append(capitalsInfo, c.point)
	}
	country.SetCapitalInfo(capitalsInfo...)
	return nil
}
func (db *Database) readCountryBorders(country *domain.Country) error {
	bdb := bordersDb{prepStmt: db.db}
	borderRecords, err := bdb.readCountryBorders(country.NumericCode())
	if err != nil {
		return err
	}
	borders := make([]string, 0, len(borderRecords))
	for _, b := range borderRecords {
		borders = append(borders, b.Alpha3Code)
	}
	country.SetBorders(borders...)
	return nil
}
func (db *Database) readCountryRegionSubregion(country *domain.Country, regionId, subregionId uint32) error {
	regDb := regionDb{prepStmt: db.db}
	records, errC := regDb.readRegionsByIds(regionId, subregionId)
	if errC != nil {
		return errC
	}
	for _, record := range records {
		if record.regionId == regionId {
			country.SetRegion(record.regionName)
		} else if record.regionId == subregionId {
			country.SetSubregion(record.regionName)
		}
	}
	return nil
}
func (db *Database) readCountryContinents(country *domain.Country) error {
	contDb := countryContinentsDB{prepStmt: db.db}
	regRecords, err := contDb.readCountryContinents(country.NumericCode())
	if err != nil {
		return err
	}
	continents := make([]uint32, 0, len(regRecords))
	for _, r := range regRecords {
		continents = append(continents, r.ContinentId)
	}
	regDb := regionDb{prepStmt: db.db}
	contRecs, errC := regDb.readRegionsByIds(continents...)
	if errC != nil {
		return errC
	}
	contStrs := make([]string, 0, len(contRecs))
	for _, rec := range contRecs {
		contStrs = append(contStrs, rec.regionName)
	}
	country.SetContinents(contStrs...)
	return nil
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
