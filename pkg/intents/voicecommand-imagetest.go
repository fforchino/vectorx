package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
)

/**********************************************************************************************************************/
/*                                                IMAGE TEST                                                          */
/**********************************************************************************************************************/

func ImageTest_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"image demo", "image test"}
	utterances[LOCALE_ITALIAN] = []string{"demo immagini", "prova immagini"}

	var intent = IntentDef{
		IntentName: "extended_intent_image_test",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    doImageTest,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func doImageTest(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	sdk_wrapper.SetBackgroundColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	sdk_wrapper.DisplayAnimatedGif(sdk_wrapper.GetDataPath("images/animated.gif"), sdk_wrapper.ANIMATED_GIF_SPEED_FASTEST, 3, false)
	sdk_wrapper.DisplayAnimatedGif(sdk_wrapper.GetDataPath("images/animated2.gif"), sdk_wrapper.ANIMATED_GIF_SPEED_FASTEST, 3, false)
	sdk_wrapper.DisplayAnimatedGif(sdk_wrapper.GetDataPath("images/animated3.gif"), sdk_wrapper.ANIMATED_GIF_SPEED_FASTEST, 3, false)
	return returnIntent
}
