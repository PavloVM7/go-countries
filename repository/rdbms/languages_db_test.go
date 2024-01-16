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
	s.dtb = languagesDb{prepStmt: s.db}
}
func Test_languagesDbTestSuite(t *testing.T) {
	suite.Run(t, new(languagesDbTestSuite))
}
func (s *languagesDbTestSuite) TestCreateLanguage() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(lang3, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, lang)
}
func (s *languagesDbTestSuite) Test_readOrCrateLanguage() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.readOrCrateLanguage(lang3, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, lang)
}
func (s *languagesDbTestSuite) Test_readOrCrateLanguage_duplicate() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.readOrCrateLanguage(lang3, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, lang)
	actual, errD := s.dtb.readOrCrateLanguage(lang3, lang3Name)
	s.Nil(errD)
	s.Equal(lang, actual)
}
func (s *languagesDbTestSuite) Test_readLanguages() {
	rec1, err1 := s.dtb.readOrCrateLanguage("eng", "English")
	s.NoError(err1)
	rec2, err2 := s.dtb.readOrCrateLanguage("nld", "Dutch")
	s.NoError(err2)
	expected := []languageRecord{rec1, rec2}
	actual, err := s.dtb.readLanguages(1, 2)
	s.NoError(err)
	s.Equal(expected, actual)
}
func (s *languagesDbTestSuite) Test_updateLanguage() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(lang3, "")
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: ""}, lang)
	updRec := languageRecord{languageId: 0, language: lang3, languageName: lang3Name}
	errU := s.dtb.updateLanguageRecord(&updRec)
	s.Nil(errU)
	actual, errS := s.dtb.readLanguage(lang3)
	s.Nil(errS)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, actual)
}
func (s *languagesDbTestSuite) TestGetLanguage() {
	langEng := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(langEng, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: langEng, languageName: lang3Name}, lang)
	actual, er := s.dtb.readLanguage(langEng)
	s.Nil(er)
	s.Equal(lang, actual)
}
func (s *languagesDbTestSuite) TestGetLanguageNotExists() {
	lng := "any"
	lang, err := s.dtb.readLanguage(lng)
	s.NotNil(err)
	s.Equal(languageRecord{languageId: 0, language: lng}, lang)
}
