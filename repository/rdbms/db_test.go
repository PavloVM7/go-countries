package rdbms

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"pm.com/go-countries/domain"
	"pm.com/go-countries/internal/config/db"
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
    borders, country_capitals, country_capital_info, top_level_domains
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

func (s *databaseBaseTestSuite) createRegions(continentName, regionName, subregionName string) (continent, region, subregion regionRecord, err error) {
	dbRegions := regionDb{prepStmt: s.db}
	conts, er := dbRegions.readOrCreateContinents(continentName)
	if er != nil {
		err = fmt.Errorf("create continent error: %w", er)
		return
	}
	continent = conts[0]
	region, err = dbRegions.CreateRegion(regionName, continent.regionId)
	if err != nil {
		err = fmt.Errorf("create region error: %w", err)
		return
	}
	subregion, err = dbRegions.CreateRegion(subregionName, region.regionId)
	if err != nil {
		err = fmt.Errorf("create soubregion error: %w", err)
	}
	return
}
func (s *databaseBaseTestSuite) createCountry(continent, region, subregion string, country CountryRecord) error {
	_, r, sr, err := s.createRegions(continent, region, subregion)
	country.RegionId = r.regionId
	country.SubregionId = sr.regionId
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
	country := createTestCountry()
	err := s.database.CreateNewCountry(&country)
	s.Nil(err)
	cdb := countriesDb{prepStmt: s.db}
	record, errR := cdb.selectCountry(country.NumericCode())
	s.Nil(errR)
	expected := createTestCountryRecord()
	s.Equal(expected, record)
}
func (s *DatabaseTestSuite) TestReadCountry() {
	country := createTestCountry()
	err := s.database.CreateNewCountry(&country)
	s.Nil(err)
	actual, regionId, subregionId, errR := s.database.ReadCountry(country.NumericCode())
	s.Nil(errR)
	s.Equal(country, actual)
	s.EqualValues(2, regionId)
	s.EqualValues(3, subregionId)
}
func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func createTestCountry() domain.Country {
	result := domain.NewCountry(528, "NL", "NLD")
	result.SetName("Netherlands", "Kingdom of the Netherlands")
	result.SetContinents("Europe")
	result.SetRegion("Europe")
	result.SetSubregion("Western Europe")
	result.SetLatLng(52.5, 5.75)
	result.SetFlag("🇳🇱")
	result.SetPopulation(16655799)
	result.SetArea(41850)
	result.SetStatus("officially-assigned")
	result.SetIndependent(true)
	result.SetLandlocked(false)
	result.SetUnMember(true)
	result.SetStartOfWeek("monday")
	result.SetFifa("NED")
	result.SetOlympicCode("NED")
	result.SetBorders("BEL", "DEU")

	result.SetCapital("Amsterdam")
	result.SetCapitalInfo(domain.LatLng{Lat: 52.35, Lng: 4.92})
	return result
}

func createTestCountryRecord() CountryRecord {
	return CountryRecord{CountryId: 528, Alpha2Code: "NL", Alpha3Code: "NLD", OlympicCode: "NED", FifaCode: "NED",
		Flag: "🇳🇱", Population: 16655799, Area: 41850.0, Independent: true, Landlocked: false, UnMember: true,
		Latitude: 52.5, Longitude: 5.75, RegionId: 2, SubregionId: 3,
		OfficialName: "Kingdom of the Netherlands", CommonName: "Netherlands", StartOfWeek: "monday",
		Status: "officially-assigned"}
}
