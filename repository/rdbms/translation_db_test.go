package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type translationsDbTestSuite struct {
	databaseBaseTestSuite
	country   CountryRecord
	languages []languageRecord
	dtb       translationDb
}

func (s *translationsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = translationDb{prepStmt: s.db}
}

func (s *translationsDbTestSuite) SetupTest() {
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
	rec, er = cdb.readOrCrateLanguage("jpn", "Japanese")
	s.NoError(er)
	s.languages = append(s.languages, rec)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.country.CountryId, s.country.Alpha2Code,
		s.country.Alpha3Code, s.country.CommonName)
}
func (s *translationsDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.country = CountryRecord{}
	s.languages = nil
}

func Test_translationsDbTestSuite(t *testing.T) {
	suite.Run(t, new(translationsDbTestSuite))
}

func (s *translationsDbTestSuite) TestCreateTranslation() {
	official := "Koninkrijk der Nederlanden"
	common := "Nederland"
	actual, errT := s.dtb.CreateTransaction(s.country.CountryId, s.languages[1].languageId, true, official, common)
	s.Nil(errT)
	s.Equal(translationRecord{id: 1, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: true,
		officialName: official, commonName: common}, actual)
}
func (s *translationsDbTestSuite) TestCreateTranslationNotNativeName() {
	official := "Koninkrijk der Nederlanden"
	common := "Nederland"
	actual, errT := s.dtb.CreateTransaction(s.country.CountryId, s.languages[1].languageId, true, official, common)
	s.Nil(errT)
	s.Equal(translationRecord{id: 1, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: true,
		officialName: official, commonName: common}, actual)
	actual2, errT2 := s.dtb.CreateTransaction(s.country.CountryId, s.languages[1].languageId, false, official, common)
	s.Nil(errT2)
	s.Equal(translationRecord{id: 2, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: false,
		officialName: official, commonName: common}, actual2)
}
func (s *translationsDbTestSuite) Test_crateTranslations() {
	translations := []*translationRecord{
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: true, officialName: "Koninkrijk der Nederlanden", commonName: "Nederland"},
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: false, officialName: "Nederland", commonName: "Nederland"},
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[2].languageId, native: false, officialName: "オランダ", commonName: "オランダ"},
	}

	errT := s.dtb.createTranslations(translations...)
	s.Nil(errT)
	expected := []*translationRecord{
		{id: 1, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: true, officialName: "Koninkrijk der Nederlanden", commonName: "Nederland"},
		{id: 2, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: false, officialName: "Nederland", commonName: "Nederland"},
		{id: 3, countryId: s.country.CountryId, languageId: s.languages[2].languageId, native: false, officialName: "オランダ", commonName: "オランダ"},
	}
	s.Equal(expected, translations)
}
func (s *translationsDbTestSuite) TestReadTranslation() {
	translations := []*translationRecord{
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: true, officialName: "Koninkrijk der Nederlanden", commonName: "Nederland"},
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: false, officialName: "Nederland", commonName: "Nederland"},
		{id: 0, countryId: s.country.CountryId, languageId: s.languages[2].languageId, native: false, officialName: "オランダ", commonName: "オランダ"},
	}

	err := s.dtb.createTranslations(translations...)
	s.Nil(err)
	records, errR := s.dtb.readTranslations(s.country.CountryId)
	s.NoError(errR)
	expected := []translationRecord{
		{id: 1, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: true, officialName: "Koninkrijk der Nederlanden", commonName: "Nederland"},
		{id: 2, countryId: s.country.CountryId, languageId: s.languages[1].languageId, native: false, officialName: "Nederland", commonName: "Nederland"},
		{id: 3, countryId: s.country.CountryId, languageId: s.languages[2].languageId, native: false, officialName: "オランダ", commonName: "オランダ"},
	}
	s.Equal(expected, records)
}
