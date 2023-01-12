package main

import (
	"context"
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"path/filepath"
	"time"
	"vectorx/pkg/intents"
)

type botConfigStruct []struct {
	ESN             string `json:"esn"`
	Location        string `json:"location"`
	Units           string `json:"units"`
	UsePlaySpecific bool   `json:"use_play_specific"`
	IsEarlyOpus     bool   `json:"is_early_opus"`
}

func main() {
	// Just called to add VIM localized strings to the engine
	intents.RegisterIntents()

	// Load the Vector Serial Numbers for which we are going to check messages from Wirepod
	serials := getMyBotSerials()

	// Check for new messages forever
	for {
		for _, serial := range serials {
			messages, err := intents.VIMAPICheckMessages(serial)
			if err == nil && len(messages) > 0 {
				for i := 0; i < len(messages); i++ {
					if !messages[i].Read {
						println(fmt.Sprintf("[%d] New message from %s: %s", messages[i].Timestamp, messages[i].From, messages[i].Message))
						var ctx = context.Background()
						var start = make(chan bool)
						var stop = make(chan bool)

						sdk_wrapper.InitSDKForWirepod(serial)
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
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}

func getMyBotSerials() []string {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	var botConfigJSON botConfigStruct = nil

	var serials []string
	botConfigFile := filepath.Join(wirepodPath, "chipper/botConfig.json")
	println("Looking at " + botConfigFile + "...")
	if _, err := os.Stat(botConfigFile); err == nil {
		botConfigJSONFile, err := os.ReadFile(botConfigFile)
		if err == nil {
			json.Unmarshal(botConfigJSONFile, &botConfigJSON)
		}
	}
	for _, botConfig := range botConfigJSON {
		println("Will take care of bot # " + botConfig.ESN)
		serials = append(serials, botConfig.ESN)
	}
	return serials
}
