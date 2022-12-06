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
const STR_ROBOT_GET_NAME = "str_robot_get_name"
const STR_ROBOT_SET_NAME = "str_robot_set_name"
const STR_ROBOT_NO_NAME = "str_robot_no_name"

const en_US = 0
const it_IT = 1

var texts = map[string][]string{
	//  key                 			en-US   it-IT
	STR_WEATHER_IN:                     {" in ", " a "},
	STR_WEATHER_FORECAST:               {"forecast", "previsioni"},
	STR_WEATHER_TOMORROW:               {"tomorrow", "domani"},
	STR_WEATHER_THE_DAY_AFTER_TOMORROW: {"day after tomorrow", "dopodomani"},
	STR_WEATHER_TONIGHT:                {"tonight", "stasera"},
	STR_WEATHER_THIS_AFTERNOON:         {"afternoon", "pomeriggio"},
	STR_ROBOT_GET_NAME:                 {"my name is %s1", "mi chiamo %s1"},
	STR_ROBOT_SET_NAME:                 {"ok. my name is %s1", "bene, mi chiamerò %s1"},
	STR_ROBOT_NO_NAME:                  {"i don't have a name yet'", "non ho ancora un nome"},
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
		} else {
			text = data[en_US]
		}
	}
	for i, p := range params {
		text = strings.Replace(text, "%s"+strconv.Itoa(i+1), p, 1)
	}
	return text
}
