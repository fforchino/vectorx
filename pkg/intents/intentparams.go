package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"strings"
)

func ParseParams(speechText string, intent IntentDef) IntentParams {
	var intentParams IntentParams
	if contains(intent.Parameters, PARAMETER_USERNAME) {
		var username string
		var nameSplitter string = ""
		if strings.Contains(speechText, getText(STR_NAME_IS)) {
			nameSplitter = getText(STR_NAME_IS)
		} else if strings.Contains(speechText, getText(STR_NAME_IS2)) {
			nameSplitter = getText(STR_NAME_IS2)
		} else if strings.Contains(speechText, getText(STR_NAME_IS3)) {
			nameSplitter = getText(STR_NAME_IS3)
		}
		if nameSplitter != "" {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				username = username + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParams.RobotName = username
		}
	} else if contains(intent.Parameters, PARAMETER_LANGUAGE) {
		if strings.Contains(speechText, getText(STR_LANGUAGE_ITALIAN)) {
			intentParams.Language = LOCALE_ITALIAN
		} else if strings.Contains(speechText, getText(STR_LANGUAGE_SPANISH)) {
			intentParams.Language = LOCALE_SPANISH
		} else if strings.Contains(speechText, getText(STR_LANGUAGE_FRENCH)) {
			intentParams.Language = LOCALE_FRENCH
		} else if strings.Contains(speechText, getText(STR_LANGUAGE_GERMAN)) {
			intentParams.Language = LOCALE_GERMAN
		} else if strings.Contains(speechText, getText(STR_LANGUAGE_ENGLISH)) {
			intentParams.Language = LOCALE_ENGLISH
		}
	} else if contains(intent.Parameters, PARAMETER_WEATHER) {
		botLocation := sdk_wrapper.GetVectorSettings()["default_location"].(string)
		botUnits := os.Getenv("WEATHERAPI_UNIT")
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon := weatherParser(speechText, botLocation, botUnits)
		wp := WeatherParams{condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon}
		intentParams.Weather = wp
	}
	return intentParams
}
