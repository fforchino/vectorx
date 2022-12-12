package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"strconv"
	"strings"
)

const STR_WEATHER_IN = "str_weather_in"
const STR_WEATHER_FORECAST = "str_weather_forecast"
const STR_WEATHER_TOMORROW = "str_weather_tomorrow"
const STR_WEATHER_THE_DAY_AFTER_TOMORROW = "str_weather_the_day_after_tomorrow"
const STR_WEATHER_TONIGHT = "str_weather_tonight"
const STR_WEATHER_THIS_AFTERNOON = "str_weather_this_afternoon"
const STR_WEATHER_AND = "str_weather_and"
const STR_NAME_IS = "str_name_is"
const STR_NAME_IS2 = "str_name_is1"
const STR_NAME_IS3 = "str_name_is2"
const STR_ROBOT_GET_NAME = "str_robot_get_name"
const STR_ROBOT_SET_NAME = "str_robot_set_name"
const STR_ROBOT_NO_NAME = "str_robot_no_name"
const STR_HELLO_WORLD = "str_hello_world"
const STR_LANGUAGE_ITALIAN = "str_language_italian"
const STR_LANGUAGE_FRENCH = "str_language_french"
const STR_LANGUAGE_SPANISH = "str_language_spanish"
const STR_LANGUAGE_GERMAN = "str_language_german"
const STR_LANGUAGE_ENGLISH = "str_language_english"

const en_US = 0
const it_IT = 1
const es_ES = 2
const fr_FR = 3
const de_DE = 4

// All text must be lowercase!

var texts = map[string][]string{
	//  key                 			en-US   it-IT   es-ES    fr-FR    de-DE
	STR_WEATHER_IN:                     {" in ", " a ", " en ", " en ", " in "},
	STR_WEATHER_FORECAST:               {"forecast", "previsioni", "pronóstico", "prévisions", "wettervorhersage"},
	STR_WEATHER_TOMORROW:               {"tomorrow", "domani", "mañana", "demain", "morgen"},
	STR_WEATHER_THE_DAY_AFTER_TOMORROW: {"day after tomorrow", "dopodomani", "el día después de mañana", "lendemain de demain", "am tag nach morgen"},
	STR_WEATHER_TONIGHT:                {"tonight", "stasera", "esta noche", "ce soir", "heute abend"},
	STR_WEATHER_THIS_AFTERNOON:         {"afternoon", "pomeriggio", "esta tarde", "après-midi", "heute nachmittag"},
	STR_NAME_IS:                        {" is ", " è ", " es ", " est ", " ist "},
	STR_NAME_IS2:                       {"'s", "sono ", "soy ", "suis ", "bin "},
	STR_NAME_IS3:                       {"names", " chiamo ", " llamo ", "appelle ", "werde"},
	STR_LANGUAGE_ITALIAN:               {"italian", "italiano", "italiano", "italien", "italienisch"},
	STR_LANGUAGE_SPANISH:               {"spanish", "spagnolo", "castellano", "espagnol", "spanisch "},
	STR_LANGUAGE_FRENCH:                {"french", "francese", "inglés", "français", "französisch"},
	STR_LANGUAGE_GERMAN:                {"german", "tedesco", "alemán", "allemand", "deutsch"},
	STR_LANGUAGE_ENGLISH:               {"english", "inglese", "inglés", "anglais", "englisch"},
	STR_ROBOT_GET_NAME:                 {"my name is %s1", "mi chiamo %s1", "mi nombre es %s1", "je m'appelle %s1", "mein name ist %s1"},
	STR_ROBOT_SET_NAME:                 {"ok. my name is %s1", "bene, mi chiamerò %s1", "bueno. mi nombre es %s1", "d'accord. mon nom est %s1", "ok. mein name ist %s1"},
	STR_ROBOT_NO_NAME:                  {"i don't have a name yet", "non ho ancora un nome", "todavía no tengo nombre", "je n'ai pas encore de nom", "ich habe noch keinen namen"},
	STR_HELLO_WORLD:                    {"hello world!", "ciao mondo!", "hola mundo!", "bonjour le monde!", "hallo welt"},
	STR_WEATHER_AND:                    {" degrees and ", " gradi e ", " grados y ", " degrés et ", " Grad und "},
}

func getText(key string) string {
	return getTextEx(key, []string{})
}

func getTextEx(key string, params []string) string {
	text := key
	var data = texts[key]
	if data != nil {
		if sdk_wrapper.GetLanguage() == sdk_wrapper.LANGUAGE_ITALIAN {
			text = data[it_IT]
		} else if sdk_wrapper.GetLanguage() == sdk_wrapper.LANGUAGE_SPANISH {
			text = data[es_ES]
		} else if sdk_wrapper.GetLanguage() == sdk_wrapper.LANGUAGE_FRENCH {
			text = data[fr_FR]
		} else if sdk_wrapper.GetLanguage() == sdk_wrapper.LANGUAGE_GERMAN {
			text = data[de_DE]
		} else {
			text = data[en_US]
		}
	}
	for i, p := range params {
		text = strings.Replace(text, "%s"+strconv.Itoa(i+1), p, 1)
	}
	return text
}
