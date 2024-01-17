package rdbms

import (
	"fmt"
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
			lc := &languagesCache{languages: make(map[string]*languageRecord)}
			lc.addLanguage(tt.langShortAndName1[0], tt.langShortAndName1[1])

			for _, lang := range tt.langShortAndName1 {
				l.addLanguage(lang, "")
			}
			for _, lang := range tt.langShortAndName2 {
				l.addLanguage(lang, "")
			}
			if l.languages["eng"] != tt.expected {
				t.Errorf("languagesCache.addLanguage() = %v, want %v", l.languages["eng"], tt.expected)
			}
		})
	}
}
