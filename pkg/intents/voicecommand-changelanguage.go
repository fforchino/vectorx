package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"os/exec"
	"path"
)

/**********************************************************************************************************************/
/*                                                CHANGE LANGUAGE                                                     */
/**********************************************************************************************************************/

func ChangeLanguage_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"let's talk in"}
	utterances[LOCALE_ITALIAN] = []string{"parliamo in"}
	utterances[LOCALE_SPANISH] = []string{"hablamos en"}
	utterances[LOCALE_FRENCH] = []string{"parlons en"}
	utterances[LOCALE_GERMAN] = []string{"sprechen auf"}

	var intent = IntentDef{
		IntentName: "extended_intent_changelanguage",
		Utterances: utterances,
		Parameters: []string{PARAMETER_LANGUAGE},
		Handler:    changeLanguage,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func changeLanguage(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	currentLanguage := sdk_wrapper.GetLanguage()
	isSameLanguage := false
	loc := "en-US"
	switch params.Language {
	case LOCALE_ITALIAN:
		sdk_wrapper.SetLocale("it-IT")
		loc = "it-IT"
		if currentLanguage == sdk_wrapper.LANGUAGE_ITALIAN {
			isSameLanguage = true
		}
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_ITALIAN)
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("languages/it.png"), 500, false)
		break
	case LOCALE_SPANISH:
		sdk_wrapper.SetLocale("es-ES")
		loc = "es-ES"
		if currentLanguage == sdk_wrapper.LANGUAGE_SPANISH {
			isSameLanguage = true
		}
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_SPANISH)
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("languages/es.png"), 500, false)
		break
	case LOCALE_FRENCH:
		sdk_wrapper.SetLocale("fr-FR")
		loc = "fr-FR"
		if currentLanguage == sdk_wrapper.LANGUAGE_FRENCH {
			isSameLanguage = true
		}
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_FRENCH)
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("languages/fr.png"), 500, false)
		break
	case LOCALE_GERMAN:
		sdk_wrapper.SetLocale("de-DE")
		loc = "de-DE"
		if currentLanguage == sdk_wrapper.LANGUAGE_GERMAN {
			isSameLanguage = true
		}
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_GERMAN)
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("languages/de.png"), 500, false)
		break
	default:
		sdk_wrapper.SetLocale("en-US")
		if currentLanguage == sdk_wrapper.LANGUAGE_ENGLISH {
			isSameLanguage = true
		}
		sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_ENGLISH)
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("languages/en.png"), 500, false)
		break
	}
	sdk_wrapper.SayText(getText(STR_HELLO_WORLD))
	if !isSameLanguage {
		// Patch and restart chipper
		vectorxPath := os.Getenv("VECTORX_HOME")
		vectorxPath = path.Join(vectorxPath, "patchChipper.sh")
		cmd := exec.Command("patchChipper.sh", loc)
		err := cmd.Run()
		if err != nil {
			println(err.Error())
		}
	}
	return returnIntent
}
