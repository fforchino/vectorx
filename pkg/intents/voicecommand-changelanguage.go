package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"os/exec"
	"path"
	"strings"
)

/**********************************************************************************************************************/
/*                                                CHANGE LANGUAGE                                                     */
/**********************************************************************************************************************/

func ChangeLanguage_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"let's talk", "let's speak"}
	utterances[LOCALE_ITALIAN] = []string{"parliamo"}
	utterances[LOCALE_SPANISH] = []string{"hablamos"}
	utterances[LOCALE_FRENCH] = []string{"parlons"}
	utterances[LOCALE_GERMAN] = []string{"sprechen"}

	var intent = IntentDef{
		IntentName: "extended_intent_changelanguage",
		Utterances: utterances,
		Parameters: []string{PARAMETER_LANGUAGE},
		Handler:    changeLanguage,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func changeLanguage(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	loc := "en-US"
	switch params.Language {
	case LOCALE_ITALIAN:
		sdk_wrapper.SetLocale("it-IT")
		loc = "it-IT"
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_ITALIAN)
		break
	case LOCALE_SPANISH:
		sdk_wrapper.SetLocale("es-ES")
		loc = "es-ES"
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_SPANISH)
		break
	case LOCALE_FRENCH:
		sdk_wrapper.SetLocale("fr-FR")
		loc = "fr-FR"
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_FRENCH)
		break
	case LOCALE_GERMAN:
		sdk_wrapper.SetLocale("de-DE")
		loc = "de-DE"
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_GERMAN)
		break
	case LOCALE_ENGLISH:
		sdk_wrapper.SetLocale("en-US")
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_ENGLISH)
		break
	default:
		return STANDARD_INTENT_SYSTEM_NOAUDIO
	}
	newLanguage := strings.Split(loc, "-")[0]
	sdk_wrapper.DisplayImageWithTransition(sdk_wrapper.GetDataPath("images/languages/"+newLanguage+".png"), 1000, sdk_wrapper.IMAGE_TRANSITION_FADE_IN, 10)
	sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/languages/" + newLanguage + ".wav"))
	sdk_wrapper.DisplayImageWithTransition(sdk_wrapper.GetDataPath("images/languages/"+newLanguage+".png"), 1000, sdk_wrapper.IMAGE_TRANSITION_FADE_OUT, 10)

	// Patch and restart chipper
	vectorxPath := os.Getenv("VECTORX_HOME")
	chipperPatcherPath := path.Join(vectorxPath, "patchChipper.sh")
	cmd := exec.Command(chipperPatcherPath, loc, "&")
	err := cmd.Start()
	if err != nil {
		println(err.Error())
	}
	return returnIntent
}
