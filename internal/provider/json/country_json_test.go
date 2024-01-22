package json

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"pm.com/go-countries/domain"
	"testing"
)

func Test_createCountryFromJson(t *testing.T) {
	name := "../../../tests/testdata/france.json"
	bytes := testFileToBytes(name)
	assert.NotEqual(t, 0, len(bytes))
	var country CountryJSON
	err := json.Unmarshal(bytes, &country)
	if err != nil {
		t.Fatal(err)
	}
	actual := createCountryFromJson(&country)
	expected := domain.NewCountry(250, "FR", "FRA")
	expected.SetName("France", "French Republic")
	expected.AddNativeName("fra", "France", "République française")
	expected.SetTopLevelDomains(".fr")
	expected.SetOlympicCode("FRA")
	expected.SetIndependent(true)
	expected.SetStatus("officially-assigned")
	expected.SetUnMember(true)
	expected.AddCurrency("EUR", "Euro", "€")
	expected.SetCapital("Paris")
	expected.SetAltSpellings("FR", "French Republic", "République française")
	expected.SetRegion("Europe")
	expected.SetSubregion("Western Europe")
	expected.SetContinents("Europe")

	expected.AddLanguage("fra", "French")
	expected.AddTranslation("ces", "Francie", "Francouzská republika")
	expected.AddTranslation("ara", "فرنسا", "الجمهورية الفرنسية")
	expected.AddTranslation("bre", "Frañs", "Republik Frañs")
	expected.AddTranslation("cym", "France", "French Republic")
	expected.AddTranslation("deu", "Frankreich", "Französische Republik")
	expected.AddTranslation("est", "Prantsusmaa", "Prantsuse Vabariik")
	expected.AddTranslation("fin", "Ranska", "Ranskan tasavalta")
	expected.AddTranslation("fra", "France", "République française")
	expected.AddTranslation("hrv", "Francuska", "Francuska Republika")
	expected.AddTranslation("hun", "Franciaország", "Francia Köztársaság")
	expected.AddTranslation("ita", "Francia", "Repubblica francese")
	expected.AddTranslation("jpn", "フランス", "フランス共和国")
	expected.AddTranslation("kor", "프랑스", "프랑스 공화국")
	expected.AddTranslation("nld", "Frankrijk", "Franse Republiek")
	expected.AddTranslation("per", "فرانسه", "جمهوری فرانسه")
	expected.AddTranslation("pol", "Francja", "Republika Francuska")
	expected.AddTranslation("por", "França", "República Francesa")
	expected.AddTranslation("rus", "Франция", "Французская Республика")
	expected.AddTranslation("slk", "Francúzsko", "Francúzska republika")
	expected.AddTranslation("spa", "Francia", "República francés")
	expected.AddTranslation("srp", "Француска", "Француска Република")
	expected.AddTranslation("swe", "Frankrike", "Republiken Frankrike")
	expected.AddTranslation("tur", "Fransa", "Fransa Cumhuriyeti")
	expected.AddTranslation("urd", "فرانس", "جمہوریہ فرانس")
	expected.AddTranslation("zho", "法国", "法兰西共和国")

	expected.SetLatLng(46.0, 2.0)
	expected.SetLandlocked(false)
	expected.SetBorders("AND", "BEL", "DEU", "ITA", "LUX", "MCO", "ESP", "CHE")
	expected.SetArea(551695.0)

	expected.AddDemonym("eng", "French", "French")
	expected.AddDemonym("fra", "Française", "Français")

	expected.SetFlag("🇫🇷")

	expected.AddMap("googleMaps", "https://goo.gl/maps/g7QxxSFsWyTPKuzd7")
	expected.AddMap("openStreetMaps", "https://www.openstreetmap.org/relation/1403916")

	expected.SetPopulation(67391582)
	expected.SetFifa("FRA")
	expected.SetCar("right", "F")
	expected.SetTimezones("UTC-10:00", "UTC-09:30", "UTC-09:00", "UTC-08:00", "UTC-04:00", "UTC-03:00",
		"UTC+01:00", "UTC+02:00", "UTC+03:00", "UTC+04:00", "UTC+05:00", "UTC+10:00", "UTC+11:00", "UTC+12:00")
	expected.AddFlag("png", "https://flagcdn.com/w320/fr.png")
	expected.AddFlag("svg", "https://flagcdn.com/fr.svg")
	expected.AddFlag("alt", "The flag of France is composed of three equal vertical bands of blue, white and red.")

	expected.AddCoatOfArm("png", "https://mainfacts.com/media/images/coats_of_arms/fr.png")
	expected.AddCoatOfArm("svg", "https://mainfacts.com/media/images/coats_of_arms/fr.svg")

	expected.SetStartOfWeek("monday")
	expected.SetCapitalInfo(domain.LatLng{Lat: 48.87, Lng: 2.33})

	assert.Equal(t, expected, actual)
}

func TestCountryJSON_unmarshal(t *testing.T) {
	var country CountryJSON
	err := json.Unmarshal(getTestCountryString(), &country)
	if err != nil {
		t.Fatal(err)
	}
	expected := CountryJSON{
		LatLng:      &LatLngJson{Lat: -10.5, Lng: 105.66666666},
		CapitalInfo: &CapitalInfoJson{LatLng: &LatLngJson{Lat: -10.42, Lng: 105.68}},
		Population:  2072,
		Area:        135.0,
		NumericCode: 162,
		Independent: false,
		UnMember:    false,
		Landlocked:  false,
		Alpha2Code:  "CX",
		Alpha3Code:  "CXR",
		Name: CountryNameJson{
			NameJson: NameJson{Official: "Territory of Christmas Island", Common: "Christmas Island"},
			NativeName: map[string]NameJson{
				"eng": {Official: "Territory of Christmas Island", Common: "Christmas Island"},
			},
		},
		Currencies: map[string]CurrencyJson{
			"AUD": {Name: "Australian dollar", Symbol: "$"},
		},
		CallingCodes: CollingCodesJson{Root: "+6", Suffixes: []string{"1"}},
		Region:       "Oceania",
		Subregion:    "Australia and New Zealand",
		Status:       "officially-assigned",
		StartOfWeek:  "monday",
		Flag:         "🇨🇽",
		Tld:          []string{".cx"},
		Capital:      []string{"Flying Fish Cove"},
		AltSpellings: []string{"CX", "Territory of Christmas Island"},
		Languages:    map[string]string{"eng": "English"},
		Continents:   []string{"Asia"},
		Timezones:    []string{"UTC+07:00"},
		Translations: map[string]NameJson{
			"ara": {Official: "جزيرة كريسماس", Common: "جزيرة كريسماس"},
			"bre": {Official: "Tiriad Enez Christmas", Common: "Enez Christmas"},
			"ces": {Official: "Teritorium Vánočního ostrova", Common: "Vánoční ostrov"},
			"cym": {Official: "Tiriogaeth yr Ynys y Nadolig", Common: "Ynys y Nadolig"},
			"deu": {Official: "Gebiet der Weihnachtsinsel", Common: "Weihnachtsinsel"},
			"est": {Official: "Jõulusaare ala", Common: "Jõulusaar"},
			"fin": {Official: "Joulusaaren alue", Common: "Joulusaari"},
			"fra": {Official: "Territoire de l'île Christmas", Common: "Île Christmas"},
			"hrv": {Official: "Teritorij Božićni otok", Common: "Božićni otok"},
			"hun": {Official: "Karácsony-sziget", Common: "Karácsony-sziget"},
			"ita": {Official: "Territorio di Christmas Island", Common: "Isola di Natale"},
			"jpn": {Official: "クリスマス島の領土", Common: "クリスマス島"},
			"kor": {Official: "크리스마스 섬", Common: "크리스마스 섬"},
			"nld": {Official: "Grondgebied van Christmas Island", Common: "Christmaseiland"},
			"per": {Official: "جزیرهٔ کریسمس", Common: "جزیرهٔ کریسمس"},
			"pol": {Official: "Wyspa Bożego Narodzenia", Common: "Wyspa Bożego Narodzenia"},
			"por": {Official: "Território da Ilha Christmas", Common: "Ilha do Natal"},
			"rus": {Official: "Территория острова Рождества", Common: "Остров Рождества"},
			"slk": {Official: "Teritórium Vianočného ostrova", Common: "Vianočnú ostrov"},
			"spa": {Official: "Territorio de la Isla de Navidad", Common: "Isla de Navidad"},
			"srp": {Official: "Божићно Острво", Common: "Божићно Острво"},
			"swe": {Official: "Julön", Common: "Julön"},
			"tur": {Official: "Christmas Adası", Common: "Christmas Adası"},
			"urd": {Official: "ریاستِ جزیرہ کرسمس", Common: "جزیرہ کرسمس"},
			"zho": {Official: "圣诞岛", Common: "圣诞岛"},
		},
		Demonyms: map[string]Demonyms{
			"eng": {F: "Christmas Islander", M: "Christmas Islander"},
		},
		Maps: map[string]string{
			"googleMaps":     "https://goo.gl/maps/ZC17hHsQZpShN5wk9",
			"openStreetMaps": "https://www.openstreetmap.org/relation/6365444",
		},
		CoatOfArms: map[string]string{
			"png": "https://mainfacts.com/media/images/coats_of_arms/cx.png",
			"svg": "https://mainfacts.com/media/images/coats_of_arms/cx.svg",
		},
		Flags: map[string]string{
			"png": "https://flagcdn.com/w320/cx.png",
			"svg": "https://flagcdn.com/cx.svg",
		},
		Car: CarJson{Signs: []string{"AUS"}, Side: "left"},
	}
	assert.Equal(t, expected, country)
}

func getTestCountryString() []byte {
	s := `
{
  "name":{
    "common":"Christmas Island",
    "official":"Territory of Christmas Island",
    "nativeName":{
      "eng":{
        "official":"Territory of Christmas Island",
        "common":"Christmas Island"
      }
    }
  },
  "tld":[".cx"],
  "cca2":"CX",
  "ccn3":"162",
  "cca3":"CXR",
  "independent":false,
  "status":"officially-assigned",
  "unMember":false,
  "currencies":{
    "AUD":{"name":"Australian dollar","symbol":"$"}
  },
  "idd":{"root":"+6","suffixes":["1"]},
  "capital":["Flying Fish Cove"],
  "altSpellings":["CX","Territory of Christmas Island"],
  "region":"Oceania",
  "subregion":"Australia and New Zealand",
  "languages":{"eng":"English"},
  "translations":{
    "ara":{"official":"جزيرة كريسماس","common":"جزيرة كريسماس"},
    "bre":{"official":"Tiriad Enez Christmas","common":"Enez Christmas"},
    "ces":{"official":"Teritorium Vánočního ostrova","common":"Vánoční ostrov"},
    "cym":{"official":"Tiriogaeth yr Ynys y Nadolig","common":"Ynys y Nadolig"},
    "deu":{"official":"Gebiet der Weihnachtsinsel","common":"Weihnachtsinsel"},
    "est":{"official":"Jõulusaare ala","common":"Jõulusaar"},
    "fin":{"official":"Joulusaaren alue","common":"Joulusaari"},
    "fra":{"official":"Territoire de l'île Christmas","common":"Île Christmas"},
    "hrv":{"official":"Teritorij Božićni otok","common":"Božićni otok"},
    "hun":{"official":"Karácsony-sziget","common":"Karácsony-sziget"},
    "ita":{"official":"Territorio di Christmas Island","common":"Isola di Natale"},
    "jpn":{"official":"クリスマス島の領土","common":"クリスマス島"},
    "kor":{"official":"크리스마스 섬","common":"크리스마스 섬"},
    "nld":{"official":"Grondgebied van Christmas Island","common":"Christmaseiland"},
    "per":{"official":"جزیرهٔ کریسمس","common":"جزیرهٔ کریسمس"},
    "pol":{"official":"Wyspa Bożego Narodzenia","common":"Wyspa Bożego Narodzenia"},
    "por":{"official":"Território da Ilha Christmas","common":"Ilha do Natal"},
    "rus":{"official":"Территория острова Рождества","common":"Остров Рождества"},
    "slk":{"official":"Teritórium Vianočného ostrova","common":"Vianočnú ostrov"},
    "spa":{"official":"Territorio de la Isla de Navidad","common":"Isla de Navidad"},
    "srp":{"official":"Божићно Острво","common":"Божићно Острво"},
    "swe":{"official":"Julön","common":"Julön"},
    "tur":{"official":"Christmas Adası","common":"Christmas Adası"},
    "urd":{"official":"ریاستِ جزیرہ کرسمس","common":"جزیرہ کرسمس"},
    "zho":{"official":"圣诞岛","common":"圣诞岛"}
  },
  "latlng":[-10.5,105.66666666],
  "landlocked":false,
  "area":135.0,
  "demonyms":{
    "eng":{"f":"Christmas Islander","m":"Christmas Islander"}
  },
  "flag":"\uD83C\uDDE8\uD83C\uDDFD",
  "maps":{
    "googleMaps":"https://goo.gl/maps/ZC17hHsQZpShN5wk9",
    "openStreetMaps":"https://www.openstreetmap.org/relation/6365444"
  },
  "population":2072,
  "car":{"signs":["AUS"],"side":"left"},
  "timezones":["UTC+07:00"],
  "continents":["Asia"],
  "flags":{
    "png":"https://flagcdn.com/w320/cx.png",
    "svg":"https://flagcdn.com/cx.svg"
  },
  "coatOfArms":{
    "png":"https://mainfacts.com/media/images/coats_of_arms/cx.png",
    "svg":"https://mainfacts.com/media/images/coats_of_arms/cx.svg"
  },
  "startOfWeek":"monday",
  "capitalInfo":{"latlng":[-10.42,105.68]},
  "postalCode":{"format":"####","regex":"^(\\d{4})$"}
}
`
	return []byte(s)
}

func TestLngNamesJson_UnmarshalJSON(t *testing.T) {
	src := `
{
	"translations": {
		"eng":{"official":"Territory of Christmas Island","common":"Territory of Christmas Island"},
		"ara":{"official":"جزيرة كريسماس","common":"جزيرة كريسماس"},
	    "fra":{"official":"Territoire de l'île Christmas","common":"Île Christmas"},
	    "zho":{"official":"圣诞岛","common":"圣诞岛"}
	}
}
`
	type ReqTest struct {
		Translations NamesJson `json:"translations"`
	}
	var res ReqTest
	require.NoError(t, json.Unmarshal([]byte(src), &res))
	assert.Equal(t, ReqTest{
		Translations: map[string]NameJson{
			"eng": {Official: "Territory of Christmas Island", Common: "Territory of Christmas Island"},
			"ara": {Official: "جزيرة كريسماس", Common: "جزيرة كريسماس"},
			"fra": {Official: "Territoire de l'île Christmas", Common: "Île Christmas"},
			"zho": {Official: "圣诞岛", Common: "圣诞岛"},
		},
	}, res)
}

func TestLatLngJson_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		src      string
		expected *LatLngJson
	}{
		{
			src:      `{"latlng":[-10.5,105.66666666]}`,
			expected: &LatLngJson{Lat: -10.5, Lng: 105.66666666},
		},
		{
			src:      `{"latlng":[-10.42,105.68]}`,
			expected: &LatLngJson{Lat: -10.42, Lng: 105.68},
		},
		{
			src:      `{"latlng":[]}`,
			expected: &LatLngJson{},
		},
	}
	type ReqTst struct {
		Latlng *LatLngJson `json:"latlng"`
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			var actual ReqTst
			require.NoError(t, json.Unmarshal([]byte(tt.src), &actual))
			assert.Equal(t, ReqTst{Latlng: tt.expected}, actual)
		})
	}
}

func TestCurrenciesJson_UnmarshalJSON(t *testing.T) {
	src := `
{
	"currencies":{
		"AUD":{"name":"Australian dollar","symbol":"$"},
		"WST":{"name":"Samoan tālā","symbol":"T"},
		"DJF":{"name":"Djiboutian franc","symbol":"Fr"}
	}
}
`
	type ReqTest struct {
		Currencies CurrenciesJson `json:"currencies"`
	}
	var res ReqTest
	require.NoError(t, json.Unmarshal([]byte(src), &res))
	assert.Equal(t, ReqTest{
		Currencies: map[string]CurrencyJson{
			"AUD": {Name: "Australian dollar", Symbol: "$"},
			"WST": {Name: "Samoan tālā", Symbol: "T"},
			"DJF": {Name: "Djiboutian franc", Symbol: "Fr"},
		},
	}, res)
}

func testFileToBytes(fileName string) []byte {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return bytes
}
