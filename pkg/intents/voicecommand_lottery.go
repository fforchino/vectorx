package intents

import (
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"image/color"
	"math/rand"
	"time"
)

/**********************************************************************************************************************/
/*                                          NEW YEAR'S EVE BINGO                                                      */
/**********************************************************************************************************************/

const MAX_NUMBERS = 90

func Lottery_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"bingo"}
	utterances[LOCALE_ITALIAN] = []string{"tombola"}
	utterances[LOCALE_SPANISH] = []string{"bingo"}
	utterances[LOCALE_FRENCH] = []string{"bingo"}
	utterances[LOCALE_GERMAN] = []string{"bingo"}

	var intent = IntentDef{
		IntentName: "extended_intent_bingo",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    bingo,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_BINGO_BINGO", []string{"bingo time!", "tombola!", "bingo!", "bingo!", "bingo!"})
	addLocalizedString("STR_BINGO_GAME_OVER", []string{"game over!", "la partita è finita!", "juego terminado!", "jeu terminé!", "spiel ist aus!"})

	return nil
}

func bingo(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	doBingo()
	return returnIntent
}

func doBingo() {
	s1 := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(s1)
	gameOver := false

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)
	isBusy := false
	sdk_wrapper.SayText(getText("STR_BINGO_BINGO"))
	sdk_wrapper.PlayAnimation("anim_greeting_hello_01", 1, false, false, false)

	// Read input asynchronously
	go func() {
		for {
			evt := sdk_wrapper.WaitForEvent()
			if evt != nil {
				evtUserIntent := evt.GetUserIntent()
				evtRobotState := evt.GetRobotState()
				if evtUserIntent != nil {
					println(fmt.Sprintf("Received intent %d", evtUserIntent.IntentId))
					println(evtUserIntent.JsonData)
					println(evtUserIntent.String())
				}
				if evtRobotState != nil {
					if evtRobotState.Status&uint32(vectorpb.RobotStatus_ROBOT_STATUS_IS_BUTTON_PRESSED) > 0 {
						isBusy = true
						gameOver = true
						return
					}
					if evtRobotState.Status&uint32(vectorpb.RobotStatus_ROBOT_STATUS_IS_PICKED_UP) > 0 ||
						evtRobotState.Status&uint32(vectorpb.RobotStatus_ROBOT_STATUS_IS_BEING_HELD) > 0 {
						println("I am being picked up.")
					}
					if evtRobotState.Status&uint32(vectorpb.RobotStatus_ROBOT_STATUS_IS_BUTTON_PRESSED) == 0 &&
						evtRobotState.TouchData.IsBeingTouched == true && !isBusy {
						isBusy = true
						go func() {
							println("I am being touched.")
							number := getRandomNumber(rnd)
							_ = sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/rattle.wav"))
							sdk_wrapper.WriteColoredText(number, 110, true, color.RGBA{255, 255, 255, 255}, 500, false)
							sdk_wrapper.SayText(number)
							time.Sleep(time.Duration(1000) * time.Millisecond)
							if c == MAX_NUMBERS {
								gameOver = true
							}
							isBusy = false
						}()
					}
				}
			}
		}
	}()

	for true {
		if gameOver {
			sdk_wrapper.SayText(getText("STR_BINGO_GAME_OVER"))
			break
		}
	}
	return
}

var numbers [MAX_NUMBERS]int
var c int = 0

func getRandomNumber(rnd *rand.Rand) string {
	n := 0
	for {
		n = rnd.Intn(MAX_NUMBERS + 1)
		if !containsNumber(numbers[:], n) {
			numbers[c] = n
			c++
			break
		}
	}
	return fmt.Sprintf("%d", n)
}

func containsNumber(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
