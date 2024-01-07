package domain

import (
	"pm.com/go-countries/tools"
	"strings"
)

type countryExt struct {
	car             Car
	continents      []string
	borders         []string
	capital         []string
	topLevelDomains []string
	altSpellings    []string
	timezones       []string
	capitalInfo     []LatLng
	currencies      []Currency
	languages       []Language
	translations    []Translation
	demonyms        []Demonym
	flags           []CodeDescription
}

func (c *countryExt) AddFlag(language, flag string) {
	language = strings.TrimSpace(language)
	flag = strings.TrimSpace(flag)
	c.flags = append(c.flags, CodeDescription{Code: language, Description: flag})
}
func (c *countryExt) Flags() []CodeDescription {
	return tools.CopyArray(c.flags)
}
func (c *countryExt) SetCar(side string, signs ...string) {
	c.car.Side = strings.TrimSpace(side)
	c.car.Signs = tools.CopyStringArraySkipEmpty(signs)
}
func (c *countryExt) Car() Car {
	return c.car
}

func (c *countryExt) AddDemonym(language, f, m string) {
	language = strings.TrimSpace(language)
	f = strings.TrimSpace(f)
	m = strings.TrimSpace(m)
	c.demonyms = append(c.demonyms, Demonym{Language: language, F: f, M: m})
}
func (c *countryExt) Demonyms() []Demonym {
	return tools.CopyArray(c.demonyms)
}
func (c *countryExt) AddTranslation(language, common, official string, native bool) {
	language = strings.TrimSpace(language)
	common = strings.TrimSpace(common)
	official = strings.TrimSpace(official)
	c.translations = append(c.translations,
		Translation{Language: language, Common: common, Official: official, Native: native})
}
func (c *countryExt) Translations() []Translation {
	return tools.CopyArray(c.translations)
}
func (c *countryExt) AddLanguage(short, name string) {
	short = strings.TrimSpace(short)
	name = strings.TrimSpace(name)
	c.languages = append(c.languages, Language{Short: short, Name: name})
}
func (c *countryExt) Languages() []Language {
	return tools.CopyArray(c.languages)
}
func (c *countryExt) SetTimezones(timezones ...string) {
	c.timezones = tools.CopyStringArraySkipEmpty(timezones)
}
func (c *countryExt) Timezones() []string {
	return tools.CopyArray(c.timezones)
}
func (c *countryExt) SetAltSpellings(spellings ...string) {
	c.altSpellings = tools.CopyStringArraySkipEmpty(spellings)
}
func (c *countryExt) AltSpellings() []string {
	return tools.CopyArray(c.altSpellings)
}
func (c *countryExt) SetTopLevelDomains(tlds ...string) {
	c.topLevelDomains = tools.CopyStringArraySkipEmpty(tlds)
}
func (c *countryExt) TopLevelDomains() []string {
	return tools.CopyArray(c.topLevelDomains)
}
func (c *countryExt) AddCurrency(short, name, symbol string) {
	short = strings.TrimSpace(short)
	name = strings.TrimSpace(name)
	symbol = strings.TrimSpace(symbol)
	c.currencies = append(c.currencies, Currency{Short: short, Name: name, Symbol: symbol})
}
func (c *countryExt) Currencies() []Currency {
	return tools.CopyArray(c.currencies)
}
func (c *countryExt) SetCapitalInfo(latLngs ...LatLng) {
	c.capitalInfo = tools.CopyArray(latLngs)
}
func (c *countryExt) CapitalInfo() []LatLng {
	return tools.CopyArray(c.capitalInfo)
}
func (c *countryExt) SetCapital(capitals ...string) {
	c.capital = tools.CopyStringArraySkipEmpty(capitals)
}
func (c *countryExt) Capital() []string {
	return tools.CopyArray(c.capital)
}

func (c *countryExt) SetContinents(continents ...string) {
	c.continents = tools.CopyStringArraySkipEmpty(continents)
}
func (c *countryExt) Continents() []string {
	return tools.CopyArray(c.continents)
}
func (c *countryExt) Borders() []string {
	if len(c.borders) == 0 {
		return []string{}
	}
	result := make([]string, len(c.borders))
	copy(result, c.borders)
	return result
}
func (c *countryExt) SetBorders(borders ...string) {
	c.borders = make([]string, 0, len(borders))
	for _, border := range borders {
		b := strings.TrimSpace(border)
		if len(b) > 0 {
			c.borders = append(c.borders, b)
		}
	}
}
