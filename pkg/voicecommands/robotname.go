package voicecommands

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"vectorx/pkg/intents"
)

func SetRobotName(params intents.IntentParams) string {
	returnIntent := "intent_imperative_negative"
	if len(params.RobotName) > 0 {
		sdk_wrapper.SetRobotName(params.RobotName)
		sdk_wrapper.SayText("Ok. My name is " + params.RobotName)
		returnIntent = "intent_imperative_affirmative"
	}
	return returnIntent
}

func SayRobotName() string {
	returnIntent := "intent_imperative_affirmative"
	robotName := sdk_wrapper.GetRobotName()
	if len(robotName) > 0 {
		sdk_wrapper.SayText("My name is " + robotName)
	} else {
		sdk_wrapper.SayText("I am just Vector")
	}
	return returnIntent
}
