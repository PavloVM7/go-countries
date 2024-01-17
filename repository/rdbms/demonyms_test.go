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
	rec, er = cdb.readOrCrateLanguage("jpn", "Japanese")
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
