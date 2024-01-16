package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type countryLanguagesDbTestSuite struct {
	databaseBaseTestSuite
	countryRecord CountryRecord
	languages     []languageRecord
	dtb           countryLanguagesDB
}

func (s *countryLanguagesDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countryLanguagesDB{prepStmt: s.db}
}
func (s *countryLanguagesDbTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	cdb := languagesDb{prepStmt: s.db}
	rec, er := cdb.readOrCrateLanguage("eng", "English")
	s.NoError(er)
	s.languages = append(s.languages, rec)
	rec, er = cdb.readOrCrateLanguage("spa", "Spanish")
	s.NoError(er)
	s.languages = append(s.languages, rec)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)
}
func (s *countryLanguagesDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
	s.languages = nil
}
func (s *countryLanguagesDbTestSuite) Test_createCountryLanguages() {
	records, err := s.dtb.createCountryLanguages(s.countryRecord.CountryId, s.languages[0].languageId, s.languages[1].languageId)
	s.NoError(err)
	expected := []countryLanguageRecord{
		{id: 1, countryId: s.countryRecord.CountryId, languageId: s.languages[0].languageId},
		{id: 2, countryId: s.countryRecord.CountryId, languageId: s.languages[1].languageId},
	}
	s.Equal(expected, records)
}
func (s *countryLanguagesDbTestSuite) Test_readCountryLanguages() {
	_, err := s.dtb.createCountryLanguages(s.countryRecord.CountryId, s.languages[0].languageId, s.languages[1].languageId)
	s.NoError(err)
	records, err := s.dtb.readCountryLanguages(s.countryRecord.CountryId)
	s.NoError(err)
	expected := []countryLanguageRecord{
		{id: 1, countryId: s.countryRecord.CountryId, languageId: s.languages[0].languageId},
		{id: 2, countryId: s.countryRecord.CountryId, languageId: s.languages[1].languageId},
	}
	s.Equal(expected, records)
}
func Test_countryLanguagesDbTestSuite(t *testing.T) {
	suite.Run(t, new(countryLanguagesDbTestSuite))
}
