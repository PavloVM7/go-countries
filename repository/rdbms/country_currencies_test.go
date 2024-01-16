package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type countryCurrenciesDbTestSuite struct {
	databaseBaseTestSuite
	countryRecord CountryRecord
	currencies    []currencyRecord
	dtb           countryCurrenciesDB
}

func (s *countryCurrenciesDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = countryCurrenciesDB{prepStmt: s.db}
}
func (s *countryCurrenciesDbTestSuite) SetupTest() {
	s.countryRecord = createTestCountryRecord()
	err := s.createCountry("Europe", "Europe", "Western Europe", s.countryRecord)
	s.Nil(err)
	fmt.Printf("country %d:%s:%s:'%s' created\n", s.countryRecord.CountryId, s.countryRecord.Alpha2Code,
		s.countryRecord.Alpha3Code, s.countryRecord.CommonName)
	cdb := currenciesDB{prepStmt: s.db}
	record, er := cdb.readOrCreateCurrency("EUR", "Euro", "€")
	s.NoError(er)
	s.currencies = append(s.currencies, record)
	record, er = cdb.readOrCreateCurrency("USD", "US Dollar", "$")
	s.NoError(er)
	s.currencies = append(s.currencies, record)
}
func (s *countryCurrenciesDbTestSuite) TearDownTest() {
	s.databaseBaseTestSuite.TearDownTest()
	s.countryRecord = CountryRecord{}
	s.currencies = nil
}
func (s *countryCurrenciesDbTestSuite) getCurrenciesIds() []uint32 {
	result := make([]uint32, 0, len(s.currencies))
	for _, currency := range s.currencies {
		result = append(result, currency.currencyId)
	}
	return result
}

func (s *countryCurrenciesDbTestSuite) Test_createCountryCurrencies() {
	ids := s.getCurrenciesIds()
	records, err := s.dtb.createCountryCurrencies(s.countryRecord.CountryId, ids...)
	s.NoError(err)
	expected := []countryCurrenciesRecord{
		{id: 1, countryId: s.countryRecord.CountryId, currencyId: 1},
		{id: 2, countryId: s.countryRecord.CountryId, currencyId: 2}}
	s.Equal(expected, records)
}

func (s *countryCurrenciesDbTestSuite) Test_readCountryCurrencies() {
	ids := s.getCurrenciesIds()
	records, err := s.dtb.createCountryCurrencies(s.countryRecord.CountryId, ids...)
	s.NoError(err)
	expected := []countryCurrenciesRecord{
		{id: 1, countryId: s.countryRecord.CountryId, currencyId: 1},
		{id: 2, countryId: s.countryRecord.CountryId, currencyId: 2}}
	s.Equal(expected, records)
	actual, errR := s.dtb.readCountryCurrencies(s.countryRecord.CountryId)
	s.NoError(errR)
	s.Equal(expected, actual)
}
func Test_countryCurrenciesDbTestSuite(t *testing.T) {
	suite.Run(t, new(countryCurrenciesDbTestSuite))
}
