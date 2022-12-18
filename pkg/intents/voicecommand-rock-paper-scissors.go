package intents

import (
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"
	opencv_ifc "vectorx/pkg/opencv-ifc"
)

/**********************************************************************************************************************/
/*                                    VECTOR PLAYS ROCK PAPER SCISSORS                                                */
/**********************************************************************************************************************/

func RockPaperScissors_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"let's play a new game"}
	utterances[LOCALE_ITALIAN] = []string{"giochiamo a morra cinese"}
	utterances[LOCALE_SPANISH] = []string{"jugamos a piedra papel o tijera"}
	utterances[LOCALE_FRENCH] = []string{"jouons à pierre papier ciseaux"}
	utterances[LOCALE_GERMAN] = []string{"spielen schere stein papier"}

	var intent = IntentDef{
		IntentName: "extended_intent_rock_paper_scissors",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    playRockPaperScissors,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_LETS_PLAY", []string{"let's play!", "giochiamo!", "jugamos!", "jouons!", "spielen!"})
	addLocalizedString("STR_ENOUGH", []string{"Ok, I think it's enough", "Penso che possa bastare", "Bien, creo que es suficiente", "bien, je pense que ça suffit!", "Gut, ich denke, es ist genug"})
	addLocalizedString("STR_I_WIN", []string{"I win", "ho vinto", "yo gano", "je gagne", "ich gewinne"})
	addLocalizedString("STR_YOU_WIN", []string{"you win", "hai vinto", "tú ganas", "Vous gagnez", "du gewinnst"})
	addLocalizedString("STR_ITS_A_DRAW", []string{"it's a draw", "pareggio", "es un empate", "C'est un match nul", "es ist eine Zeichnung"})
	addLocalizedString("STR_YOU_PUT", []string{"you put", "hai messo", "pones", "tu mets", "du legst"})
	addLocalizedString("STR_I_PUT", []string{"I put", "ho messo", "puse", "je mets", "ich lege"})
	addLocalizedString("STR_SORRY_I_DONT_GET_IT", []string{"sorry, I don't get it", "scusa non l'ho capita", "Lo siento, no lo entiendo", "Désolé, je ne comprends pas", "Entschuldigung, ich verstehe es nicht"})
	addLocalizedString("STR_ROCK", []string{"rock", "roccia", "roca", "rock", "Felsen"})
	addLocalizedString("STR_PAPER", []string{"paper", "carta", "papel", "papier", "Papier"})
	addLocalizedString("STR_SCISSORS", []string{"scissors", "forbici", "tijeras", "les ciseaux", "Schere"})

	return nil
}

func playRockPaperScissors(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	sdk_wrapper.SayText(getText("STR_LETS_PLAY"))
	playGame(10)
	sdk_wrapper.SayText(getText("STR_ENOUGH"))
	return returnIntent
}

func playGame(numSteps int) {
	opencv_ifc.CreateClient()

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)

	myScore := 0
	userScore := 0
	options := [3]string{
		"rock",
		"paper",
		"scissors",
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i <= numSteps; i++ {
		sdk_wrapper.SayText("1, 2, 3!")

		myMove := options[r1.Intn(len(options))]
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/rps/"+myMove+".png"), 5000, false)
		fName := sdk_wrapper.GetTemporaryFilename("rps", "jpg", true)
		sdk_wrapper.SaveHiResCameraPicture(fName)
		f, err := os.Open(fName)

		if err == nil {
			defer f.Close()
			image, _, _ := image.Decode(f)

			var handInfo map[string]interface{}
			jsonData := opencv_ifc.SendImageToImageServer(&image)
			println("OpenCV server response: " + jsonData)
			json.Unmarshal([]byte(jsonData), &handInfo)
			numFingers := -1
			numFingers = int(handInfo["raisedfingers"].(float64))
			win := 0
			answer := ""
			userMove := ""

			println(fmt.Sprintf("num fingers %d", numFingers))

			switch numFingers {
			case 0:
				// User plays "rock"
				userMove = "rock"
				if myMove == "paper" {
					win = 1
				} else if myMove == "scissors" {
					win = -1
				}
				break
			case 2:
				// User plays "scissors"
				userMove = "scissors"
				if myMove == "rock" {
					win = 1
				} else if myMove == "paper" {
					win = -1
				}
				break
			case 5:
				// User plays "paper"
				userMove = "paper"
				if myMove == "scissors" {
					win = 1
				} else if myMove == "rock" {
					win = -1
				}
				break
			default:
				answer = getText("STR_SORRY_I_DONT_GET_IT")
				break
			}

			userMoveLocalized := localizeMove(userMove)
			myMoveLocalized := localizeMove(myMove)

			if answer == "" {
				answer = getText("STR_YOU_PUT") + " " + userMoveLocalized + ". "

				switch win {
				case -1:
					answer = answer + getText("STR_YOU_WIN") + "!"
					userScore++
					break
				case 1:
					answer = answer + getText("STR_I_WIN") + "!"
					myScore++
					break
				default:
					answer = answer + getText("STR_ITS_A_DRAW") + "!"
					break
				}
			}
			sdk_wrapper.SayText(getText("STR_I_PUT") + " " + myMoveLocalized + "!")
			sdk_wrapper.SayText(answer)
			sdk_wrapper.WriteText(fmt.Sprintf("%d - %d", myScore, userScore), 64, true, 5000, true)
		}
	}
}

func localizeMove(move string) string {
	locMove := move
	switch move {
	case "rock":
		locMove = getText("STR_ROCK")
		break
	case "paper":
		locMove = getText("STR_PAPER")
		break
	case "scissors":
		locMove = getText("STR_SCISSORS")
		break
	}
	return locMove
}
