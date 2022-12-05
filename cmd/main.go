package main

import (
	"context"
	"flag"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"vectorx/pkg/intents"
	"vectorx/pkg/voicecommands"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var intent = flag.String("intent", "", "Extended intent to fire")
	var speechText = flag.String("speechText", "", "Speech text")
	flag.Parse()

	println("SERIAL: " + *serial)
	println("INTENT: " + *intent)
	println("SPEECH TEXT: " + *speechText)

	if len(*speechText) > 0 {
		sdk_wrapper.InitSDKForWirepod(*serial)
		params := intents.ParseParams(*speechText, *intent)

		ctx := context.Background()
		start := make(chan bool)
		stop := make(chan bool)

		go func() {
			_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
		}()

		for {
			select {
			case <-start:
				returnIntent := "intent_system_noaudio"
				switch *intent {
				case "extended_intent_rolladie":
					returnIntent = voicecommands.RollADie()
					break
				case "extended_intent_set_robot_name":
					returnIntent = voicecommands.SetRobotName(params)
					break
				case "extended_intent_say_robot_name":
					returnIntent = voicecommands.SayRobotName()
					break
				}

				stop <- true
				print(returnIntent)
				return
			}
		}
	}
}
