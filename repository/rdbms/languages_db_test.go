package rdbms

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type languagesDbTestSuite struct {
	databaseBaseTestSuite
	dtb languagesDb
}

func (s *languagesDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = languagesDb{db: s.db}
}
func Test_languagesDbTestSuite(t *testing.T) {
	suite.Run(t, new(languagesDbTestSuite))
}
func (s *languagesDbTestSuite) TestCreateLanguage() {
	lang3 := "eng"
	lang, err := s.dtb.CreateLanguage(lang3)
	s.Nil(err)
	s.Equal(LanguageRecord{LanguageId: 1, Language: lang3}, lang)
}
func (s *languagesDbTestSuite) TestCreateLanguageDuplicate() {
	lang3 := "eng"
	lang, err := s.dtb.CreateLanguage(lang3)
	s.Nil(err)
	s.Equal(LanguageRecord{LanguageId: 1, Language: lang3}, lang)
	langD, errD := s.dtb.CreateLanguage(lang3)
	s.NotNil(errD)
	s.Equal(LanguageRecord{LanguageId: 0, Language: lang3}, langD)
}

func (s *languagesDbTestSuite) TestGetLanguage() {
	langEng := "eng"
	lang, err := s.dtb.CreateLanguage(langEng)
	s.Nil(err)
	s.Equal(LanguageRecord{LanguageId: 1, Language: langEng}, lang)
	actual, er := s.dtb.GetLanguage(langEng)
	s.Nil(er)
	s.Equal(lang, actual)
}
func (s *languagesDbTestSuite) TestGetLanguageNotExists() {
	lng := "any"
	lang, err := s.dtb.GetLanguage(lng)
	s.NotNil(err)
	s.Equal(LanguageRecord{LanguageId: 0, Language: lng}, lang)
}
