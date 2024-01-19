package json

import (
	"encoding/json"
	"io"
	"pm.com/go-countries/domain"
)

type CountriesJsonReader struct {
	r io.Reader
}

func (cjr *CountriesJsonReader) Read() ([]domain.Country, error) {
	countries, err := parseJsonData(cjr.r)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Country, 0, len(countries))
	for _, jsonCountry := range countries {
		result = append(result, createCountryFromJson(&jsonCountry))
	}
	return result, nil
}

func parseJsonData(r io.Reader) ([]CountryJSON, error) {
	decoder := json.NewDecoder(r)
	var res []CountryJSON
	if err := decoder.Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
