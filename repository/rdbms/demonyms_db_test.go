package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type demonymsDbTestSuite struct {
	databaseBaseTestSuite
	country   CountryRecord
	languages []languageRecord
	dtb       demonymsDB
}

func (s *demonymsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = demonymsDB{prepStmt: s.db}
}

func (s *demonymsDbTestSuite) SetupTest() {
	s.country = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.country)
	s.Nil(err)
	cdb := languagesDb{prepStmt: s.db}
	rec, er := cdb.readOrCrateLanguage("eng", "English")
	s.NoError(er)
	s.languages = append(s.languages, rec)
	rec, er = cdb.readOrCrateLanguage("nld", "Dutch")
	s.NoError(er)
	s.languages = append(s.languages, rec)
	rec, er = cdb.readOrCrateLanguage("fra", "French")
	s.NoError(er)
	s.languages = append(s.languages, rec)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.country.CountryId, s.country.Alpha2Code,
		s.country.Alpha3Code, s.country.CommonName)
}
func (s *demonymsDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.country = CountryRecord{}
	s.languages = nil
}

func Test_demonymsDbTestSuite(t *testing.T) {
	suite.Run(t, new(demonymsDbTestSuite))
}
func (s *demonymsDbTestSuite) Test_createDemonyms() {
	demonyms := []*demonymRecord{
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[0].languageId, female: "Dutch", male: "Dutch"},
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[2].languageId, female: "Néerlandaise", male: "Néerlandais"},
	}
	err := s.dtb.createDemonyms(demonyms)
	s.NoError(err)
	s.EqualValues(1, demonyms[0].id)
	s.EqualValues(2, demonyms[1].id)
}
func (s *demonymsDbTestSuite) Test_readDemonyms() {
	demonyms := []*demonymRecord{
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[0].languageId, female: "Dutch", male: "Dutch"},
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[2].languageId, female: "Néerlandaise", male: "Néerlandais"},
	}
	err := s.dtb.createDemonyms(demonyms)
	s.NoError(err)
	actual, err := s.dtb.readDemonyms(s.country.CountryId)
	s.NoError(err)
	s.Equal([]demonymRecord{*demonyms[0], *demonyms[1]}, actual)
}
