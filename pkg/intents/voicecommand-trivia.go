package intents

import (
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"os"
	"strconv"
	"strings"
	"time"
)

const TRIVIA_ANSWER_UNKNOWN = -1
const TRIVIA_ANSWER_1 = 1
const TRIVIA_ANSWER_2 = 2
const TRIVIA_ANSWER_3 = 3
const TRIVIA_ANSWER_4 = 4
const TRIVIA_ANSWER_QUIT = 5
const TRIVIA_GAME_NAME = "TRIVIA_GAME"

var TRIVIA_SERVER_URL = os.Getenv("VIM_SERVER") // "https://www.wondergarden.app/VECTOR-TRIVIA" //"http://192.168.43.65/VECTOR-TRIVIA"
var TriviaDebug = true
var CurrentQuestion = 0

func Trivia_Register(intentList *[]IntentDef) error {
	addLocalizedString("STR_OK_LETS_GO", []string{"Ok, let's go!", "Perfetto, andiamo!", "", "", ""})
	addLocalizedString("STR_GAME_OVER", []string{"Game over", "Fine partita", "", "", ""})
	addLocalizedString("STR_QUESTION_NUM", []string{"Question %s1", "Domanda numero %s1", "", "", ""})

	addLocalizedString("FIRST", []string{"first", "prima", "", "", ""})
	addLocalizedString("SECOND", []string{"second", "seconda", "", "", ""})
	addLocalizedString("THIRD", []string{"third", "terza", "", "", ""})
	addLocalizedString("FOURTH", []string{"fourth", "quarta", "", "", ""})
	addLocalizedString("QUIT", []string{"quit", "esci", "", "", ""})

	registerTriviaIntent(intentList)
	registerTriviaAnswers(intentList)

	return nil
}

func triviaGameStarted() bool {
	return sdk_wrapper.GetCurrentGame() == TRIVIA_GAME_NAME
}

func setTriviaGameStart() bool {
	retVal := false
	if !triviaGameStarted() {
		sdk_wrapper.SetCurrentGame(TRIVIA_GAME_NAME)
		retVal = true
	}
	return retVal
}

func setTriviaGameEnd() bool {
	retVal := false
	if triviaGameStarted() {
		sdk_wrapper.SetCurrentGame("")
		retVal = true
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
		OSKRTriggersUserInput: true,
	}
	*intentList = append(*intentList, intent)

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

func registerTriviaAnswers(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{
		getTextIn("STR_QUIT", []string{}, en_US),
		getTextIn("STR_FIRST", []string{}, en_US),
		getTextIn("STR_SECOND", []string{}, en_US),
		getTextIn("STR_THIRD", []string{}, en_US),
		getTextIn("STR_FOURTH", []string{}, en_US),
	}
	utterances[LOCALE_ITALIAN] = []string{
		getTextIn("STR_QUIT", []string{}, it_IT),
		getTextIn("STR_FIRST", []string{}, it_IT),
		getTextIn("STR_SECOND", []string{}, it_IT),
		getTextIn("STR_THIRD", []string{}, it_IT),
		getTextIn("STR_FOURTH", []string{}, it_IT),
	}
	utterances[LOCALE_SPANISH] = []string{
		getTextIn("STR_QUIT", []string{}, es_ES),
		getTextIn("STR_FIRST", []string{}, es_ES),
		getTextIn("STR_SECOND", []string{}, es_ES),
		getTextIn("STR_THIRD", []string{}, es_ES),
		getTextIn("STR_FOURTH", []string{}, es_ES),
	}
	utterances[LOCALE_FRENCH] = []string{
		getTextIn("STR_QUIT", []string{}, fr_FR),
		getTextIn("STR_FIRST", []string{}, fr_FR),
		getTextIn("STR_SECOND", []string{}, fr_FR),
		getTextIn("STR_THIRD", []string{}, fr_FR),
		getTextIn("STR_FOURTH", []string{}, fr_FR),
	}
	utterances[LOCALE_GERMAN] = []string{
		getTextIn("STR_QUIT", []string{}, de_DE),
		getTextIn("STR_FIRST", []string{}, de_DE),
		getTextIn("STR_SECOND", []string{}, de_DE),
		getTextIn("STR_THIRD", []string{}, de_DE),
		getTextIn("STR_FOURTH", []string{}, de_DE),
	}

	var intent = IntentDef{
		IntentName:            "extended_intent_trivia_input",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               handleTriviaInput,
		OSKRTriggersUserInput: true,
	}
	*intentList = append(*intentList, intent)
	return nil
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

		if userAnswer == TRIVIA_ANSWER_UNKNOWN {
			go func() {
				time.Sleep(3 * time.Second)
				sdk_wrapper.TriggerWakeWord()
			}()
		}
		if userAnswer == TRIVIA_ANSWER_1 {
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
			go func() {
				time.Sleep(3 * time.Second)
				CurrentQuestion = CurrentQuestion + 1
				gotoQuestion(CurrentQuestion)
			}()
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
		CurrentQuestion = questionNum
		sdk_wrapper.SayText(getTextEx("STR_QUESTION_NUM", []string{strconv.Itoa(questionNum)}))
	}
}
