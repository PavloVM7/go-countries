package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopyStringArraySkipEmpty(t *testing.T) {
	source := []string{"\tstring 1  ", " string 2\n", "\t    \r\n", "string 3"}
	actual := CopyStringArraySkipEmpty(source)
	assert.Equal(t, []string{"string 1", "string 2", "string 3"}, actual)
}
func TestCopyStringArraySkipEmpty_same(t *testing.T) {
	source := []string{"string 1", "string 2", "string 3"}
	actual := CopyStringArraySkipEmpty(source)
	assert.Equal(t, source, actual)
}
func TestCopyStringArraySkipEmpty_empty(t *testing.T) {
	source := []string{" ", "", "\r\n", "\t"}
	actual := CopyStringArraySkipEmpty(source)
	assert.Equal(t, []string{}, actual)
}

func TestCopyArray_int(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	actual := CopyArray(expected)
	assert.Equal(t, expected, actual)
	assert.NotSame(t, expected, actual)
}

func TestCopyArray_empty(t *testing.T) {
	expected := make([]string, 0)
	actual := CopyArray(expected)
	assert.Equal(t, expected, actual)
	assert.NotSame(t, expected, actual)
}

func TestCopyArray_nil(t *testing.T) {
	var expected []interface{}
	actual := CopyArray(expected)
	assert.Nil(t, actual)
}
