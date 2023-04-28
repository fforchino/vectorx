package intents

import (
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
	"math"
	"time"
)

/**********************************************************************************************************************/
/*                                                HELLO WORLD                                                         */
/**********************************************************************************************************************/

func Balance_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"balance"}
	utterances[LOCALE_ITALIAN] = []string{"bilancia"}
	utterances[LOCALE_SPANISH] = []string{"escala de peso"}
	utterances[LOCALE_FRENCH] = []string{"échelle de poids"}
	utterances[LOCALE_GERMAN] = []string{"Gewichtsskala"}

	var intent = IntentDef{
		IntentName:            "extended_intent_hello_world",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               doBalance,
		OSKRTriggersUserInput: false,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_BALANCE_WEIGHT", []string{"%s1 grams", "%s1 grammi", "%s1 gramos", "%s1 grammes", "%s1 Gramm"})

	return nil
}

func doBalance(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	var maxAcc float64 = 0

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)
	sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/balance/balance.png"), 5000, true)

	// Read input asynchronously
	go func() {
		for {
			evt := sdk_wrapper.WaitForEvent()
			if evt != nil {
				evtRobotState := evt.GetRobotState()
				if evtRobotState != nil {
					// Update mximum acceleration on the Y axis
					maxAcc = math.Max(float64(evtRobotState.Accel.GetY()), maxAcc)
				}
			}
		}
	}()

	weight := "0.0"
	for i := 0; i <= 30; i++ {
		weight = fmt.Sprintf("%.1f", math.Round(maxAcc))
		sdk_wrapper.WriteText(weight, 64, true, 200, false)
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	sdk_wrapper.SayText(getTextEx("STR_BALANCE_WEIGHT", []string{weight}))
	sdk_wrapper.WriteText(weight, 64, true, 3000, true)

	return returnIntent
}
