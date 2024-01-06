package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"pm.com/go-countries/configs/db"
	"testing"
)

type databaseBaseTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *databaseBaseTestSuite) SetupSuite() {
	fmt.Println(">>> From Base SetupSuite")
	var config db.PostgresConfig
	var err error
	err = config.Read()
	s.T().Log("config error:", err) // ToDo: need to create env variables
	s.db, err = Connect("localhost", 5432, "admin", "admin", "countries_db")
	if err != nil {
		panic(err)
	}

	_, err = s.execSqlFile("./sql/references.sql")
	if err != nil {
		panic(err)
	}
	_, err = s.execSqlFile("./sql/countries.sql")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database created")
}

func (s *databaseBaseTestSuite) TearDownSuite() {
	fmt.Println(">>> From Base TearDownSuite")
	if s.db != nil {
		res, er := s.db.Exec("DROP DATABASE countries_db")
		_, _ = fmt.Fprintln(os.Stdout, "drop db result:", res, ", error:", er)
		if err := s.db.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println("db closed")
		}
	}
}

func (s *databaseBaseTestSuite) TearDownTest() {
	fmt.Println("--- Truncate tables")
	res, err := s.db.Exec(`TRUNCATE translations, country_continents, countries, languages, regions,
    borders, top_level_domains
    RESTART IDENTITY CASCADE;`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables truncated.")
	_, _ = fmt.Fprintln(os.Stdout, "truncate result:", res)
}

func (s *databaseBaseTestSuite) execSqlFile(file string) (sql.Result, error) {
	bts, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return s.db.Exec(string(bts))
}

func (s *databaseBaseTestSuite) createRegions(continentName, regionName, subregionName string) (continent, region, subregion RegionRecord, err error) {
	dbRegions := regionDb{prepStmt: s.db}
	continent, err = dbRegions.CreateContinent(continentName)
	if err != nil {
		err = fmt.Errorf("create continent error: %w", err)
		return
	}
	region, err = dbRegions.CreateRegion(regionName, continent.RegionId)
	if err != nil {
		err = fmt.Errorf("create region error: %w", err)
		return
	}
	subregion, err = dbRegions.CreateRegion(subregionName, region.RegionId)
	if err != nil {
		err = fmt.Errorf("create soubregion error: %w", err)
	}
	return
}
func (s *databaseBaseTestSuite) createCountry(continent, region, subregion string, country CountryRecord) error {
	_, r, sr, err := s.createRegions(continent, region, subregion)
	country.RegionId = r.RegionId
	country.SubregionId = sr.RegionId
	if err != nil {
		return err
	}
	datb := countriesDb{prepStmt: s.db}
	err = datb.createCountry(&country)
	if err != nil {
		return fmt.Errorf("create country error: %w", err)
	}
	return nil
}

type DatabaseTestSuite struct {
	databaseBaseTestSuite
	database *Database
}

func (s *DatabaseTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.database = NewDatabase(s.db)
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

func (s *DatabaseTestSuite) TestCrateCountry() {
}
func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func createTestCountry() CountryRecord {
	return CountryRecord{CountryId: 528, Alpha2Code: "NL", Alpha3Code: "NLD", OlympicCode: "NED", FifaCode: "NED",
		Flag: "🇳🇱", Population: 16655799, Area: 41850.0, Independent: true, Landlocked: false, UnMember: true,
		Latitude: 52.5, Longitude: 5.75, RegionId: 2, SubregionId: 3,
		OfficialName: "Kingdom of the Netherlands", CommonName: "Netherlands", StartOfWeek: "monday",
		Status: "officially-assigned"}
}
