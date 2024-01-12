package json

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"os"
	"strings"
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
	maxLanguageNameLen := 0
	maxLenLanguage := ""
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
		if len(c.Capital) != 1 {
			t.Log(c.Name.Common, ":", len(c.Capital), ", capitals:", strings.Join(c.Capital, ","), ", LatLng:", c.LatLng)
		}
		for k, v := range c.Languages {
			if len(v) > maxLanguageNameLen {
				maxLanguageNameLen = len(v)
				maxLenLanguage = fmt.Sprintf("%s:'%s'", k, v)
			}
		}
		if len(c.Continents) != 1 {
			t.Log(len(c.Continents), ", continents:", c.Continents, ", region:", c.Region, ", subregion:", c.Subregion, ", name:", c.Name.Common)
		}
	}
	t.Log("maxNameLength =", maxNameLength, name)
	t.Log("maxBorders =", borders, ",", len(borders))
	t.Log("maxCapitals:", capital, "LatLng:", capitalLatLng)
	t.Log("maxLanguageNameLen:", maxLanguageNameLen, maxLenLanguage)
}

func createTestReader() io.Reader {
	name := "../../../tests/testdata/all-countries.json"
	bts, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bts)
}
