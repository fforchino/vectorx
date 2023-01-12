package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
)

/**********************************************************************************************************************/
/*                            VIDEO DOORPHONE INTERFACE USING VOCAL COMMANDS (TEST/DEMO)                              */
/**********************************************************************************************************************/

func Urmet_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"open the door"}
	utterances[LOCALE_ITALIAN] = []string{"apri la porta"}
	utterances[LOCALE_SPANISH] = []string{"abre la puerta"}
	utterances[LOCALE_FRENCH] = []string{"ouvre la porte"}
	utterances[LOCALE_GERMAN] = []string{"öffne die Tür"}

	var intent = IntentDef{
		IntentName: "extended_intent_urmet",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    urmetHandler,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_URMET_OPEN_THE_DOOR", []string{"open the door", "apri la porta", "abre la puerta", "ouvre la porte", "öffne die Tür"})
	addLocalizedString("STR_URMET_WAKEWORD", []string{"hey urmet", "hey urmet", "hey urmet", "hey urmet", "hey urmet"})

	return nil
}

func urmetHandler(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	sdk_wrapper.SayText(getText("STR_URMET_WAKEWORD"))
	return returnIntent
}
