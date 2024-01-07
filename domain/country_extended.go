package domain

import (
	"pm.com/go-countries/tools"
	"strings"
)

type countryExt struct {
	continents  []string
	borders     []string
	capital     []string
	capitalInfo []LatLng
	currencies  []Currency
}

func (c *countryExt) AddCurrency(short, name, symbol string) {
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
