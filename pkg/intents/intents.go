package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"strings"
)

type IntentParams struct {
	RobotName string
}

func ParseParams(speechText string, intent string) IntentParams {
	robotLocale := sdk_wrapper.GetLocale()
	if strings.HasPrefix(robotLocale, "it_") {
		return ParseParamsIt(speechText, intent)
	}
	return ParseParamsEn(speechText, intent)
}
