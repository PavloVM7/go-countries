package domain

import (
	"github.com/PavloVM7/go-collections/pkg/collections"
	"pm.com/go-countries/tools"
	"slices"
	"sort"
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
	maps            []CodeDescription
	coatOfArms      []CodeDescription
	// ToDo: International dialing codes (callingCodes / idd)
}

func (c *countryExt) AddCoatOfArm(picType, ref string) {
	picType = strings.TrimSpace(picType)
	ref = strings.TrimSpace(ref)
	if picType == "" || ref == "" {
		return
	}
	c.coatOfArms = append(c.coatOfArms, CodeDescription{Code: picType, Description: ref})
	sort.Slice(c.coatOfArms, func(i, j int) bool { return c.coatOfArms[i].Code < c.coatOfArms[j].Code })
}
func (c *countryExt) CoatOfArms() []CodeDescription {
	return tools.CopyArray(c.coatOfArms)
}
func (c *countryExt) AddMap(name, ref string) {
	name = strings.TrimSpace(name)
	ref = strings.TrimSpace(ref)
	if name == "" || ref == "" {
		return
	}
	c.maps = append(c.maps, CodeDescription{Code: name, Description: ref})
	sort.Slice(c.maps, func(i, j int) bool { return c.maps[i].Code < c.maps[j].Code })
}
func (c *countryExt) Maps() []CodeDescription {
	return tools.CopyArray(c.maps)
}
func (c *countryExt) AddFlag(pic, ref string) {
	pic = strings.TrimSpace(pic)
	ref = strings.TrimSpace(ref)
	if pic == "" || ref == "" {
		return
	}
	c.flags = append(c.flags, CodeDescription{Code: pic, Description: ref})
	sort.Slice(c.flags, func(i, j int) bool { return c.flags[i].Code < c.flags[j].Code })
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
	if language == "" || f == "" || m == "" {
		return
	}
	c.demonyms = append(c.demonyms, Demonym{Language: language, F: f, M: m})
	sort.Slice(c.demonyms, func(i, j int) bool { return c.demonyms[i].Language < c.demonyms[j].Language })
}
func (c *countryExt) Demonyms() []Demonym {
	return tools.CopyArray(c.demonyms)
}
func (c *countryExt) AddTranslation(language, common, official string) {
	c.addTranslationInner(language, common, official, false)
}
func (c *countryExt) AddNativeName(language, common, official string) {
	c.addTranslationInner(language, common, official, true)
}
func (c *countryExt) addTranslationInner(language, common, official string, native bool) {
	language = strings.TrimSpace(language)
	common = strings.TrimSpace(common)
	official = strings.TrimSpace(official)
	if language == "" || common == "" || official == "" {
		return
	}
	c.translations = append(c.translations,
		Translation{Language: language, Common: common, Official: official, Native: native})
	sort.Slice(c.translations, func(i, j int) bool {
		if c.translations[i].Language == c.translations[j].Language {
			return c.translations[i].Native
		}
		return c.translations[i].Language < c.translations[j].Language
	})

}
func (c *countryExt) Translations() []Translation {
	result := make([]Translation, 0, len(c.translations))
	for _, translation := range c.translations {
		if !translation.Native {
			result = append(result, translation)
		}
	}
	return result
}
func (c *countryExt) NativeNames() []Translation {
	result := make([]Translation, 0, len(c.translations))
	for _, translation := range c.translations {
		if translation.Native {
			result = append(result, translation)
		}
	}
	slices.Clip(result)
	return result
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
	filtered := tools.CopyStringArraySkipEmpty(capitals)
	c.capital = make([]string, 0, len(filtered))
	set := collections.NewSetItems[string](filtered...)
	for _, capital := range filtered {
		if set.Contains(capital) {
			c.capital = append(c.capital, capital)
			set.Remove(capital)
		}
	}
	c.capital = slices.Clip(c.capital)
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
