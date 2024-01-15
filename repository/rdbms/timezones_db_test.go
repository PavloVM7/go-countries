package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type timezonesDBTestSuite struct {
	databaseBaseTestSuite
	countryRecord CountryRecord
	dtb           timezonesDB
}

func (s *timezonesDBTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = timezonesDB{prepStmt: s.db}
}
func (s *timezonesDBTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
}
func (s *timezonesDBTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)
}
func (s *timezonesDBTestSuite) TestCreateTimezone() {
	records, err := s.dtb.createTimezones(s.countryRecord.CountryId, "UTC+01:00", "UTC+02:00")
	s.Nil(err)
	expected := []timezoneRecord{
		{id: 1, countryId: s.countryRecord.CountryId, tz: "UTC+01:00"},
		{id: 2, countryId: s.countryRecord.CountryId, tz: "UTC+02:00"},
	}
	s.Equal(expected, records)
}
func (s *timezonesDBTestSuite) TestReadTimezones() {
	records, err := s.dtb.createTimezones(s.countryRecord.CountryId, "UTC+01:00", "UTC+02:00")
	s.Nil(err)
	expected := []timezoneRecord{
		{id: 1, countryId: s.countryRecord.CountryId, tz: "UTC+01:00"},
		{id: 2, countryId: s.countryRecord.CountryId, tz: "UTC+02:00"},
	}
	s.Equal(expected, records)
	records2, err2 := s.dtb.readTimezones(s.countryRecord.CountryId)
	s.Nil(err2)
	s.Equal(expected, records2)
}
func Test_timezonesDBTestSuite(t *testing.T) {
	suite.Run(t, new(timezonesDBTestSuite))
}
