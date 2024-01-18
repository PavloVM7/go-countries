package rdbms

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_languagesCache_addLanguage(t *testing.T) {
	tests := []struct {
		name              string
		langShortAndName1 []string
		langShortAndName2 []string
		expected          *languageRecord
	}{
		{
			name:              "add one language with full name",
			langShortAndName1: []string{"eng", "English"},
			langShortAndName2: nil,
			expected:          &languageRecord{languageId: 0, language: "eng", languageName: "English"},
		},
		{
			name:              "add first language has full name",
			langShortAndName1: []string{"eng", "English"},
			langShortAndName2: []string{"eng", ""},
			expected:          &languageRecord{languageId: 0, language: "eng", languageName: "English"},
		},
		{
			name:              "add second language has full name",
			langShortAndName1: []string{"eng", ""},
			langShortAndName2: []string{"eng", "English"},
			expected:          &languageRecord{languageId: 0, language: "eng", languageName: "English"},
		},
	}
	for i, tt := range tests {
		tt.name = fmt.Sprintf("%d_%s", i, tt.name)
		t.Run(tt.name, func(t *testing.T) {
			lc := newLanguagesCache()
			lc.addLanguage(tt.langShortAndName1[0], tt.langShortAndName1[1])
			if len(tt.langShortAndName2) > 0 {
				lc.addLanguage(tt.langShortAndName2[0], tt.langShortAndName2[1])
			}
			assert.Equal(t, tt.expected, lc.languages["eng"])
			assert.Equal(t, []*languageRecord{tt.expected}, lc.getLanguages())
		})
	}
}
