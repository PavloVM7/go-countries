package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountryExt_SetBorders(t *testing.T) {
	var cext countryExt
	expected := []string{"AND", "BEL", "DEU", "ITA", "LUX", "MCO", "ESP", "CHE"}
	cext.SetBorders(expected...)
	actual := cext.Borders()
	assert.Equal(t, expected, actual)
}
