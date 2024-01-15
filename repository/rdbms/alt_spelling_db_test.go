package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type altSpellingDBTestSuite struct {
	databaseBaseTestSuite
	countryRecord CountryRecord
	dtb           altSpellingDB
}

func (s *altSpellingDBTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = altSpellingDB{prepStmt: s.db}
}
func (s *altSpellingDBTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
}
func (s *altSpellingDBTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)
}

func (s *altSpellingDBTestSuite) TestCreateAltSpelling() {
	records, err := s.dtb.createAltSpellings(s.countryRecord.CountryId, "NL", "Holland", "Nederland",
		"The Netherlands")
	s.Nil(err)
	expected := []altSpellingRecord{
		{id: 1, countryId: s.countryRecord.CountryId, altSpelling: "NL"},
		{id: 2, countryId: s.countryRecord.CountryId, altSpelling: "Holland"},
		{id: 3, countryId: s.countryRecord.CountryId, altSpelling: "Nederland"},
		{id: 4, countryId: s.countryRecord.CountryId, altSpelling: "The Netherlands"},
	}
	s.Equal(expected, records)
}
func (s *altSpellingDBTestSuite) TestReadAltSpelling() {
	records, err := s.dtb.createAltSpellings(s.countryRecord.CountryId, "NL", "Holland", "Nederland",
		"The Netherlands")
	s.Nil(err)
	expected := []altSpellingRecord{
		{id: 1, countryId: s.countryRecord.CountryId, altSpelling: "NL"},
		{id: 2, countryId: s.countryRecord.CountryId, altSpelling: "Holland"},
		{id: 3, countryId: s.countryRecord.CountryId, altSpelling: "Nederland"},
		{id: 4, countryId: s.countryRecord.CountryId, altSpelling: "The Netherlands"},
	}
	s.Equal(expected, records)
	records2, err2 := s.dtb.readAltSpellings(s.countryRecord.CountryId)
	s.Nil(err2)
	s.Equal(expected, records2)
}
func Test_altSpellingDBTestSuite(t *testing.T) {
	suite.Run(t, new(altSpellingDBTestSuite))
}
