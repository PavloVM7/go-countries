package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryExt_AddTranslation(t *testing.T) {
	var cext countryExt
	cext.AddTranslation("nld ", " Nederland", "Koninkrijk der Nederlanden ", true)
	cext.AddTranslation("deu", "Niederlande", "Niederlande", false)
	actual := cext.Translations()
	assert.Equal(t, []Translation{
		{"nld", "Nederland", "Koninkrijk der Nederlanden", true},
		{"deu", "Niederlande", "Niederlande", false}}, actual)
}

func TestCountryExt_AddLanguage(t *testing.T) {
	var cext countryExt
	cext.AddLanguage(" eng ", " English  ")
	cext.AddLanguage(" nld", " Dutch ")
	cext.AddLanguage("fra ", " French  ")
	actual := cext.Languages()
	assert.Equal(t, []Language{
		{"eng", "English"},
		{"nld", "Dutch"},
		{"fra", "French"}}, actual)
}

func TestCountryExt_SetTimezones(t *testing.T) {
	var cext countryExt
	cext.SetTimezones("UTC-04:00", "UTC-03:00", " UTC-01:00", " UTC UTC+01:00 ", " ", "")
	actual := cext.Timezones()
	assert.Equal(t, []string{"UTC-04:00", "UTC-03:00", "UTC-01:00", "UTC UTC+01:00"}, actual)
}

func TestCountryExt_SetAltSpellings(t *testing.T) {
	var cext countryExt
	cext.SetAltSpellings("", "NL", " Holland ", "Nederland ", " The Netherlands ")
	actual := cext.AltSpellings()
	assert.Equal(t, []string{"NL", "Holland", "Nederland", "The Netherlands"}, actual)
}

func TestCountryExt_SetTopLevelDomains(t *testing.T) {
	var cext countryExt
	cext.SetTopLevelDomains(" .fr  ", " .gp ", "\n", " ", "\t")
	actual := cext.TopLevelDomains()
	assert.Equal(t, []string{".fr", ".gp"}, actual)
}

func TestCountryExt_AddCurrency(t *testing.T) {
	var cext countryExt
	cext.AddCurrency("EUR", "Euro", "€")
	cext.AddCurrency("USD", "United States dollar", "$")
	cext.AddCurrency("GBP", "British pound", "£")
	actual := cext.Currencies()
	assert.Equal(t, []Currency{
		{"EUR", "Euro", "€"},
		{"USD", "United States dollar", "$"},
		{"GBP", "British pound", "£"},
	}, actual)
}

func TestCountryExt_SetCapitalInfo(t *testing.T) {
	var cext countryExt
	expected := []LatLng{{Lat: 52.35, Lng: 4.92}}
	cext.SetCapitalInfo(expected...)
	actual := cext.CapitalInfo()
	assert.Equal(t, expected, actual)
}

func TestCountryExt_SetCapital(t *testing.T) {
	var cext countryExt
	capitals := []string{"London ", " Dublin", ""}
	cext.SetCapital(capitals...)
	actual := cext.Capital()
	assert.Equal(t, []string{"London", "Dublin"}, actual)
}

func TestCountryExt_SetContinents(t *testing.T) {
	var cext countryExt
	continents := []string{"Europe", "Asia", ""}
	cext.SetContinents(continents...)
	actual := cext.Continents()
	assert.Equal(t, []string{"Europe", "Asia"}, actual)
}

func TestCountryExt_SetBorders(t *testing.T) {
	var cext countryExt
	expected := []string{"AND", "BEL", "DEU", "ITA", "LUX", "MCO", "ESP", "CHE"}
	cext.SetBorders(expected...)
	actual := cext.Borders()
	assert.Equal(t, expected, actual)
}
