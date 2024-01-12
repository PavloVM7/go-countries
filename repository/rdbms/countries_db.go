package rdbms

import (
	"fmt"
)

type countriesDb struct {
	prepStmt prepStatementI
}

const (
	insertCountry = `INSERT INTO countries (country_id, alpha2_code, alpha3_code, olympic_code, 
                       fifa_code, flag, population, area, independent, landlocked, un_member, latitude, longitude, 
                       region_id, subregion_id, official_name, common_name, start_of_week, status) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
)

func (db *countriesDb) createCountry(record *CountryRecord) error {
	stmt, err := db.prepStmt.Prepare(insertCountry)
	if err != nil {
		return err
	}
	defer closeWithShowError(stmt)
	_, err = stmt.Exec(record.CountryId, record.Alpha2Code, record.Alpha3Code, record.OlympicCode, record.FifaCode,
		record.Flag, record.Population, record.Area, record.Independent, record.Landlocked, record.UnMember,
		record.Latitude, record.Longitude, record.RegionId, record.SubregionId, record.OfficialName, record.CommonName,
		record.StartOfWeek, record.Status)
	return err
}
func (db *countriesDb) selectCountry(countryId uint16) (CountryRecord, error) {
	var result CountryRecord
	stmt, err := db.prepStmt.Prepare("SELECT * FROM countries WHERE country_id=$1")
	if err != nil {
		return result, err
	}
	defer closeWithShowError(stmt)
	err = toCountryRecord(stmt.QueryRow(countryId), &result)
	if err != nil {
		return result, fmt.Errorf("country not found by id=%d", countryId)
	}
	return result, nil
}
