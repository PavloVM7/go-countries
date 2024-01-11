package rdbms

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"testing"
)

type countryContinentDbTestSuite struct {
	databaseBaseTestSuite
	dtb countryContinentsDB
}

func (s *countryContinentDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countryContinentsDB{prepStmt: s.db}
}

func (s *countryContinentDbTestSuite) TestCreateCountryContinents_transaction() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)

	tx, err := s.databaseBaseTestSuite.db.Begin()
	s.Nil(err)
	s.NotNil(tx)
	defer func(tx *sql.Tx) {
		showError(tx.Rollback())
	}(tx)

	ccDb := countryContinentsDB{prepStmt: tx}

	actual, errC := ccDb.createCountryContinents(country.CountryId, 1)
	s.Nil(errC)
	err = tx.Commit()
	s.Nil(err)
	s.Equal([]CountryContinentRecord{{CountryId: country.CountryId, ContinentId: 1}}, actual)
}
func (s *countryContinentDbTestSuite) TestCreateCountryContinents() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, errC := s.dtb.createCountryContinents(country.CountryId, 1)
	s.Nil(errC)
	s.Equal([]CountryContinentRecord{{CountryId: country.CountryId, ContinentId: 1}}, actual)
}
func (s *countryContinentDbTestSuite) TestCreateCountryContinentError() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, errC := s.dtb.createCountryContinents(country.CountryId, 1, 200)
	s.NotNil(errC)
	s.Nil(actual)
}
func (s *countryContinentDbTestSuite) TestGetContinents() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	expected, errC := s.dtb.createCountryContinents(country.CountryId, 1)
	s.Nil(errC)
	s.Equal([]CountryContinentRecord{{CountryId: country.CountryId, ContinentId: 1}}, expected)
	actual, errG := s.dtb.readCountryContinents(country.CountryId)
	s.Nil(errG)
	s.Equal(expected, actual)
}
func Test_countryContinentDbTestSuite(t *testing.T) {
	suite.Run(t, new(countryContinentDbTestSuite))
}
