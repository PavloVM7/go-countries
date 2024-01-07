package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
