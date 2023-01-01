package main

import (
	"context"
	"flag"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"time"
	"vectorx/pkg/intents"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	// Just called to add VIM localized strings to the engine
	intents.RegisterIntents()

	sdk_wrapper.InitSDKForWirepod(*serial)

	// Check for new messages forever
	for {
		messages, err := intents.VIMAPICheckMessages()
		if err == nil && len(messages) > 0 {
			for i := 0; i < len(messages); i++ {
				if !messages[i].Read {
					println(fmt.Sprintf("[%d] New message from %s: %s", messages[i].Timestamp, messages[i].From, messages[i].Message))
					var ctx = context.Background()
					var start = make(chan bool)
					var stop = make(chan bool)

					go func() {
						_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
					}()
					done := false
					for done == false {
						select {
						case <-start:
							intents.VIMAPIPlayMessage(messages[i])
							stop <- true
							done = true
						}
					}
					println("Message processed")
				}
			}
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
