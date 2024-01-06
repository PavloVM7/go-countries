package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountry_SetStartOfWeekend(t *testing.T) {
	country := createTestCountryFRA()
	startOfWeek := "monday"
	country.SetStartOfWeek(startOfWeek)
	assert.Equal(t, startOfWeek, country.StartOfWeek())
}

func TestCountry_SetStatus(t *testing.T) {
	country := createTestCountryFRA()
	status := "officially-assigned"
	country.SetStatus(status)
	assert.Equal(t, status, country.Status())
}

func TestCountry_LatLng(t *testing.T) {
	country := createTestCountryFRA()
	latExp := float32(46.0)
	lngExp := float32(2.0)
	country.SetLatLng(latExp, lngExp)
	lat, lng := country.LatLng()
	assert.Equal(t, latExp, lat)
	assert.Equal(t, lngExp, lng)
}

func TestCountry_SetFlag(t *testing.T) {
	flag := "🇳🇱"
	country := createTestCountryNLD()
	country.SetFlag(flag)
	assert.Equal(t, flag, country.Flag())
}

func TestCountry_SetSubregion(t *testing.T) {
	subregion := "Western Europe"
	country := createTestCountryNLD()
	country.SetSubregion(subregion)
	assert.Equal(t, subregion, country.Subregion())
}
func TestCountry_SetRegion(t *testing.T) {
	region := "Europe"
	country := createTestCountryNLD()
	country.SetRegion(region)
	assert.Equal(t, region, country.Region())
}

func TestCountry_SetUnMember(t *testing.T) {
	country := createTestCountryNLD()
	country.SetUnMember(true)
	assert.True(t, country.IsUnMember())
}

func TestCountry_SetLandlocked(t *testing.T) {
	country := createTestCountryNLD()
	country.SetLandlocked(false)
	assert.False(t, country.IsLandlocked())
}

func TestCountry_SetFifa(t *testing.T) {
	fifa := "NED"
	country := createTestCountryNLD()
	country.SetFifa(fifa)
	assert.Equal(t, fifa, country.Fifa())
}

func TestCountry_SetIndependent(t *testing.T) {
	country := createTestCountryNLD()
	country.SetIndependent(true)
	assert.True(t, country.IsIndependent())
}
func TestCountry_SetOlympicCode(t *testing.T) {
	olympicCode := "NED"
	country := createTestCountryNLD()
	country.SetOlympicCode(olympicCode)
	assert.Equal(t, olympicCode, country.OlympicCode())
}

func TestCountry_SetPopulation(t *testing.T) {
	population := uint32(16655799)
	country := createTestCountryNLD()
	country.SetPopulation(population)
	assert.Equal(t, population, country.Population())
}

func TestCountry_SetArea(t *testing.T) {
	country := createTestCountryNLD()
	area := float32(41850.0)
	country.SetArea(area)
	assert.Equal(t, area, country.Area())
}

func TestCountry_CommonName(t *testing.T) {
	country := createTestCountryNLD()
	common := "Netherlands"
	official := "Kingdom of the Netherlands"
	country.SetName(common, official)
	assert.Equal(t, common, country.CommonName())
}
func TestCountry_OfficialName(t *testing.T) {
	country := createTestCountryNLD()
	common := "Netherlands"
	official := "Kingdom of the Netherlands"
	country.SetName(common, official)
	assert.Equal(t, official, country.OfficialName())
}

func TestCountry_SetName(t *testing.T) {
	country := createTestCountryNLD()
	common := "Netherlands"
	official := "Kingdom of the Netherlands"
	country.SetName(common, official)
	expected := Country{
		numericCode: 528,
		alpha2Code:  "NL",
		alpha3Code:  "NLD",
		name:        countryName{common: common, official: official},
	}
	assert.Equal(t, expected, country)
}

func TestCountry_Alpha3Code(t *testing.T) {
	country := createTestCountryNLD()
	assert.Equal(t, "NLD", country.Alpha3Code())
}

func TestCountry_Alpha2Code(t *testing.T) {
	country := createTestCountryNLD()
	assert.Equal(t, "NL", country.Alpha2Code())
}

func TestNewCountry(t *testing.T) {
	expected := Country{
		numericCode: 248,
		alpha2Code:  "AX",
		alpha3Code:  "ALA",
	}
	actual := NewCountry(248, "AX", "ALA")
	assert.Equal(t, expected, actual)
}

func createTestCountryNLD() Country {
	return NewCountry(528, "NL", "NLD")
}
func createTestCountryFRA() Country {
	return NewCountry(250, "FR", "FRA")
}
