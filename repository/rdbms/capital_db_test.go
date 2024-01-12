package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type countryCapitalsDbTestSuite struct {
	databaseBaseTestSuite
	dtb           countryCapitalsDB
	countryRecord CountryRecord
}

func (s *countryCapitalsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countryCapitalsDB{prepStmt: s.db}
}
func (s *countryCapitalsDbTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)
}
func (s *countryCapitalsDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
}
func (s *countryCapitalsDbTestSuite) Test_createCountryCapitals() {
	actual, err := s.dtb.createCapitals(s.countryRecord.CountryId, "Amsterdam")
	s.Nil(err)
	expected := []countryCapitalRecord{{id: 1, countryId: s.countryRecord.CountryId, capital: "Amsterdam"}}
	s.Equal(expected, actual)
}
func (s *countryCapitalsDbTestSuite) Test_createCountryCapitalsError() {
	_, err := s.dtb.createCapitals(s.countryRecord.CountryId, "Amsterdam", "Amsterdam")
	s.NotNil(err)
}
func (s *countryCapitalsDbTestSuite) Test_readCountryCapitals() {
	_, err := s.dtb.createCapitals(s.countryRecord.CountryId, "Amsterdam")
	s.Nil(err)
	actual, err := s.dtb.readCountryCapitals(s.countryRecord.CountryId)
	s.Nil(err)
	expected := []countryCapitalRecord{{id: 1, countryId: s.countryRecord.CountryId, capital: "Amsterdam"}}
	s.Equal(expected, actual)
}
func Test_countryCapitalsDbTestSuite(t *testing.T) {
	suite.Run(t, new(countryCapitalsDbTestSuite))
}
