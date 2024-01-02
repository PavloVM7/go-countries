package json

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"os"
	"testing"
)

func Test_parseJsonData(t *testing.T) {
	reader := createTestReader()
	actual, err := parseJsonData(reader)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 250, len(actual))
	maxNameLength := 0
	var name NameJson
	var borders []string
	var capital []string
	var capitalLatLng LatLngJson
	for _, c := range actual {
		nameLength := int(math.Max(float64(len(c.Name.Official)), float64(len(c.Name.Common))))
		if nameLength > maxNameLength {
			maxNameLength = nameLength
			name = c.Name.NameJson
		}
		if len(c.Borders) > len(borders) {
			borders = c.Borders
		}
		if len(c.Capital) > len(capital) {
			capital = c.Capital
			capitalLatLng = *c.LatLng
		}
		t.Log("continents:", c.Continents, ", region:", c.Region, ", subregion:", c.Subregion)
	}
	t.Log("maxNameLength =", maxNameLength, name)
	t.Log("maxBorders =", borders, ",", len(borders))
	t.Log("maxCapitals:", capital, "LatLng:", capitalLatLng)
}

func createTestReader() io.Reader {
	name := "../../../tests/testdata/all-countries.json"
	bts, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bts)
}
