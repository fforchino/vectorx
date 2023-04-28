package main

import (
	"context"
	"flag"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"strings"
	"time"
	"vectorx/pkg/intents"
	"vectorx/pkg/stats"
)

var Debug = true
var Ctx = context.Background()
var Start = make(chan bool)
var Stop = make(chan bool)

func main() {
	sdkInit := false
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var locale = flag.String("locale", "", "STT Locale in use")
	var speechText = flag.String("speechText", "", "Speech text")
	flag.Parse()

	if Debug {
		println("SERIAL: " + *serial)
		println("LOCALE: " + *locale)
		println("SPEECH TEXT: " + *speechText)
	}
	language := *locale
	if strings.Contains(language, "-") {
		language = strings.Split(language, "-")[0]
	}

	if len(*speechText) > 0 {
		// Remove "" if any
		if strings.HasPrefix(*speechText, "\"") && strings.HasSuffix(*speechText, "\"") {
			*speechText = (*speechText)[1 : len(*speechText)-1]
		}

		// Register vectorx intents
		intents.RegisterIntents()

		// Find out whether the speech text matches any registered intent
		xIntent, err := intents.IntentMatch(*speechText, language)

		if err == nil {
			// Ok, we have a match. Then extract the parameters (if any) from the intent, but before that init the
			// SDK because we may need to get Vector's settings (e.g. for weather forecast we need the default location)

			stats.StatsIntentHandled(true)
			// Init SDK
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

			// Extract params
			intents.GetWirepodBotInfo(*serial)
			params := intents.ParseParams(*speechText, xIntent)

			engine, voice := sdk_wrapper.GetTTSConfiguration()
			sdk_wrapper.SetTTSEngine(engine)
			sdk_wrapper.SetTTSVoice(voice)

			go func() {
				_ = sdk_wrapper.Robot.BehaviorControl(Ctx, Start, Stop)
			}()

			for {
				select {
				case <-Start:
					returnIntent := xIntent.Handler(xIntent, *speechText, params)
					// Seems that we have to force back en_US locale or "Hey Vector" won't work anymore
					sdk_wrapper.SetLocale("en-US")
					// Ok, intent handled. Return the intent that Wirepod has to send to the robot
					fmt.Println("{\"status\": \"ok\", \"returnIntent\": \"" + returnIntent + "\"}")
					time.Sleep(3 * time.Second)
					sdk_wrapper.TriggerWakeWord()
					Stop <- true
				}
				return
			}
		} else {
			stats.StatsIntentHandled(false)
			// Intent cannot be handled by VectorX. Wirepod may continue its intent parsing chain
			if sdkInit {
				sdk_wrapper.SetLocale("en-US")
			}
			fmt.Println("{\"status\": \"ko\", \"returnIntent\": \"\"}")
		}
	} else {
		// Intent cannot be handled by VectorX. Wirepod may continue its intent parsing chain
		sdk_wrapper.SetLocale("en-US")
		fmt.Println("{\"status\": \"ko\", \"returnIntent\": \"\"}")
		sdk_wrapper.SetLocale("en-US")
	}
}
