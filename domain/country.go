package domain

import "strings"

type Country struct {
	area        float32
	population  uint32
	numericCode uint16
	independent bool
	landlocked  bool
	unMember    bool
	alpha2Code  string
	alpha3Code  string
	region      string
	subregion   string
	olympicCode string
	fifa        string
	flag        string
	startOfWeek string // Day of the start of week (Sunday/Monday)
	status      string // assigment status ( officially-assigned user-assigned)
	name        countryName
	latLng      LatLng
	countryExt
}

// SetStartOfWeek is used to set the start of the week in the country
func (c *Country) SetStartOfWeek(startOfWeek string) {
	c.startOfWeek = startOfWeek
}

// StartOfWeek is used to get the start of the week in the country
func (c *Country) StartOfWeek() string {
	return c.startOfWeek
}

// SetStatus is used to set the status of the country
func (c *Country) SetStatus(status string) {
	c.status = status
}

// Status is used to get the status of the country
func (c *Country) Status() string {
	return c.status
}

// SetLatLng is used to set the latitude and longitude of the country
func (c *Country) SetLatLng(lat, lng float32) {
	c.latLng.Lat = lat
	c.latLng.Lng = lng
}

// LatLng is used to get the latitude and longitude of the country
func (c *Country) LatLng() (lat, lng float32) {
	lat, lng = c.latLng.Lat, c.latLng.Lng
	return
}

// Flag is used to get the flag emoji of the country
func (c *Country) Flag() string {
	return c.flag
}

// SetFlag is used to set the flag emoji of the country
func (c *Country) SetFlag(flag string) {
	c.flag = flag
}

// Region is used to get the UN demographic region of the country
func (c *Country) Region() string {
	return c.region
}

// Subregion is used to get the UN demographic subregion of the country
func (c *Country) Subregion() string {
	return c.subregion
}

// SetSubregion is used to set the UN demographic subregion of the country
func (c *Country) SetSubregion(subregion string) {
	c.subregion = strings.TrimSpace(subregion)
}

// SetRegion is used to set the UN demographic region of the country
func (c *Country) SetRegion(region string) {
	c.region = strings.TrimSpace(region)
}

// IsUnMember is used to check if the country is a UN member
func (c *Country) IsUnMember() bool {
	return c.unMember
}

// SetUnMember is used to set the country UN member status
func (c *Country) SetUnMember(unMember bool) {
	c.unMember = unMember
}
func (c *Country) IsLandlocked() bool {
	return c.landlocked
}
func (c *Country) SetLandlocked(landlocked bool) {
	c.landlocked = landlocked
}
func (c *Country) Fifa() string {
	return c.fifa
}
func (c *Country) SetFifa(fifa string) {
	c.fifa = fifa
}
func (c *Country) IsIndependent() bool {
	return c.independent
}
func (c *Country) SetIndependent(independent bool) {
	c.independent = independent
}
func (c *Country) OlympicCode() string {
	return c.olympicCode
}
func (c *Country) SetOlympicCode(olympicCode string) {
	c.olympicCode = olympicCode
}
func (c *Country) Population() uint32 {
	return c.population
}
func (c *Country) SetPopulation(population uint32) {
	c.population = population
}
func (c *Country) SetArea(area float32) {
	c.area = area
}
func (c *Country) Area() float32 {
	return c.area
}
func (c *Country) SetName(common, official string) {
	c.name = countryName{common: common, official: official}
}
func (c *Country) CommonName() string {
	return c.name.common
}
func (c *Country) OfficialName() string {
	return c.name.official
}
func (c *Country) Alpha3Code() string {
	return c.alpha3Code
}
func (c *Country) Alpha2Code() string {
	return c.alpha2Code
}
func (c *Country) NumericCode() uint16 {
	return c.numericCode
}
func NewCountry(numericCode uint16, alpha2Code string, alpha3Code string) Country {
	return Country{
		numericCode: numericCode,
		alpha2Code:  alpha2Code,
		alpha3Code:  alpha3Code,
	}
}

type CountryExtended struct {
	Borders []string
}
