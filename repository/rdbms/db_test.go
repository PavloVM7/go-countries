package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"pm.com/go-countries/configs/db"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
	db       *sql.DB
	database *Database
}

func (s *DatabaseTestSuite) SetupSuite() {
	fmt.Println(">>> From SetupSuite")
	var config db.PostgresConfig
	var err error
	err = config.Read()
	s.T().Log("config error:", err) // ToDo: need to create env variables
	s.db, err = Connect("localhost", 5432, "admin", "admin", "countries_db")
	if err != nil {
		panic(err)
	}
	s.database = NewDatabase(s.db)
	res, er := s.execSqlFile("./sql/references.sql")
	_, _ = fmt.Fprintln(os.Stderr, "create tables:", res, ", error:", er)
}

func (s *DatabaseTestSuite) TearDownSuite() {
	fmt.Println(">>> From TearDownSuite")
	if s.database != nil {
		res, er := s.database.db.Exec("DROP DATABASE countries_db")
		_, _ = fmt.Fprintln(os.Stdout, "drop db result:", res, ", error:", er)
		if err := s.database.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (s *DatabaseTestSuite) TearDownTest() {
	fmt.Println("--- Truncate tables")
	res, err := s.database.db.Exec("TRUNCATE languages RESTART IDENTITY; TRUNCATE regions RESTART IDENTITY;")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	_, _ = fmt.Fprintln(os.Stdout, "truncate result:", res)
}

func (s *DatabaseTestSuite) execSqlFile(file string) (sql.Result, error) {
	bts, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return s.database.db.Exec(string(bts))
}

func (s *DatabaseTestSuite) TestCreateLanguage() {
	lang3 := "eng"
	lang, err := s.database.CreateLanguage(lang3)
	s.Nil(err)
	s.Equal(LanguageRecord{LanguageId: 1, Language: lang3}, lang)
}
func (s *DatabaseTestSuite) TestCreateLanguageDuplicate() {
	lang3 := "eng"
	lang, err := s.database.CreateLanguage(lang3)
	s.Nil(err)
	s.Equal(LanguageRecord{LanguageId: 1, Language: lang3}, lang)
	langD, errD := s.database.CreateLanguage(lang3)
	s.NotNil(errD)
	s.Equal(LanguageRecord{LanguageId: -1, Language: lang3}, langD)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
