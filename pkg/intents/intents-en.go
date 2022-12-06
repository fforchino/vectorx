package intents

import "strings"

func ParseParamsEn(speechText string, intent IntentDef) IntentParams {
	var intentParams IntentParams
	if contains(intent.Parameters, PARAMETER_USERNAME) {
		var username string
		var nameSplitter string
		if strings.Contains(speechText, "is") {
			nameSplitter = "is"
		} else if strings.Contains(speechText, "'s") {
			nameSplitter = "'s"
		} else if strings.Contains(speechText, "names") {
			nameSplitter = "names"
		}
		if strings.Contains(speechText, "is") || strings.Contains(speechText, "'s") || strings.Contains(speechText, "names") {
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
	}
	return intentParams
}
