package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type tldDbTestSuite struct {
	databaseBaseTestSuite
	countryRecord CountryRecord
	dtb           tldDb
}

func (s *tldDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = tldDb{prepStmt: s.db}
}
func (s *tldDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
}
func (s *tldDbTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)
}

func Test_tldDbTestSuite(t *testing.T) {
	suite.Run(t, new(tldDbTestSuite))
}

func (s *tldDbTestSuite) TestCreateTopLevelDomains() {
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.dtb.createTopLevelDomains(s.countryRecord.CountryId, tlds...)
	s.Nil(errQ)
	expected := []tldRecord{
		{id: 1, countryId: s.countryRecord.CountryId, tld: ".nl"},
		{id: 2, countryId: s.countryRecord.CountryId, tld: ".nld"},
	}
	s.Equal(expected, actual)
}
func (s *tldDbTestSuite) TestGetTopLevelDomains() {
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.dtb.createTopLevelDomains(s.countryRecord.CountryId, tlds...)
	s.Nil(errQ)
	expected := []tldRecord{
		{id: 1, countryId: s.countryRecord.CountryId, tld: ".nl"},
		{id: 2, countryId: s.countryRecord.CountryId, tld: ".nld"},
	}
	s.Equal(expected, actual)
	actual2, errQ2 := s.dtb.readTopLevelDomains(s.countryRecord.CountryId)
	s.Nil(errQ2)
	s.Equal(actual, actual2)
}
func (s *tldDbTestSuite) TestGetTopLevelDomainsNotExist() {
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.dtb.createTopLevelDomains(s.countryRecord.CountryId, tlds...)
	s.Nil(errQ)
	expected := []tldRecord{
		{id: 1, countryId: s.countryRecord.CountryId, tld: ".nl"},
		{id: 2, countryId: s.countryRecord.CountryId, tld: ".nld"},
	}
	s.Equal(expected, actual)
	actual2, errQ2 := s.dtb.readTopLevelDomains(s.countryRecord.CountryId + 1)
	s.Nil(errQ2)
	s.Equal([]tldRecord{}, actual2)
}
