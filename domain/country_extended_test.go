package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryExt_SetContinents(t *testing.T) {
	var cext countryExt
	continents := []string{"Europe", "Asia", ""}
	cext.SetContinents(continents...)
	actual := cext.Continents()
	assert.Equal(t, []string{"Europe", "Asia"}, actual)
}

func TestCountryExt_SetBorders(t *testing.T) {
	var cext countryExt
	expected := []string{"AND", "BEL", "DEU", "ITA", "LUX", "MCO", "ESP", "CHE"}
	cext.SetBorders(expected...)
	actual := cext.Borders()
	assert.Equal(t, expected, actual)
}
