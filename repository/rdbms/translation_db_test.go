package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type translationsDbTestSuite struct {
	databaseBaseTestSuite
	dtb translationDb
}

func (s *translationsDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = translationDb{db: s.db}
}
func Test_translationsDbTestSuite(t *testing.T) {
	suite.Run(t, new(translationsDbTestSuite))
}
func (s *translationsDbTestSuite) TestCreateTranslation() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	countryId := country.CountryId
	language := "nld"
	official := "Koninkrijk der Nederlanden"
	common := "Nederland"
	lng, errL := s.createLanguage(language)
	s.Nil(errL)
	actual, errT := s.dtb.CreateTransaction(countryId, lng.LanguageId, true, official, common)
	s.Nil(errT)
	s.Equal(TranslationRecord{Id: 1, CountryId: countryId, LanguageId: lng.LanguageId, Native: true,
		OfficialName: official, CommonName: common}, actual)
}
func (s *translationsDbTestSuite) TestCreateTranslationNotNativeName() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	countryId := country.CountryId
	language := "nld"
	official := "Koninkrijk der Nederlanden"
	common := "Nederland"
	lng, errL := s.createLanguage(language)
	s.Nil(errL)
	actual, errT := s.dtb.CreateTransaction(countryId, lng.LanguageId, true, official, common)
	s.Nil(errT)
	s.Equal(TranslationRecord{Id: 1, CountryId: countryId, LanguageId: lng.LanguageId, Native: true,
		OfficialName: official, CommonName: common}, actual)
	actual2, errT2 := s.dtb.CreateTransaction(countryId, lng.LanguageId, false, official, common)
	s.Nil(errT2)
	s.Equal(TranslationRecord{Id: 2, CountryId: countryId, LanguageId: lng.LanguageId, Native: false,
		OfficialName: official, CommonName: common}, actual2)
}
func (s *translationsDbTestSuite) createLanguage(language string) (LanguageRecord, error) {
	langDb := languagesDb{db: s.db}
	return langDb.CreateLanguage(language)
}
