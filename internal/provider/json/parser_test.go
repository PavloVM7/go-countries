package json

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
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
	for i, c := range actual {
		if len(c.Capital) > 1 {
			t.Logf("%d. %+v", i, c)
		}
	}
}

func createTestReader() io.Reader {
	name := "../../../tests/testdata/all-countries.json"
	bts, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bts)
}
