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
func (l *languagesCache) getLanguages() []*languageRecord {
	languages := make([]*languageRecord, 0, len(l.languages))
	for _, record := range l.languages {
		languages = append(languages, record)
	}
	return languages
}
