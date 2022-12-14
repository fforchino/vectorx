package intents

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"time"
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
		IntentName: "extended_intent_hello_world",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    playRockPaperScissors,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_LETS_PLAY", []string{"let's play!", "giochiamo!", "jugamos!", "jouons!", "spielen"})
	addLocalizedString("STR_ENOUGH", []string{"Ok, I think it's enough", "Penso che possa bastare", "Bien, creo que es suficiente", "bien, je pense que ça suffit!", "Gut, ich denke, es ist genug"})

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
	// Start http client
	var client *http.Client
	//setup a mocked http client.
	println("")
	println("Setup HTTP client")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}
		println(fmt.Sprintf("%s", b))
	}))
	defer ts.Close()
	client = ts.Client()

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)
	sdk_wrapper.EnableCameraStream()

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
		image := sdk_wrapper.ProcessCameraStream()
		if image != nil {
			var handInfo map[string]interface{}
			jsonData := sendImageToImageServer(client, &image)
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
				answer = "Sorry... I don't get it"
				break
			}

			if answer == "" {
				answer = "You put " + userMove + ". "

				switch win {
				case -1:
					answer = answer + "You win!"
					userScore++
					break
				case 1:
					answer = answer + "I win!"
					myScore++
					break
				default:
					answer = answer + "It's a draw!"
					break
				}
			}
			sdk_wrapper.SayText("I put " + myMove + "!")
			sdk_wrapper.SayText(answer)
			sdk_wrapper.WriteText(fmt.Sprintf("%d - %d", myScore, userScore), 64, true, 5000, true)
		}
	}
}

func sendImageToImageServer(client *http.Client, img *image.Image) string {
	//println("Encoding new frame")
	// Convert image to jpg and obtain the bytes
	var imageBuf bytes.Buffer
	_ = jpeg.Encode(&imageBuf, *img, nil)

	// Prepare the reader instances to encode
	values := map[string]io.Reader{
		"file": bytes.NewReader(imageBuf.Bytes()),
	}

	// Upload and get back the json response
	resp, err := Upload(client, "http://localhost:8090", values)
	if err != nil {
		println("Response error!")
		return ""
	}

	// Return json string
	//println("Response received: " + resp)
	return resp
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (response string, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return "", err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return "", err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	//println("Encoded data")
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		println("Error when performing HTTP request")
		return "", err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	//println("POSTing...")
	res, err := client.Do(req)
	if err != nil {
		println(err.Error())
		return "", err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		println(fmt.Errorf("bad status: %s", res.Status))
	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	return "", err
}
