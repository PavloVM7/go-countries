package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type countryContinentDbTestSuite struct {
	databaseBaseTestSuite
	dtb countryContinentsDB
}

func (s *countryContinentDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countryContinentsDB{db: s.db}
}

func (s *countryContinentDbTestSuite) TestCreateCountryContinent() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, errC := s.dtb.CreateCountryContinent(country.CountryId, 1)
	s.Nil(errC)
	s.Equal([]CountryContinentRecord{{CountryId: country.CountryId, ContinentId: 1}}, actual)
}
func (s *countryContinentDbTestSuite) TestCreateCountryContinentError() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, errC := s.dtb.CreateCountryContinent(country.CountryId, 1, 200)
	s.NotNil(errC)
	s.Nil(actual)
}
func (s *countryContinentDbTestSuite) TestGetContinents() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	expected, errC := s.dtb.CreateCountryContinent(country.CountryId, 1)
	s.Nil(errC)
	s.Equal([]CountryContinentRecord{{CountryId: country.CountryId, ContinentId: 1}}, expected)
	actual, errG := s.dtb.GetContinents(country.CountryId)
	s.Nil(errG)
	s.Equal(expected, actual)
}
func Test_countryContinentDbTestSuite(t *testing.T) {
	suite.Run(t, new(countryContinentDbTestSuite))
}
