package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"pm.com/go-countries/domain"
	"testing"
)

type countryCapitalInfoDbTestSuite struct {
	databaseBaseTestSuite
	dtb           capitalInfoDb
	countryRecord CountryRecord
	capitals      []countryCapitalRecord
}

func (s *countryCapitalInfoDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = capitalInfoDb{prepStmt: s.db}
}
func (s *countryCapitalInfoDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
}
func (s *countryCapitalInfoDbTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	tmp := countryCapitalsDB{prepStmt: s.db}
	s.capitals, err = tmp.createCapitals(s.countryRecord.CountryId, "Amsterdam")
	s.Nil(err)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)

}

func (s *countryCapitalInfoDbTestSuite) TestCreateCapitalInfo() {
	points := []domain.LatLng{{Lat: 52.35, Lng: 4.92}}
	err := s.dtb.createCapitalsInfo([]uint32{s.capitals[0].id}, points)
	s.Nil(err)
}
func (s *countryCapitalInfoDbTestSuite) Test_readCapitalInfo() {
	points := []domain.LatLng{{Lat: 52.35, Lng: 4.92}}
	err := s.dtb.createCapitalsInfo([]uint32{s.capitals[0].id}, points)
	s.Nil(err)
	actual, err := s.dtb.readCapitalInfo(s.capitals[0].id)
	s.Nil(err)
	s.Equal([]capitalInfoRecord{{capitalId: s.capitals[0].id, point: domain.LatLng{Lat: 52.35, Lng: 4.92}}}, actual)
}
func Test_capitalInfoDb(t *testing.T) {
	suite.Run(t, new(countryCapitalInfoDbTestSuite))
}
