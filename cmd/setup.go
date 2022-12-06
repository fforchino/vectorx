package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type intentsStruct []struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Utterances  []string `json:"utterances"`
	Intent      string   `json:"intent"`
	Params      struct {
		ParamName  string `json:"paramname"`
		ParamValue string `json:"paramvalue"`
	} `json:"params"`
	Exec           string   `json:"exec"`
	ExecArgs       []string `json:"execargs"`
	IsSystemIntent bool     `json:"issystem"`
}

func contains(iss intentsStruct, i string) bool {
	for _, itx := range iss {
		if itx.Name == i {
			return true
		}
	}
	return false
}

func main() {
	var customIntentJSON intentsStruct = nil
	var vectorxIntentJSON intentsStruct = nil

	wirepodPath := os.Getenv("WIREPOD_HOME")

	// Load WP custom intents
	customIntentsFile := filepath.Join(wirepodPath, "chipper/customIntents.json")
	println("Looking at " + customIntentsFile + "...")
	if _, err := os.Stat(customIntentsFile); err == nil {
		customIntentJSONFile, err := os.ReadFile(customIntentsFile)
		if err == nil {
			json.Unmarshal(customIntentJSONFile, &customIntentJSON)
		}
	}
	// Load VECTORX custom intents
	vectorxIntentsFile := "./vectorxSetup.json"
	if _, err := os.Stat(vectorxIntentsFile); err == nil {
		vectorxIntentsJSONFile, err := os.ReadFile(vectorxIntentsFile)
		if err == nil {
			json.Unmarshal(vectorxIntentsJSONFile, &vectorxIntentJSON)
		}
	}

	// Overwrite vectorx intents in wirepod custom intents
	if nil != vectorxIntentJSON {
		for i, v := range vectorxIntentJSON {
			cwd, _ := os.Getwd()
			v.Exec = strings.Replace(v.Exec, "__REPLACEME__", cwd, -1)
			if !contains(customIntentJSON, v.Name) {
				println("Appending intent " + v.Name)
				customIntentJSON = append(customIntentJSON, v)
			} else {
				println("Overwriting intent " + v.Name)
				customIntentJSON[i] = v
			}
		}
	}

	// Flush
	finalJsonBytes, _ := json.Marshal(customIntentJSON)
	os.WriteFile(customIntentsFile, finalJsonBytes, 0644)
}
