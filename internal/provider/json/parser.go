package json

import (
	"encoding/json"
	"io"
)

func parseJsonData(r io.Reader) ([]CountryJSON, error) {
	decoder := json.NewDecoder(r)
	var res []CountryJSON
	if err := decoder.Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
