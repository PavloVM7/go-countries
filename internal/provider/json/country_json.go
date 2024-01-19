package json

import (
	"encoding/json"
	"pm.com/go-countries/domain"
)

type CountryJSON struct {
	LatLng       *LatLngJson         `json:"latlng"`
	CapitalInfo  *CapitalInfoJson    `json:"capitalInfo"`
	Population   uint32              `json:"population"`
	Area         float32             `json:"area"`
	NumericCode  uint16              `json:"ccn3,string"`
	Independent  bool                `json:"independent"`
	UnMember     bool                `json:"unMember"`
	Landlocked   bool                `json:"landlocked"`
	Name         CountryNameJson     `json:"name"`
	Currencies   CurrenciesJson      `json:"currencies"`
	CallingCodes CollingCodesJson    `json:"idd"`
	Alpha2Code   string              `json:"cca2"`
	Alpha3Code   string              `json:"cca3"`
	Cioc         string              `json:"cioc"`
	Fifa         string              `json:"fifa"`
	Region       string              `json:"region"`
	Subregion    string              `json:"subregion"`
	Status       string              `json:"status"`
	StartOfWeek  string              `json:"startOfWeek"`
	Flag         string              `json:"flag"`
	Tld          []string            `json:"tld"`
	Capital      []string            `json:"capital"`
	AltSpellings []string            `json:"altSpellings"`
	Continents   []string            `json:"continents"`
	Timezones    []string            `json:"timezones"`
	Borders      []string            `json:"borders"`
	Gini         map[string]float32  `json:"gini"`
	Languages    map[string]string   `json:"languages"`
	Translations NamesJson           `json:"translations"`
	Demonyms     map[string]Demonyms `json:"demonyms"`
	Maps         map[string]string   `json:"maps"`
	CoatOfArms   map[string]string   `json:"coatOfArms"`
	Flags        map[string]string   `json:"flags"`
	Car          CarJson             `json:"car"`
}

type CountryNameJson struct {
	NameJson
	NativeName NamesJson `json:"nativeName"`
}

type NamesJson map[string]NameJson

func (nms *NamesJson) UnmarshalJSON(data []byte) error {
	var res map[string]NameJson
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*nms = res
	return nil
}

type NameJson struct {
	Official string `json:"official"`
	Common   string `json:"common"`
}
type CapitalInfoJson struct {
	LatLng *LatLngJson `json:"latlng"`
}

type LatLngJson domain.LatLng

func (ll *LatLngJson) UnmarshalJSON(data []byte) error {
	var ar [2]float32
	if err := json.Unmarshal(data, &ar); err != nil {
		return err
	}
	ll.Lat = ar[0]
	ll.Lng = ar[1]
	return nil
}

type CurrenciesJson map[string]CurrencyJson

func (c *CurrenciesJson) UnmarshalJSON(data []byte) error {
	var res map[string]CurrencyJson
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*c = res
	return nil
}

type CurrencyJson struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type CollingCodesJson struct {
	Root     string   `json:"root"`
	Suffixes []string `json:"suffixes"`
}

type Demonyms struct {
	F string
	M string
}
type CarJson struct {
	Signs []string
	Side  string
}

func createCountryFromJson(jsonCountry *CountryJSON) domain.Country {
	result := domain.NewCountry(jsonCountry.NumericCode, jsonCountry.Alpha2Code, jsonCountry.Alpha3Code)
	result.SetName(jsonCountry.Name.Common, jsonCountry.Name.Official)
	for k, v := range jsonCountry.Name.NativeName {
		result.AddNativeName(k, v.Common, v.Official)
	}
	result.SetTopLevelDomains(jsonCountry.Tld...)
	result.SetOlympicCode(jsonCountry.Cioc)
	result.SetIndependent(jsonCountry.Independent)
	result.SetStatus(jsonCountry.Status)
	result.SetUnMember(jsonCountry.UnMember)
	for k, currency := range jsonCountry.Currencies {
		result.AddCurrency(k, currency.Name, currency.Symbol)
	}
	result.SetCapital(jsonCountry.Capital...)
	result.SetAltSpellings(jsonCountry.AltSpellings...)
	result.SetRegion(jsonCountry.Region)
	result.SetSubregion(jsonCountry.Subregion)
	for k, v := range jsonCountry.Languages {
		result.AddLanguage(k, v)
	}
	for k, v := range jsonCountry.Translations {
		result.AddTranslation(k, v.Common, v.Official)
	}
	//result.SetBorders(jsonCountry.Borders...)
	//result.SetCallingCodes(jsonCountry.CallingCodes.Root, jsonCountry.CallingCodes.Suffixes...)
	//result.SetCapitalInfo(jsonCountry.CapitalInfo.LatLng.Lat, jsonCountry.CapitalInfo.LatLng.Lng)
	//jsonCountry.Currencies.UnmarshalCurrencies(result)
	result.SetContinents(jsonCountry.Continents...)
	return result
}
