package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type countryCapitalsDbTestSuite struct {
	databaseBaseTestSuite
	dtb countryCapitalsDB
}

func (s *countryCapitalsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countryCapitalsDB{prepStmt: s.db}
}
func (s *countryCapitalsDbTestSuite) Test_createCountryCapitals() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, err := s.dtb.createCapitals(country.CountryId, "Amsterdam")
	s.Nil(err)
	expected := []countryCapitalRecord{{id: 1, countryId: country.CountryId, capital: "Amsterdam"}}
	s.Equal(expected, actual)
}
func (s *countryCapitalsDbTestSuite) Test_createCountryCapitalsError() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	_, err = s.dtb.createCapitals(country.CountryId, "Amsterdam", "Amsterdam")
	s.NotNil(err)
}
func (s *countryCapitalsDbTestSuite) Test_readCountryCapitals() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	_, err = s.dtb.createCapitals(country.CountryId, "Amsterdam")
	s.Nil(err)
	actual, err := s.dtb.readCountryCapitals(country.CountryId)
	s.Nil(err)
	expected := []countryCapitalRecord{{id: 1, countryId: country.CountryId, capital: "Amsterdam"}}
	s.Equal(expected, actual)
}
func Test_countryCapitalsDbTestSuite(t *testing.T) {
	suite.Run(t, new(countryCapitalsDbTestSuite))
}
