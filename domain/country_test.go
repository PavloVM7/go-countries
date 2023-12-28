package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
