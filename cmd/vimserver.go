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

	// Check for new messages forever
	for {
		// Load the Vector Serial Numbers for which we are going to check messages from Wirepod
		serials := getMyBotSerials()

		if intents.VIMEnabled {
			//println("VIM Enabled")
			for _, serial := range serials {
				isChatty, lastMessageId := isBotInChatMood(serial)
				if isChatty {
					//println(fmt.Sprintf("Checking messages for %s %d", serial, lastMessageId))
					messages, err := intents.VIMAPICheckMessages(serial, lastMessageId)
					if err == nil && len(messages) > 0 {
						var ctx = context.Background()
						var start = make(chan bool)
						var stop = make(chan bool)
						sdk_wrapper.InitSDKForWirepod(serial)
						go func() {
							_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
						}()
						select {
						case <-start:
							for i := 0; i < len(messages); i++ {
								if !messages[i].Read {
									println(fmt.Sprintf("[%d] New message from %s: %s", messages[i].Timestamp, messages[i].From, messages[i].Message))
									intents.VIMAPIPlayMessage(messages[i])
								}
							}
							println("Messages for bot processed")
							stop <- true
						}
					}
				}
			}
		}

		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}

// Gets a list of serials managed by this instance of wirepod

func getMyBotSerials() []string {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	var botConfigJSON botConfigStruct = nil

	var serials []string
	botConfigFile := filepath.Join(wirepodPath, "chipper/botConfig.json")
	//println("Looking at " + botConfigFile + "...")
	if _, err := os.Stat(botConfigFile); err == nil {
		botConfigJSONFile, err := os.ReadFile(botConfigFile)
		if err == nil {
			json.Unmarshal(botConfigJSONFile, &botConfigJSON)
		}
	}
	for _, botConfig := range botConfigJSON {
		//println("Will take care of bot # " + botConfig.ESN)
		serials = append(serials, botConfig.ESN)
	}
	return serials
}

// Checks whether the given bot is logged into chat

func isBotInChatMood(serial string) (bool, int32) {
	// Peek into the given vector custom settings and read the value
	customSettingsPath := filepath.Join(os.Getenv("VECTORX_HOME"), "vectorfs")
	customSettingsPath = filepath.Join(customSettingsPath, "nvm")
	customSettingsPath = filepath.Join(customSettingsPath, serial)
	customSettingsPath = filepath.Join(customSettingsPath, "custom_settings.json")

	println(customSettingsPath)
	botCustomSettingsJSONFile, err := os.ReadFile(customSettingsPath)
	if err == nil {
		var botCustomSettings sdk_wrapper.CustomSettings
		err := json.Unmarshal(botCustomSettingsJSONFile, &botCustomSettings)
		if err == nil {
			//println("OK")
			return botCustomSettings.LoggedInToChat, botCustomSettings.LastChatMessageRead
		}
	}
	return false, -1
}
