package main

import (
	"context"
	"flag"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"strings"
	"vectorx/pkg/intents"
)

var Debug = false
var Ctx = context.Background()
var Start = make(chan bool)
var Stop = make(chan bool)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var locale = flag.String("locale", "", "STT Locale in use")
	var speechText = flag.String("speechText", "", "Speech text")
	flag.Parse()

	if Debug {
		println("SERIAL: " + *serial)
		println("LOCALE: " + *locale)
		println("SPEECH TEXT: " + *speechText)
	}

	if len(*speechText) > 0 {
		// Remove "" if any
		if strings.HasPrefix(*speechText, "\"") && strings.HasSuffix(*speechText, "\"") {
			*speechText = (*speechText)[1 : len(*speechText)-1]
		}

		// Register vectorx intents
		intents.RegisterIntents()
		intents.GetWirepodBotInfo(*serial)

		err := sdk_wrapper.InitSDKForWirepod(*serial)
		if err != nil {
			println("FATAL: could not load Vector settings from JDOCS")
			return
		}

		robotLocale := sdk_wrapper.GetLocale()
		if Debug {
			println("ROBOT LOCALE: " + robotLocale)
		}
		if robotLocale != *locale {
			if Debug {
				println("Different locales! Setting " + *locale)
			}
			sdk_wrapper.SetLocale(*locale)
		}
		if Debug {
			println("ROBOT LOCALE: " + robotLocale)
		}
		sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)

		// Make sure "locale" is just the language name
		if strings.Contains(*locale, "-") {
			*locale = strings.Split(*locale, "-")[0]
		}
		// Find out whether the speech text matches any registered intent
		xIntent, err := intents.IntentMatch(*speechText, *locale)

		if err == nil {
			// Ok, we have a match. Then extract the parameters (if any) from the intent...
			params := intents.ParseParams(*speechText, xIntent)

			go func() {
				_ = sdk_wrapper.Robot.BehaviorControl(Ctx, Start, Stop)
			}()

			for {
				select {
				case <-Start:
					returnIntent := xIntent.Handler(xIntent, *speechText, params)
					// Seems that we have to force back en_US locale or "Hey Vector" won't work anymore
					sdk_wrapper.SetLocale("en_US")
					// Ok, intent handled. Return the intent that Wirepod has to send to the robot
					fmt.Println("{\"status\": \"ok\", \"returnIntent\": \"" + returnIntent + "\"}")
					Stop <- true
				}
				return
			}
		} else {
			// Intent cannot be handled by VectorX. Wirepod may continue its intent parsing chain
			sdk_wrapper.SetLocale("en_US")
			fmt.Println("{\"status\": \"ko\", \"returnIntent\": \"\"}")
			sdk_wrapper.SetLocale("en_US")
		}
	} else {
		// Intent cannot be handled by VectorX. Wirepod may continue its intent parsing chain
		sdk_wrapper.SetLocale("en_US")
		fmt.Println("{\"status\": \"ko\", \"returnIntent\": \"\"}")
		sdk_wrapper.SetLocale("en_US")
	}
}
