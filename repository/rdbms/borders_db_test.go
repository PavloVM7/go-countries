package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type bordersDbTestSuite struct {
	databaseBaseTestSuite
	dtb bordersDb
}

func (s *bordersDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = bordersDb{db: s.db}
}
func Test_bordersDbTestSuite(t *testing.T) {
	suite.Run(t, new(bordersDbTestSuite))
}
func (s *bordersDbTestSuite) TestCreateBorders() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual, errB := s.dtb.CreteBorders(country.CountryId, borders...)
	s.Nil(errB)
	expected := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected, actual)
}
func (s *bordersDbTestSuite) TestCreateBorders1duplicate() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.dtb.CreteBorders(country.CountryId, borders[0])
	s.Nil(err1)
	expected1 := []BorderRecord{{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]}}
	s.Equal(expected1, actual1)
	actual2, err2 := s.dtb.CreteBorders(country.CountryId, borders...)
	s.NotNil(err2)
	expected2 := []BorderRecord{{Id: 3, CountryId: country.CountryId, Alpha3Code: borders[1]}}
	s.Equal(expected2, actual2)
}
func (s *bordersDbTestSuite) TestCreateBordersDuplicate() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.dtb.CreteBorders(country.CountryId, borders...)
	s.Nil(err1)
	expected1 := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected1, actual1)
	actual2, err2 := s.dtb.CreteBorders(country.CountryId, borders...)
	s.NotNil(err2)
	s.NotNil(actual2)
	s.Equal(0, len(actual2))
}
func (s *bordersDbTestSuite) TestGetBorders() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.dtb.CreteBorders(country.CountryId, borders...)
	s.Nil(err1)
	expected1 := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected1, actual1)
	actual2, err2 := s.dtb.GetBorders(country.CountryId)
	s.Nil(err2)
	s.Equal(expected1, actual2)
}
func (s *bordersDbTestSuite) TestGetBordersNotExist() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.dtb.CreteBorders(country.CountryId, borders...)
	s.Nil(err1)
	expected1 := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected1, actual1)
	actual2, err2 := s.dtb.GetBorders(country.CountryId + 1)
	s.Nil(err2)
	s.Equal([]BorderRecord{}, actual2)
}
