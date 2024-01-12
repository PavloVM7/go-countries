package rdbms

import (
	"errors"
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
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(lang3, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, lang)
}
func (s *languagesDbTestSuite) TestCreateLanguageDuplicate() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(lang3, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, lang)
	langD, errD := s.dtb.createLanguage(lang3, lang3Name)
	e := toPqError(errD)
	s.T().Log("error:", errD, e.Severity, e.Code, "msg:", e.Message, "hint:", e.Hint, "Detail:", e.Detail)
	s.NotNil(errD)
	if errors.Is(err, ErrDuplicateKey) {
		s.T().Log("error:", errD)
	}
	s.Equal(languageRecord{languageId: 0, language: lang3, languageName: lang3Name}, langD)
}
func (s *languagesDbTestSuite) TestCreateLanguageNameDuplicate() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(lang3, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, lang)
	langD, errD := s.dtb.createLanguage("any", lang3Name)
	s.NotNil(errD)
	s.Equal(languageRecord{languageId: 0, language: "any", languageName: lang3Name}, langD)
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
func (s *languagesDbTestSuite) Test_updateLanguage() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(lang3, "")
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: ""}, lang)
	updRec := languageRecord{languageId: 0, language: lang3, languageName: lang3Name}
	errU := s.dtb.updateLanguageRecord(&updRec)
	s.Nil(errU)
	actual, errS := s.dtb.readeLanguage(lang3)
	s.Nil(errS)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, actual)
}
func (s *languagesDbTestSuite) Test_readAndUpdateLanguage() {
	lang3 := "eng"
	lang3Name := "English"
	lang, err := s.dtb.readOrCrateLanguage(lang3, "")
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: ""}, lang)
	selectStmt, er := s.dtb.prepareSelectLanguage()
	s.Nil(er)
	defer func() {
		s.Nil(selectStmt.Close())
	}()
	actual, errU := s.dtb.readAndUpdateLanguage(selectStmt, lang3, lang3Name)
	s.Nil(errU)
	s.Equal(languageRecord{languageId: 1, language: lang3, languageName: lang3Name}, actual)
}
func (s *languagesDbTestSuite) TestGetLanguage() {
	langEng := "eng"
	lang3Name := "English"
	lang, err := s.dtb.createLanguage(langEng, lang3Name)
	s.Nil(err)
	s.Equal(languageRecord{languageId: 1, language: langEng, languageName: lang3Name}, lang)
	actual, er := s.dtb.readeLanguage(langEng)
	s.Nil(er)
	s.Equal(lang, actual)
}
func (s *languagesDbTestSuite) TestGetLanguageNotExists() {
	lng := "any"
	lang, err := s.dtb.readeLanguage(lng)
	s.NotNil(err)
	s.Equal(languageRecord{languageId: 0, language: lng}, lang)
}
