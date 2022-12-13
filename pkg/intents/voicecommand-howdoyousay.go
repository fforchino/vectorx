package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"strings"
)

/**********************************************************************************************************************/
/*                                                IMAGE TEST                                                          */
/**********************************************************************************************************************/

func HowDoYouSay_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"how do you say"}
	utterances[LOCALE_ITALIAN] = []string{"come si dice"}
	utterances[LOCALE_SPANISH] = []string{"como dicen"}
	utterances[LOCALE_FRENCH] = []string{"comme on dit"}
	utterances[LOCALE_GERMAN] = []string{"wie sie"}

	var intent = IntentDef{
		IntentName: "extended_intent_how_do_you_say",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    HowDoYouSay,
	}
	*intentList = append(*intentList, intent)

	addLocalizedString("STR_HOWDOYOUSAY_HOW_DO_YOU_SAY", []string{"how do you say ", "come si dice ", "como dicen ", "comme on dit ", "wie sie "})
	addLocalizedString("STR_HOWDOYOUSAY_IN", []string{" in ", " a ", " en ", " en ", " auf "})

	return nil
}

func HowDoYouSay(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	// Extract word and target language from the sentence
	var word string
	var language string
	wordSplitter := getText("STR_HOWDOYOUSAY_HOW_DO_YOU_SAY")
	languageSplitter := getText("STR_HOWDOYOUSAY_IN")

	if strings.Contains(speechText, wordSplitter) && strings.Contains(speechText, languageSplitter) {
		splitPhrase := strings.SplitAfter(speechText, wordSplitter)
		tmp := strings.TrimSpace(splitPhrase[1])
		splitPhrase2 := strings.SplitAfter(tmp, languageSplitter)
		word = splitPhrase2[0]
		language = splitPhrase2[1]
		trans := getText(STR_GENERIC_I_DONT_KNOW)
		currentLanguage := sdk_wrapper.GetLanguage()
		if word != "" {
			if language == getText(STR_LANGUAGE_ENGLISH) {
				sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_ENGLISH)
				trans = sdk_wrapper.Translate(word, currentLanguage, language)
			} else if language == getText(STR_LANGUAGE_ITALIAN) {
				sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_ITALIAN)
				trans = sdk_wrapper.Translate(word, currentLanguage, language)
			} else if language == getText(STR_LANGUAGE_SPANISH) {
				sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_SPANISH)
				trans = sdk_wrapper.Translate(word, currentLanguage, language)
			} else if language == getText(STR_LANGUAGE_FRENCH) {
				sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_FRENCH)
				trans = sdk_wrapper.Translate(word, currentLanguage, language)
			} else if language == getText(STR_LANGUAGE_GERMAN) {
				sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_GERMAN)
				trans = sdk_wrapper.Translate(word, currentLanguage, language)
			}
		}
		sdk_wrapper.SayText(trans)
		sdk_wrapper.SetLanguage(currentLanguage)
	}

	return returnIntent
}
