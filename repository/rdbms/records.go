package rdbms

type CountryRecord struct {
	CountryId    uint16
	Alpha2Code   string
	Alpha3Code   string
	OlympicCode  string
	FifaCode     string
	Flag         string
	Population   int32
	Area         float32
	Independent  bool
	Landlocked   bool
	UnMember     bool
	Latitude     float32
	Longitude    float32
	RegionId     uint32
	SubregionId  uint32
	OfficialName string
	CommonName   string
	StartOfWeek  string
	Status       string
}

type RegionRecord struct {
	RegionId   uint32
	ParentId   uint32
	RegionName string
}

type LanguageRecord struct {
	LanguageId uint16
	Language   string
}

type TranslationRecord struct {
	Id           uint32
	CountryId    uint16
	LanguageId   uint16
	Native       bool
	OfficialName string
	CommonName   string
}
