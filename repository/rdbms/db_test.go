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
	res, err := s.database.db.Exec(`TRUNCATE translations, country_continents, countries, languages, regions,
    borders, top_level_domains
    RESTART IDENTITY CASCADE;`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables truncated.")
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
	s.Equal(LanguageRecord{LanguageId: 0, Language: lang3}, langD)
}

func (s *DatabaseTestSuite) TestGetLanguage() {
	langEng := "eng"
	lang, err := s.database.CreateLanguage(langEng)
	s.Nil(err)
	s.Equal(LanguageRecord{LanguageId: 1, Language: langEng}, lang)
	actual, er := s.database.GetLanguage(langEng)
	s.Nil(er)
	s.Equal(lang, actual)
}
func (s *DatabaseTestSuite) TestGetLanguageNotExists() {
	lng := "any"
	lang, err := s.database.GetLanguage(lng)
	s.NotNil(err)
	s.Equal(LanguageRecord{LanguageId: 0, Language: lng}, lang)
}
func (s *DatabaseTestSuite) TestCreateContinent() {
	name := "Europe"
	reg, err := s.database.CreateContinent(name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 1, ParentId: 0, RegionName: name}, reg)
}
func (s *DatabaseTestSuite) TestCreateRegion() {
	name := "Europe"
	continent, err := s.database.CreateContinent(name)
	s.Nil(err)
	s.Equal(RegionRecord{RegionId: 1, ParentId: 0, RegionName: name}, continent)
	reg, errReg := s.database.CreateRegion(name, continent.RegionId)
	s.Nil(errReg)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, reg)
}
func (s *DatabaseTestSuite) TestCreateSubregion() {
	name := "Europe"
	continent, err := s.database.CreateContinent(name)
	s.Nil(err)
	reg, errReg := s.database.CreateRegion(name, continent.RegionId)
	s.Nil(errReg)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, reg)
	subregionName := "Western Europe"
	subregion, errSubreg := s.database.CreateRegion(subregionName, reg.RegionId)
	s.Nil(errSubreg)
	s.Equal(RegionRecord{RegionId: 3, ParentId: 2, RegionName: subregionName}, subregion)
}
func (s *DatabaseTestSuite) TestGetContinent() {
	name := "Europe"
	continent, errC := s.database.CreateContinent(name)
	s.Nil(errC)
	s.Equal(RegionRecord{RegionId: 1, ParentId: 0, RegionName: name}, continent)
	region, errR := s.database.CreateRegion(name, continent.RegionId)
	s.Nil(errR)
	subregionName := "Western Europe"
	_, errS := s.database.CreateRegion(subregionName, region.RegionId)
	s.Nil(errS)
	actual, err := s.database.GetContinent(name)
	s.Nil(err)
	s.Equal(continent, actual)
}
func (s *DatabaseTestSuite) TestGetRegion() {
	name := "Europe"
	continent, errC := s.database.CreateContinent(name)
	s.Nil(errC)
	region, errR := s.database.CreateRegion(name, continent.RegionId)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, region)
	s.Nil(errR)
	subregionName := "Western Europe"
	_, errS := s.database.CreateRegion(subregionName, region.RegionId)
	s.Nil(errS)
	actual, err := s.database.GetRegion(name)
	s.Nil(err)
	s.Equal(region, actual)
}
func (s *DatabaseTestSuite) TestGetSubregion() {
	name := "Europe"
	continent, errC := s.database.CreateContinent(name)
	s.Nil(errC)
	region, errR := s.database.CreateRegion(name, continent.RegionId)
	s.Equal(RegionRecord{RegionId: 2, ParentId: 1, RegionName: name}, region)
	s.Nil(errR)
	subregionName := "Western Europe"
	subregion, errS := s.database.CreateRegion(subregionName, region.RegionId)
	s.Nil(errS)
	s.Equal(RegionRecord{RegionId: 3, ParentId: 2, RegionName: subregionName}, subregion)
	actual, err := s.database.GetSubregion(subregionName)
	s.Nil(err)
	s.Equal(subregion, actual)
}
func (s *DatabaseTestSuite) TestCrateCountry() {
	name := "Europe"
	subregionName := "Western Europe"
	_, _, _, err := s.createRegions(name, name, subregionName)
	s.Nil(err)
	record := createTestCountry()
	err = s.database.CreateCountry(&record)
	s.Nil(err)
	actual, errG := s.database.GetCountry(record.CountryId)
	s.Nil(errG)
	s.Equal(record, actual)
}
func (s *DatabaseTestSuite) createRegions(continentName, regionName, subregionName string) (continent, region, subregion RegionRecord, err error) {
	continent, err = s.database.CreateContinent(continentName)
	if err != nil {
		err = fmt.Errorf("create continent error: %w", err)
		return
	}
	region, err = s.database.CreateRegion(regionName, continent.RegionId)
	if err != nil {
		err = fmt.Errorf("create region error: %w", err)
		return
	}
	subregion, err = s.database.CreateRegion(subregionName, region.RegionId)
	if err != nil {
		err = fmt.Errorf("create soubregion error: %w", err)
	}
	return
}
func (s *DatabaseTestSuite) createCountry(continent, region, subregion string, country CountryRecord) error {
	_, r, sr, err := s.createRegions(continent, region, subregion)
	country.RegionId = r.RegionId
	country.SubregionId = sr.RegionId
	if err != nil {
		return err
	}
	err = s.database.CreateCountry(&country)
	if err != nil {
		return fmt.Errorf("create country error: %w", err)
	}
	return nil
}
func (s *DatabaseTestSuite) TestGetCountryNotFound() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	actual, errGet1 := s.database.GetCountry(country.CountryId)
	s.Nil(errGet1)
	s.Equal(country, actual)
	idNotExists := country.CountryId + 1
	_, errG2 := s.database.GetCountry(idNotExists)
	s.NotNil(errG2)
	s.Equal(fmt.Sprintf("country not found by id=%d", idNotExists), errG2.Error())
}
func (s *DatabaseTestSuite) TestCreateTranslation() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	countryId := country.CountryId
	language := "nld"
	official := "Koninkrijk der Nederlanden"
	common := "Nederland"
	lng, errL := s.database.CreateLanguage(language)
	s.Nil(errL)
	actual, errT := s.database.CreateTransaction(countryId, lng.LanguageId, true, official, common)
	s.Nil(errT)
	s.Equal(TranslationRecord{Id: 1, CountryId: countryId, LanguageId: lng.LanguageId, Native: true,
		OfficialName: official, CommonName: common}, actual)
}
func (s *DatabaseTestSuite) TestCreateTranslationNotNativeName() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	countryId := country.CountryId
	language := "nld"
	official := "Koninkrijk der Nederlanden"
	common := "Nederland"
	lng, errL := s.database.CreateLanguage(language)
	s.Nil(errL)
	actual, errT := s.database.CreateTransaction(countryId, lng.LanguageId, true, official, common)
	s.Nil(errT)
	s.Equal(TranslationRecord{Id: 1, CountryId: countryId, LanguageId: lng.LanguageId, Native: true,
		OfficialName: official, CommonName: common}, actual)
	actual2, errT2 := s.database.CreateTransaction(countryId, lng.LanguageId, false, official, common)
	s.Nil(errT2)
	s.Equal(TranslationRecord{Id: 2, CountryId: countryId, LanguageId: lng.LanguageId, Native: false,
		OfficialName: official, CommonName: common}, actual2)
}
func (s *DatabaseTestSuite) TestCreateBorders() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual, errB := s.database.CreteBorders(country.CountryId, borders...)
	s.Nil(errB)
	expected := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected, actual)
}
func (s *DatabaseTestSuite) TestCreateBorders1duplicate() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.database.CreteBorders(country.CountryId, borders[0])
	s.Nil(err1)
	expected1 := []BorderRecord{{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]}}
	s.Equal(expected1, actual1)
	actual2, err2 := s.database.CreteBorders(country.CountryId, borders...)
	s.NotNil(err2)
	expected2 := []BorderRecord{{Id: 3, CountryId: country.CountryId, Alpha3Code: borders[1]}}
	s.Equal(expected2, actual2)
}
func (s *DatabaseTestSuite) TestCreateBordersDuplicate() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.database.CreteBorders(country.CountryId, borders...)
	s.Nil(err1)
	expected1 := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected1, actual1)
	actual2, err2 := s.database.CreteBorders(country.CountryId, borders...)
	s.NotNil(err2)
	s.NotNil(actual2)
	s.Equal(0, len(actual2))
}
func (s *DatabaseTestSuite) TestGetBorders() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.database.CreteBorders(country.CountryId, borders...)
	s.Nil(err1)
	expected1 := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected1, actual1)
	actual2, err2 := s.database.GetBorders(country.CountryId)
	s.Nil(err2)
	s.Equal(expected1, actual2)
}
func (s *DatabaseTestSuite) TestGetBordersNotExist() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	borders := []string{"BEL", "DEU"}
	actual1, err1 := s.database.CreteBorders(country.CountryId, borders...)
	s.Nil(err1)
	expected1 := []BorderRecord{
		{Id: 1, CountryId: country.CountryId, Alpha3Code: borders[0]},
		{Id: 2, CountryId: country.CountryId, Alpha3Code: borders[1]},
	}
	s.Equal(expected1, actual1)
	actual2, err2 := s.database.GetBorders(country.CountryId + 1)
	s.Nil(err2)
	s.Equal([]BorderRecord{}, actual2)
}
func (s *DatabaseTestSuite) TestCreateTopLevelDomains() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.database.CreateTopLevelDomains(country.CountryId, tlds...)
	s.Nil(errQ)
	expected := []TldRecord{
		{Id: 1, CountryId: country.CountryId, Tld: ".nl"},
		{Id: 2, CountryId: country.CountryId, Tld: ".nld"},
	}
	s.Equal(expected, actual)
}
func (s *DatabaseTestSuite) TestGetTopLevelDomains() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.database.CreateTopLevelDomains(country.CountryId, tlds...)
	s.Nil(errQ)
	expected := []TldRecord{
		{Id: 1, CountryId: country.CountryId, Tld: ".nl"},
		{Id: 2, CountryId: country.CountryId, Tld: ".nld"},
	}
	s.Equal(expected, actual)
	actual2, errQ2 := s.database.GetTopLevelDomains(country.CountryId)
	s.Nil(errQ2)
	s.Equal(actual, actual2)
}
func (s *DatabaseTestSuite) TestGetTopLevelDomainsNotExist() {
	country := createTestCountry()
	err := s.createCountry("Europe", "Europe", "Western Europe", country)
	s.Nil(err)
	tlds := []string{".nl ", ".nld"}
	actual, errQ := s.database.CreateTopLevelDomains(country.CountryId, tlds...)
	s.Nil(errQ)
	expected := []TldRecord{
		{Id: 1, CountryId: country.CountryId, Tld: ".nl"},
		{Id: 2, CountryId: country.CountryId, Tld: ".nld"},
	}
	s.Equal(expected, actual)
	actual2, errQ2 := s.database.GetTopLevelDomains(country.CountryId + 1)
	s.Nil(errQ2)
	s.Equal([]TldRecord{}, actual2)
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
