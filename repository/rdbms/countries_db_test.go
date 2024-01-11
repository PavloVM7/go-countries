package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type countriesDbTestSuite struct {
	databaseBaseTestSuite
	dtb countriesDb
}

func (s *countriesDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countriesDb{prepStmt: s.db}
}

func (s *countriesDbTestSuite) TestGetCountryNotFound() {
	country := createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, errGet1 := s.dtb.selectCountry(country.CountryId)
	s.Nil(errGet1)
	s.Equal(country, actual)
	idNotExists := country.CountryId + 1
	_, errG2 := s.dtb.selectCountry(idNotExists)
	s.NotNil(errG2)
	s.Equal(fmt.Sprintf("country not found by id=%d", idNotExists), errG2.Error())
}

func (s *countriesDbTestSuite) Test_createCountry() {
	name := "Europe"
	subregionName := "Western Europe"
	_, _, _, err := s.createRegions(name, name, subregionName)
	s.Nil(err)
	record := createTestCountryRecord()
	err = s.dtb.createCountry(&record)
	s.Nil(err)
	actual, errG := s.dtb.selectCountry(record.CountryId)
	s.Nil(errG)
	s.Equal(record, actual)
}

func Test_countriesDbTestSuite(t *testing.T) {
	suite.Run(t, new(countriesDbTestSuite))
}
