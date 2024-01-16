package rdbms

import (
	"database/sql"
	"github.com/stretchr/testify/suite"
	"testing"
)

type currenciesDbTestSuite struct {
	databaseBaseTestSuite
	dtb currenciesDB
}

func (s *currenciesDbTestSuite) SetupSuite() {
	s.databaseBaseTestSuite.SetupSuite()
	s.dtb = currenciesDB{prepStmt: s.db}
}
func (s *currenciesDbTestSuite) TestCreateCurrency() {
	currency, err := s.dtb.createCurrency("USD", "US Dollar", "$")
	s.NoError(err)
	s.Equal("USD", currency.short)
	s.Equal("US Dollar", currency.name)
	s.Equal("$", currency.symbol)
	s.EqualValues(1, currency.currencyId)
}
func (s *currenciesDbTestSuite) TestGetCurrency() {
	expected := currencyRecord{currencyId: 1, short: "USD", name: "US Dollar", symbol: "$"}
	actual, err := s.dtb.createCurrency(expected.short, expected.name, expected.symbol)
	s.NoError(err)
	s.Equal(expected, actual)
	currency, err := s.dtb.getCurrency("USD")
	s.NoError(err)
	s.Equal(expected, currency)
}
func (s *currenciesDbTestSuite) TestGetCurrency_not_found() {
	currency, err := s.dtb.getCurrency("USD")
	s.Error(err)
	s.ErrorIs(err, sql.ErrNoRows)
	s.Equal(currencyRecord{short: "USD"}, currency)
}
func (s *currenciesDbTestSuite) TestReadOrCreateCurrency() {
	expected := currencyRecord{currencyId: 1, short: "EUR", name: "Euro", symbol: "€"}
	actual, err := s.dtb.readOrCreateCurrency(expected.short, expected.name, expected.symbol)
	s.NoError(err)
	s.Equal(expected, actual)
	currency, err := s.dtb.readOrCreateCurrency(expected.short, expected.name, expected.symbol)
	s.NoError(err)
	s.Equal(expected, currency)
}
func (s *currenciesDbTestSuite) TestReadCurrencies() {
	ids := []uint32{1, 2}
	expected := []currencyRecord{
		{currencyId: 1, short: "USD", name: "US Dollar", symbol: "$"},
		{currencyId: 2, short: "EUR", name: "Euro", symbol: "€"}}
	_, er := s.dtb.readOrCreateCurrency("USD", "US Dollar", "$")
	s.NoError(er)
	_, er = s.dtb.readOrCreateCurrency("EUR", "Euro", "€")
	s.NoError(er)
	actual, err := s.dtb.readCurrencies(ids...)
	s.NoError(err)
	s.Equal(expected, actual)
}
func Test_currenciesDbTestSuite(t *testing.T) {
	suite.Run(t, new(currenciesDbTestSuite))
}
