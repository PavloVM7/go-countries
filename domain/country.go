package domain

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
	name        countryName
}

func (c *Country) SetRegion(region, subregion string) {
	c.region = region
	c.subregion = subregion
}
func (c *Country) IsUnMember() bool {
	return c.unMember
}
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
