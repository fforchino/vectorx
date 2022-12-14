package intents

import (
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
	} else if contains(intent.Parameters, PARAMETER_ROBOTNAME) {
		var username string
		var nameSplitter string = ""
		if strings.Contains(speechText, getText(STR_ROBOT_NAME_IS)) {
			nameSplitter = getText(STR_ROBOT_NAME_IS)
		} else if strings.Contains(speechText, getText(STR_ROBOT_NAME_IS2)) {
			nameSplitter = getText(STR_ROBOT_NAME_IS2)
		} else if strings.Contains(speechText, getText(STR_ROBOT_NAME_IS3)) {
			nameSplitter = getText(STR_ROBOT_NAME_IS3)
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
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon := weatherParser(speechText, BotLocation, BotUnits)
		wp := WeatherParams{condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit, icon}
		intentParams.Weather = wp
	} else if contains(intent.Parameters, PARAMETER_CHAT_TARGET) {
		var username string
		var nameSplitter string = ""
		if strings.Contains(speechText, getText(STR_SET_CHAT_TARGET)) {
			nameSplitter = getText(STR_SET_CHAT_TARGET)
		}
		if nameSplitter != "" {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			intentParams.ChatTargetName = username
		}
	}
	return intentParams
}
