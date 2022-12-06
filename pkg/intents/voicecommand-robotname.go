package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
)

func RobotName_Register(intentList *[]IntentDef) error {
	registerSetRobotName(intentList)
	registerSayRobotName(intentList)
	return nil
}

/**********************************************************************************************************************/
/*                                            SET ROBOT NAME                                                          */
/**********************************************************************************************************************/

func registerSetRobotName(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"you are", "your name is", "your name's"}
	utterances[LOCALE_ITALIAN] = []string{"tu sei", "ti chiami", "il tuo nome è"}

	var intent = IntentDef{
		IntentName: "extended_intent_set_robot_name",
		Utterances: utterances,
		Parameters: []string{PARAMETER_USERNAME},
		Handler:    setRobotName,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func setRobotName(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	if len(params.RobotName) > 0 {
		sdk_wrapper.SetRobotName(params.RobotName)
		sdk_wrapper.SayText(getTextEx(STR_ROBOT_SET_NAME, []string{params.RobotName}))
		returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                            SAY ROBOT NAME                                                          */
/**********************************************************************************************************************/

func registerSayRobotName(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"who are you", "what's your name", "what is your name"}
	utterances[LOCALE_ITALIAN] = []string{"chi sei", "come ti chiami"}

	var intent = IntentDef{
		IntentName: "extended_intent_say_robot_name",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    sayRobotName,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func sayRobotName(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	robotName := sdk_wrapper.GetRobotName()
	if len(robotName) > 0 {
		sdk_wrapper.SayText(getTextEx(STR_ROBOT_GET_NAME, []string{robotName}))
	} else {
		sdk_wrapper.SayText(getText(STR_ROBOT_NO_NAME))
	}
	return returnIntent
}
