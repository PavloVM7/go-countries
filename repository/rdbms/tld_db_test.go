package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type tldDbTestSuite struct {
	databaseBaseTestSuite
	dtb tldDb
}

func (s *tldDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = tldDb{db: s.db}
}
func Test_tldDbTestSuite(t *testing.T) {
	suite.Run(t, new(tldDbTestSuite))
}

func (s *tldDbTestSuite) TestCreateTopLevelDomains() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.dtb.CreateTopLevelDomains(country.CountryId, tlds...)
	s.Nil(errQ)
	expected := []TldRecord{
		{Id: 1, CountryId: country.CountryId, Tld: ".nl"},
		{Id: 2, CountryId: country.CountryId, Tld: ".nld"},
	}
	s.Equal(expected, actual)
}
func (s *tldDbTestSuite) TestGetTopLevelDomains() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.dtb.CreateTopLevelDomains(country.CountryId, tlds...)
	s.Nil(errQ)
	expected := []TldRecord{
		{Id: 1, CountryId: country.CountryId, Tld: ".nl"},
		{Id: 2, CountryId: country.CountryId, Tld: ".nld"},
	}
	s.Equal(expected, actual)
	actual2, errQ2 := s.dtb.GetTopLevelDomains(country.CountryId)
	s.Nil(errQ2)
	s.Equal(actual, actual2)
}
func (s *tldDbTestSuite) TestGetTopLevelDomainsNotExist() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.dtb.CreateTopLevelDomains(country.CountryId, tlds...)
	s.Nil(errQ)
	expected := []TldRecord{
		{Id: 1, CountryId: country.CountryId, Tld: ".nl"},
		{Id: 2, CountryId: country.CountryId, Tld: ".nld"},
	}
	s.Equal(expected, actual)
	actual2, errQ2 := s.dtb.GetTopLevelDomains(country.CountryId + 1)
	s.Nil(errQ2)
	s.Equal([]TldRecord{}, actual2)
}
