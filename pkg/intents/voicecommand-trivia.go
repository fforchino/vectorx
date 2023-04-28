package intents

import (
	"encoding/json"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"strconv"
	"strings"
)

type TriviaGameData struct {
	GameName        string `json:"gameName"`
	CurrentQuestion int    `json:"currentQuestion"`
	Score           int    `json:"score"`
	State           string `json:"state"`
}

const STATE_TRIVIA_GAME_STARTED = "started"
const STATE_TRIVIA_GAME_NOT_STARTED = "none"

const TRIVIA_ANSWER_UNKNOWN = -1
const TRIVIA_ANSWER_1 = 1
const TRIVIA_ANSWER_2 = 2
const TRIVIA_ANSWER_3 = 3
const TRIVIA_ANSWER_4 = 4
const TRIVIA_ANSWER_QUIT = 5
const TRIVIA_GAME_NAME = "TRIVIA_GAME"

var TRIVIA_SERVER_URL = os.Getenv("VIM_SERVER") // "https://www.wondergarden.app/VECTOR-TRIVIA" //"http://192.168.43.65/VECTOR-TRIVIA"
var TriviaDebug = true
var GameConfig = TriviaGameData{GameName: TRIVIA_GAME_NAME,
	CurrentQuestion: 1,
	Score:           0,
	State:           STATE_TRIVIA_GAME_NOT_STARTED,
}

func Trivia_Register(intentList *[]IntentDef) error {
	addLocalizedString("STR_OK_LETS_GO", []string{"Ok, let's go!", "Perfetto, andiamo!", "", "", ""})
	addLocalizedString("STR_GAME_OVER", []string{"Game over", "Fine partita", "", "", ""})
	addLocalizedString("STR_QUESTION_NUM", []string{"Question %s1", "Domanda numero %s1", "", "", ""})

	addLocalizedString("STR_FIRST", []string{"first", "prima", "", "", ""})
	addLocalizedString("STR_SECOND", []string{"second", "seconda", "", "", ""})
	addLocalizedString("STR_THIRD", []string{"third", "terza", "", "", ""})
	addLocalizedString("STR_FOURTH", []string{"fourth", "quarta", "", "", ""})
	addLocalizedString("STR_QUIT", []string{"quit", "esci", "", "", ""})
	addLocalizedString("STR_CORRECT_ANSWER", []string{"correct!", "giusto!", "", "", ""})
	addLocalizedString("STR_WRONG_ANSWER", []string{"wrong!", "sbagliato!", "", "", ""})

	registerTriviaIntent(intentList)

	return nil
}

func triviaGameStarted() bool {
	err := json.Unmarshal([]byte(sdk_wrapper.GetCurrentGame()), &GameConfig)
	return (err == nil && GameConfig.GameName == TRIVIA_GAME_NAME && GameConfig.State == STATE_TRIVIA_GAME_STARTED)
}

func saveConfig() bool {
	retVal := false
	b, err := json.Marshal(GameConfig)
	if err == nil {
		sdk_wrapper.SetCurrentGame(string(b))
		retVal = true
	}
	return retVal
}

func setTriviaGameStart() bool {
	retVal := false
	if !triviaGameStarted() {
		retVal = saveConfig()
	}
	return retVal
}

func setTriviaGameEnd() bool {
	retVal := false
	if triviaGameStarted() {
		GameConfig.State = STATE_TRIVIA_GAME_NOT_STARTED
		GameConfig.Score = 0
		GameConfig.CurrentQuestion = 1
		b, err := json.Marshal(GameConfig)
		if err == nil {
			sdk_wrapper.SetCurrentGame(string(b))
			retVal = true
		}
	}
	return retVal
}

/**********************************************************************************************************************/
/*                                            TRIGGER TRIVIA GAME                                                     */
/**********************************************************************************************************************/

func registerTriviaIntent(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"trivia game"}
	utterances[LOCALE_ITALIAN] = []string{"gioco delle domande"}
	utterances[LOCALE_SPANISH] = []string{""}
	utterances[LOCALE_FRENCH] = []string{""}
	utterances[LOCALE_GERMAN] = []string{""}

	var intent = IntentDef{
		IntentName:            "extended_intent_trivia",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               triggerTriviaGame,
		OSKRTriggersUserInput: handleTriviaNextInput,
	}
	*intentList = append(*intentList, intent)

	// If the game is started, register a catch-all intent that will capture any input
	// Note: only a catchall intent should be active at a time!!!
	if triviaGameStarted() {
		registerCatchallIntent(intentList)
	}

	return nil
}

func triggerTriviaGame(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	sdk_wrapper.SayText(getText("STR_OK_LETS_GO"))
	if setTriviaGameStart() {
		gotoQuestion(1)
		returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                            HANDLE ANSWERS                                                          */
/**********************************************************************************************************************/

func registerCatchallIntent(intentList *[]IntentDef) {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"*"}
	utterances[LOCALE_ITALIAN] = []string{"*"}
	utterances[LOCALE_SPANISH] = []string{"*"}
	utterances[LOCALE_FRENCH] = []string{"*"}
	utterances[LOCALE_GERMAN] = []string{"*"}

	var intent = IntentDef{
		IntentName:            "extended_intent_trivia_input",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               handleTriviaInput,
		OSKRTriggersUserInput: handleTriviaNextInput,
	}
	*intentList = append(*intentList, intent)
}

func handleTriviaNextInput() bool {
	return triviaGameStarted()
}

func handleTriviaInput(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE

	if triviaGameStarted() {
		userAnswer := TRIVIA_ANSWER_UNKNOWN
		if strings.Contains(speechText, getText("STR_QUIT")) {
			// Quit the game
			sdk_wrapper.SayText(getText("STR_GAME_OVER"))
			setTriviaGameEnd()
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
			userAnswer = TRIVIA_ANSWER_QUIT
		} else if strings.Contains(speechText, getText("STR_FIRST")) {
			userAnswer = TRIVIA_ANSWER_1
		} else if strings.Contains(speechText, getText("STR_SECOND")) {
			userAnswer = TRIVIA_ANSWER_2
		} else if strings.Contains(speechText, getText("STR_THIRD")) {
			userAnswer = TRIVIA_ANSWER_3
		} else if strings.Contains(speechText, getText("STR_FOURTH")) {
			userAnswer = TRIVIA_ANSWER_4
		}

		if userAnswer == TRIVIA_ANSWER_1 {
			sdk_wrapper.SayText(getText("STR_CORRECT_ANSWER"))
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
			gotoQuestion(GameConfig.CurrentQuestion + 1)
		} else if userAnswer == TRIVIA_ANSWER_UNKNOWN {
			returnIntent = STANDARD_INTENT_IMPERATIVE_NEGATIVE
			gotoQuestion(GameConfig.CurrentQuestion)
		} else {
			sdk_wrapper.SayText(getText("STR_WRONG_ANSWER"))
			returnIntent = STANDARD_INTENT_IMPERATIVE_NEGATIVE
			gotoQuestion(GameConfig.CurrentQuestion + 1)
		}

	}

	return returnIntent
}

/**********************************************************************************************************************/
/*                                            HANDLE ANSWERS                                                          */
/**********************************************************************************************************************/

func gotoQuestion(questionNum int) {
	if triviaGameStarted() {
		// Ask question
		GameConfig.CurrentQuestion = questionNum
		saveConfig()
		sdk_wrapper.SayText(getTextEx("STR_QUESTION_NUM", []string{strconv.Itoa(questionNum)}))
	}
}
