package rdbms

import "pm.com/go-countries/domain"

type CountryRecord struct {
	Population   uint32
	RegionId     uint32
	SubregionId  uint32
	Area         float32
	Latitude     float32
	Longitude    float32
	CountryId    uint16
	Independent  bool
	Landlocked   bool
	UnMember     bool
	Alpha2Code   string
	Alpha3Code   string
	OlympicCode  string
	FifaCode     string
	Flag         string
	OfficialName string
	CommonName   string
	StartOfWeek  string
	Status       string
}

func newCountryRecord(country *domain.Country) CountryRecord {
	result := CountryRecord{
		CountryId:    country.NumericCode(),
		Alpha2Code:   country.Alpha2Code(),
		Alpha3Code:   country.Alpha3Code(),
		OlympicCode:  country.OlympicCode(),
		FifaCode:     country.Fifa(),
		CommonName:   country.CommonName(),
		OfficialName: country.OfficialName(),
		Flag:         country.Flag(),
		Area:         country.Area(),
		Population:   country.Population(),
		Independent:  country.IsIndependent(),
		Landlocked:   country.IsLandlocked(),
		UnMember:     country.IsUnMember(),
		StartOfWeek:  country.StartOfWeek(),
		Status:       country.Status(),
	}
	result.Latitude, result.Longitude = country.LatLng()
	return result
}

func toCountryRecord(scn scannable, result *CountryRecord) error {
	return scn.Scan(&result.CountryId, &result.Alpha2Code, &result.Alpha3Code, &result.OlympicCode, &result.FifaCode,
		&result.Flag, &result.Population, &result.Area, &result.Independent, &result.Landlocked, &result.UnMember,
		&result.Latitude, &result.Longitude, &result.RegionId, &result.SubregionId, &result.OfficialName, &result.CommonName,
		&result.StartOfWeek, &result.Status)
}
