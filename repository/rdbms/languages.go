package rdbms

type languagesCache struct {
	languages map[string]*languageRecord
}

func newLanguagesCache() *languagesCache {
	return &languagesCache{
		languages: make(map[string]*languageRecord, 128),
	}
}
func (l *languagesCache) addLanguage(language, languageName string) {
	record := l.languages[language]
	if record == nil {
		l.languages[language] = &languageRecord{languageId: 0, language: language, languageName: languageName}
		return
	}
	if record.languageName == "" && languageName != "" {
		record.languageName = languageName
	}
}
func (l *languagesCache) getLanguage(language string) *languageRecord {
	return l.languages[language]
}
