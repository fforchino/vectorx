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
	utterances[LOCALE_SPANISH] = []string{"tú eres", "te llamas", "tu nombre es"}
	utterances[LOCALE_FRENCH] = []string{"tu es", "tu t'appelles", "ton nom est"}
	utterances[LOCALE_GERMAN] = []string{"Du bist", "du nennst dich", "dein Name ist"}

	var intent = IntentDef{
		IntentName: "extended_intent_set_robot_name",
		Utterances: utterances,
		Parameters: []string{PARAMETER_USERNAME},
		Handler:    setRobotName,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_ROBOT_GET_NAME", []string{"my name is %s1", "mi chiamo %s1", "mi nombre es %s1", "je m'appelle %s1", "mein name ist %s1"})
	addLocalizedString("STR_ROBOT_SET_NAME", []string{"ok. my name is %s1", "bene, mi chiamerò %s1", "bueno. mi nombre es %s1", "d'accord. mon nom est %s1", "ok. mein name ist %s1"})
	addLocalizedString("STR_ROBOT_NO_NAME", []string{"i don't have a name yet", "non ho ancora un nome", "todavía no tengo nombre", "je n'ai pas encore de nom", "ich habe noch keinen namen"})

	return nil
}

func setRobotName(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	if len(params.RobotName) > 0 {
		sdk_wrapper.SetRobotName(params.RobotName)
		sdk_wrapper.SayText(getTextEx("STR_ROBOT_SET_NAME", []string{params.RobotName}))
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
	utterances[LOCALE_SPANISH] = []string{"quién eres", "cuál es tu nombre"}
	utterances[LOCALE_FRENCH] = []string{"qui tu es", "quel est ton nom"}
	utterances[LOCALE_GERMAN] = []string{"Wer du bist", "Wie heißt du?"}

	var intent = IntentDef{
		IntentName: "extended_intent_say_robot_name",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    sayRobotName,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func sayRobotName(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	robotName := sdk_wrapper.GetRobotName()
	if len(robotName) > 0 {
		sdk_wrapper.SayText(getTextEx("STR_ROBOT_GET_NAME", []string{robotName}))
	} else {
		sdk_wrapper.SayText(getText("STR_ROBOT_NO_NAME"))
	}
	return returnIntent
}
