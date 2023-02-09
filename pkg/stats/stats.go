package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Stats struct {
	IntentsReceived int `json:"intents-received"`
	IntentsMatched  int `json:"intents-matched"`
}

func StatsIntentHandled(matched bool) {
	vPath := os.Getenv("VECTORX_HOME")
	statsFile := filepath.Join(vPath, "stats.json")
	data, err := os.ReadFile(statsFile)
	var jsonObj Stats = Stats{IntentsReceived: 0, IntentsMatched: 0}
	if err == nil {
		err = json.Unmarshal(data, &jsonObj)
	}
	if matched {
		jsonObj.IntentsMatched++
	}
	jsonObj.IntentsReceived++
	data, err = json.Marshal(jsonObj)
	if err == nil {
		err = os.WriteFile(statsFile, data, 0644)
	}
}

func GetStats() (Stats, error) {
	data := GetStatsJson()
	var jsonObj Stats
	err := json.Unmarshal([]byte(data), &jsonObj)
	return jsonObj, err
}

func GetStatsJson() string {
	vPath := os.Getenv("VECTORX_HOME")
	statsFile := filepath.Join(vPath, "stats.json")
	data, err := os.ReadFile(statsFile)
	if err == nil {
		return string(data)
	}
	return "{}"
}
