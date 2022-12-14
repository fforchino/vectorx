package intents

import (
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"path"
	"strings"
)

// Wirepod supported locales

const LOCALE_ENGLISH = "en"
const LOCALE_ITALIAN = "it"
const LOCALE_SPANISH = "es"
const LOCALE_FRENCH = "fr"
const LOCALE_GERMAN = "de"

// Supported parameters

const PARAMETER_USERNAME = "PARAMETER_USERNAME"
const PARAMETER_LANGUAGE = "PARAMETER_LANGUAGE"
const PARAMETER_WEATHER = "PARAMETER_WEATHER"

// Standard intents a production Vector understands

const STANDARD_INTENT_NAMES_USERNAME_EXTEND = "intent_names_username_extend"
const STANDARD_INTENT_WEATHER_EXTEND = "intent_weather_extend"
const STANDARD_INTENT_NAMES_ASK = "intent_names_ask"
const STANDARD_INTENT_IMPERATIVE_EYECOLOR = "intent_imperative_eyecolor"
const STANDARD_INTENT_CHARACTER_AGE = "intent_character_age"
const STANDARD_INTENT_EXPLORE_START = "intent_explore_start"
const STANDARD_INTENT_SYSTEM_CHARGER = "intent_system_charger"
const STANDARD_INTENT_SYSTEM_SLEEP_ = "intent_system_sleep"
const STANDARD_INTENT_GREETING_GOODMORNING = "intent_greeting_goodmorning"
const STANDARD_INTENT_GREETING_GOODNIGHT = "intent_greeting_goodnight"
const STANDARD_INTENT_GREETING_GOODBYE = "intent_greeting_goodbye"
const STANDARD_INTENT_SEASONAL_HAPPYNEWYEAR = "intent_seasonal_happynewyear"
const STANDARD_INTENT_SEASONAL_HAPPY_HOLIDAYS = "intent_seasonal_happyholidays"
const STANDARD_INTENT_AMAZON_SIGNIN = "intent_amazon_signin"
const STANDARD_INTENT_AMAZON_SIGNOUT = "intent_amazon_signout"
const STANDARD_INTENT_IMPERATIVE_FORWARD = "intent_imperative_forward"
const STANDARD_INTENT_IMPERATIVE_TURNAROUND = "intent_imperative_turnaround"
const STANDARD_INTENT_IMPERATIVE_TURNLEFT = "intent_imperative_turnleft"
const STANDARD_INTENT_IMPERATIVE_TURNRIGHT = "intent_imperative_turnright"
const STANDARD_INTENT_PLAY_ROLLCUBE = "intent_play_rollcube"
const STANDARD_INTENT_PLAY_POPAWHEELIE = "intent_play_popawheelie"
const STANDARD_INTENT_PLAY_FISTBUMP = "intent_play_fistbump"
const STANDARD_INTENT_PLAY_BLACKJACK = "intent_play_blackjack"
const STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE = "intent_imperative_affirmative"
const STANDARD_INTENT_IMPERATIVE_NEGATIVE = "intent_imperative_negative"
const STANDARD_INTENT_PHOTO_TAKE_EXTEND_ = "intent_photo_take_extend"
const STANDARD_INTENT_IMPERATIVE_PRAISE = "intent_imperative_praise"
const STANDARD_INTENT_IMPERATIVE_ABUSE = "intent_imperative_abuse"
const STANDARD_INTENT_IMPERATIVE_APOLOGIZE = "intent_imperative_apologize"
const STANDARD_INTENT_IMPERATIVE_BACKUP = "intent_imperative_backup"
const STANDARD_INTENT_IMPERATIVE_VOLUMEDOWN = "intent_imperative_volumedown"
const STANDARD_INTENT_IMPERATIVE_VOLUMEUP = "intent_imperative_volumeup"
const STANDARD_INTENT_IMPERATIVE_LOOKATME = "intent_imperative_lookatme"
const STANDARD_INTENT_IMPERATIVE_VOLUMELEVEL_EXTEND = "intent_imperative_volumelevel_extend"
const STANDARD_INTENT_IMPERATIVE_SHUTUP = "intent_imperative_shutup"
const STANDARD_INTENT_GREETING_HELLO = "intent_greeting_hello"
const STANDARD_INTENT_IMPERATIVE_COME = "intent_imperative_come"
const STANDARD_INTENT_IMPERATIVE_LOVE = "intent_imperative_love"
const STANDARD_INTENT_PROMPTQUESTION = "intent_knowledge_promptquestion"
const STANDARD_INTENT_CHECKTIMER = "intent_clock_checktimer"
const STANDARD_INTENT_GLOBAL_STOP_EXTEND = "intent_global_stop_extend"
const STANDARD_INTENT_SETTIMER_EXTEND = "intent_clock_settimer_extend"
const STANDARD_INTENT_CLOCK_TIME = "intent_clock_time"
const STANDARD_INTENT_IMPERATIVE_QUIET = "intent_imperative_quiet"
const STANDARD_INTENT_IMPERATIVE_DANCE = "intent_imperative_dance"
const STANDARD_INTENT_PLAY_PICKUPCUBE = "intent_play_pickupcube"
const STANDARD_INTENT_IMPERATIVE_FETCHCUBE = "intent_imperative_fetchcube"
const STANDARD_INTENT_IMPERATIVE_FINDCUBE = "intent_imperative_findcube"
const STANDARD_INTENT_PLAY_ANYTRICK = "intent_play_anytrick"
const STANDARD_INTENT_RECORDMESSAGE_EXTEND = "intent_message_recordmessage_extend"
const STANDARD_INTENT_PLAYMESSAGE_EXTEND = "intent_message_playmessage_extend"
const STANDARD_INTENT_BLACKJACK_HIT = "intent_blackjack_hit"
const STANDARD_INTENT_BLACKJACK_STAND = "intent_blackjack_stand"
const STANDARD_INTENT_KEEPAWAY = "intent_play_keepaway"
const STANDARD_INTENT_SYSTEM_NOAUDIO = "intent_system_noaudio"

type WeatherParams struct {
	Condition               string
	IsForecast              string
	LocalDatetime           string
	SpeakableLocationString string
	Temperature             string
	TemperatureUnit         string
	Icon                    string
}

type IntentParams struct {
	RobotName string
	Language  string
	Weather   WeatherParams
}

type IntentHandlerFunc func(IntentDef, string, IntentParams) string

type IntentDef struct {
	IntentName string
	Utterances map[string][]string
	Parameters []string
	Handler    IntentHandlerFunc
}

var intents []IntentDef

var BotLocation string = ""
var BotUnits string = ""

func RegisterIntents() {
	HelloWorld_Register(&intents)
	RollaDie_Register(&intents)
	RobotName_Register(&intents)
	HowDoYouSay_Register(&intents)
	ChangeLanguage_Register(&intents)
	Weather_Register(&intents)
	RockPaperScissors_Register(&intents)
}

func IntentMatch(speechText string, locale string) (IntentDef, error) {
	var candidates1 []IntentDef
	var candidates2 []IntentDef
	maxLen := 0
	cIntent := 0
	for _, intent := range intents {
		if hasPerfectMatch(intent.Utterances[locale], speechText) {
			candidates1 = append(candidates1, intent)
		}
	}
	// Return the intent with longer utterance matched
	if len(candidates1) > 0 {
		for idx, intent := range candidates1 {
			if len(intent.Utterances[locale]) > maxLen {
				maxLen = len(intent.Utterances[locale])
				cIntent = idx
			}
		}
		return candidates1[cIntent], nil
	}

	for _, intent := range intents {
		if hasPartialMatch(intent.Utterances[locale], speechText) {
			candidates2 = append(candidates1, intent)
		}
	}
	maxLen = 0
	cIntent = 0
	// Return the intent with longer utterance matched
	if len(candidates2) > 0 {
		for idx, intent := range candidates2 {
			if len(intent.Utterances[locale]) > maxLen {
				maxLen = len(intent.Utterances[locale])
				cIntent = idx
			}
		}
		return candidates2[cIntent], nil
	}
	return IntentDef{}, fmt.Errorf("Intent not found")
}

// Try to get Wirepod config for this bot if any. Else get params from the robot itself
func GetWirepodBotInfo(serialNo string) {
	botConfigPath := os.Getenv("WIREPOD_HOME")
	botConfigPath = path.Join(botConfigPath, "chipper", "botConfig.json")
	if _, err := os.Stat(botConfigPath); err == nil {
		type botConfigJSON []struct {
			ESN             string `json:"ESN"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}

		var botConfig botConfigJSON

		byteValue, err := os.ReadFile(botConfigPath)
		if err != nil {
			println(err)
		}

		json.Unmarshal(byteValue, &botConfig)
		for _, bot := range botConfig {
			if strings.ToLower(bot.ESN) == serialNo {
				println("Found bot config for " + bot.ESN)
				BotLocation = bot.Location
				BotUnits = bot.Units
			}
		}
		if BotLocation == "" {
			BotLocation = sdk_wrapper.GetVectorSettings()["default_location"].(string)
			BotUnits = os.Getenv("WEATHERAPI_UNIT")
		}
	}
}

/**********************************************************************************************************************/
/*                                                PRIVATE FUNCTIONS                                                   */
/**********************************************************************************************************************/

func hasPerfectMatch(utterances []string, phrase string) bool {
	for _, s := range utterances {
		if strings.ToLower(s) == strings.ToLower(phrase) {
			return true
		}
	}
	return false
}

func hasPartialMatch(utterances []string, phrase string) bool {
	for _, s := range utterances {
		if strings.Contains(strings.ToLower(phrase), strings.ToLower(s)) {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
