package intents

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"
	opencv_ifc "vectorx/pkg/opencv-ifc"
)

/**********************************************************************************************************************/
/*                                          FOLLOW (INDEX) FINGER                                                     */
/**********************************************************************************************************************/

func Pong_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"follow the finger"}
	utterances[LOCALE_ITALIAN] = []string{"segui il dito"}
	utterances[LOCALE_SPANISH] = []string{"sigue el dedo"}
	utterances[LOCALE_FRENCH] = []string{"suivre le doigt"}
	utterances[LOCALE_GERMAN] = []string{"Folgen Sie dem Finger"}

	var intent = IntentDef{
		IntentName: "extended_intent_play_pong",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    playPong,
	}
	*intentList = append(*intentList, intent)

	return nil
}

func playPong(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	doPong(1000, true)
	return returnIntent
}

func doPong(numSteps int, useFx bool) {
	opencv_ifc.CreateClient()

	s1 := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(s1)

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)
	paddleFile, _ := os.Open(sdk_wrapper.GetDataPath("images/pong/paddle.png"))
	ballFile, _ := os.Open(sdk_wrapper.GetDataPath("images/pong/ball.png"))
	playFieldFile, _ := os.Open(sdk_wrapper.GetDataPath("images/pong/playfield.png"))
	paddle, _, _ := image.Decode(paddleFile)
	ball, _, _ := image.Decode(ballFile)
	playField, _, _ := image.Decode(playFieldFile)

	var scores []image.Image = []image.Image{}

	for i := 0; i <= 9; i++ {
		dFile, _ := os.Open(sdk_wrapper.GetDataPath(fmt.Sprintf("images/pong/digits_%d.png", i)))
		score, _, _ := image.Decode(dFile)
		scores = append(scores, score)
	}

	humanScore := 0
	vectorScore := 0
	const WIDTH = sdk_wrapper.VECTOR_SCREEN_WIDTH
	const HEIGHT = sdk_wrapper.VECTOR_SCREEN_HEIGHT
	PADDLE_WIDTH := paddle.Bounds().Dx()
	PADDLE_HEIGHT := paddle.Bounds().Dy()
	SPEED := 3
	VECTOR_PADDLE_SPEED := 3
	BALL_WIDTH := ball.Bounds().Dx()
	BALL_HEIGHT := ball.Bounds().Dy()

	// Human paddle coordinates
	humanPaddle := image.Point{0, HEIGHT / 2}
	// Vector paddle coordinates
	vectorPaddle := image.Point{WIDTH - PADDLE_WIDTH, HEIGHT / 2}
	// Ball coordinates
	ballObj := image.Point{WIDTH / 2, HEIGHT / 2}
	// Ball speed
	bSpeed := image.Point{X: -1 * SPEED, Y: 0}

	// Read input asynchronously
	go func() {
		for true {
			img, _ := sdk_wrapper.GetHiResCameraPicture()
			var handInfo map[string]interface{}
			jsonData := opencv_ifc.SendImageToImageServer(&img)
			//println("OpenCV server response: " + jsonData)
			json.Unmarshal([]byte(jsonData), &handInfo)
			index_x := int(handInfo["index_x"].(float64))
			if index_x != -1 {
				// Increment human paddle position
				humanPaddle.Y = HEIGHT * index_x / img.Bounds().Dx()
				if humanPaddle.Y < PADDLE_HEIGHT/2 {
					humanPaddle.Y = PADDLE_HEIGHT / 2
				} else if humanPaddle.Y > HEIGHT-PADDLE_HEIGHT/2 {
					humanPaddle.Y = HEIGHT - PADDLE_HEIGHT/2
				}
			}
		}
	}()

	for i := 0; i <= numSteps; i++ {
		fx := ""
		// Increment ball position
		ballObj.X += bSpeed.X
		ballObj.Y += bSpeed.Y

		// Increment Vector's paddle position, and check bounds
		if (vectorPaddle.Y + PADDLE_WIDTH) < ballObj.Y {
			vectorPaddle.Y += VECTOR_PADDLE_SPEED
		}
		if (vectorPaddle.Y + PADDLE_WIDTH) > ballObj.Y {
			vectorPaddle.Y -= VECTOR_PADDLE_SPEED
		}

		if vectorPaddle.Y < PADDLE_HEIGHT/2 {
			vectorPaddle.Y = PADDLE_HEIGHT / 2
		} else if vectorPaddle.Y > HEIGHT-PADDLE_HEIGHT/2 {
			vectorPaddle.Y = HEIGHT - PADDLE_HEIGHT/2
		}

		// Check bouncing
		if ballObj.X <= PADDLE_WIDTH && bSpeed.X < 0 {
			// Ball hits Human's wall
			if ballObj.Y >= humanPaddle.Y-PADDLE_HEIGHT/2 && ballObj.Y <= humanPaddle.Y+PADDLE_HEIGHT/2 {
				// Paddle hits the ball
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_ping.pcm")
				bSpeed.X = bSpeed.X * -1
				bSpeed.Y = (humanPaddle.Y - ballObj.Y) / 4
				println(fmt.Sprintf(">> HUMAN HITS THE BALL, new speed %d,%d", bSpeed.X, bSpeed.Y))
			} else {
				// Ball lost
				if vectorScore < 9 {
					vectorScore++
				}
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_out.pcm")
				ballObj.X = WIDTH / 2
				ballObj.Y = HEIGHT/2 + -1*HEIGHT/5 + rnd.Intn(HEIGHT/5*2)
				bSpeed.X = -1 * SPEED
				bSpeed.Y = -3 + rnd.Intn(6)
				println(fmt.Sprintf(">> BALL LOST, new speed %d,%d", bSpeed.X, bSpeed.Y))
			}
		} else if ballObj.X+BALL_WIDTH >= WIDTH-PADDLE_WIDTH && bSpeed.X > 0 {
			if ballObj.Y >= vectorPaddle.Y-PADDLE_HEIGHT/2 && ballObj.Y <= vectorPaddle.Y+PADDLE_HEIGHT/2 {
				// Paddle hits the ball
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_pong.pcm")
				bSpeed.X = bSpeed.X * -1
				bSpeed.Y = (vectorPaddle.Y - ballObj.Y) / 4
				println(fmt.Sprintf(">> VECTOR HITS THE BALL, new speed %d,%d", bSpeed.X, bSpeed.Y))
			} else {
				// Ball lost
				if humanScore < 9 {
					humanScore++
				}
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_out.pcm")
				ballObj.X = WIDTH / 2
				ballObj.Y = HEIGHT/2 + -1*HEIGHT/5 + rnd.Intn(HEIGHT/5*2)
				bSpeed.X = -1 * SPEED
				bSpeed.Y = -3 + rnd.Intn(6)
				println(fmt.Sprintf(">> BALL LOST, new speed %d,%d", bSpeed.X, bSpeed.Y))
			}
		} else if ballObj.Y <= BALL_HEIGHT || ballObj.Y+BALL_HEIGHT >= HEIGHT {
			// Ball hits top or bottom part of the screen, bounce back
			println(fmt.Sprintf(">> BALL BOUNCE, new speed %d,%d", bSpeed.X, bSpeed.Y))
			fx = sdk_wrapper.GetDataPath("audio/pong/ball_bounce.pcm")
			bSpeed.Y = bSpeed.Y * -1
		}

		// Draw
		dc := gg.NewContext(WIDTH, HEIGHT)
		dc.DrawImage(playField, 0, 0)
		dc.DrawImage(paddle, humanPaddle.X, humanPaddle.Y-PADDLE_HEIGHT/2)
		dc.DrawImage(paddle, vectorPaddle.X, vectorPaddle.Y-PADDLE_HEIGHT/2)
		dc.DrawImage(ball, ballObj.X, ballObj.Y)
		dc.DrawImage(scores[humanScore], WIDTH/4-5, 0)
		dc.DrawImage(scores[vectorScore], WIDTH/4*3-5, 0)

		buf := new(bytes.Buffer)
		bitmap := convertPixelsToRawBitmap(dc.Image(), 100)
		for _, ui := range bitmap {
			binary.Write(buf, binary.LittleEndian, ui)
		}
		_, _ = sdk_wrapper.Robot.Conn.DisplayFaceImageRGB(
			context.Background(),
			&vectorpb.DisplayFaceImageRGBRequest{
				FaceData:         buf.Bytes(),
				DurationMs:       100,
				InterruptRunning: true,
			},
		)

		// Play FX
		if useFx && fx != "" {
			go func() {
				sdk_wrapper.PlaySound(fx)
			}()
		}
		//println(fmt.Sprintf("Step %d/%d. User pos @ %d,%d, ball pos %d,%d ballspeed @ %d,%d", i, numSteps, humanPaddle.X, humanPaddle.Y, ballObj.X, ballObj.Y, bSpeed.X, bSpeed.Y))
	}
}

func convertPixesTo16BitRGB(r uint32, g uint32, b uint32, a uint32, opacityPercentage uint16) uint16 {
	R, G, B := uint16(r/257), uint16(g/8193), uint16(b/257)

	R = R * opacityPercentage / 100
	G = G * opacityPercentage / 100
	B = B * opacityPercentage / 100

	//The format appears to be: 000bbbbbrrrrrggg

	var Br uint16 = (uint16(B & 0xF8)) << 5 // 5 bits for blue  [8..12]
	var Rr uint16 = (uint16(R & 0xF8))      // 5 bits for red   [3..7]
	var Gr uint16 = (uint16(G))             // 3 bits for green [0..2]

	out := uint16(Br | Rr | Gr)
	//println(fmt.Sprintf("%d,%d,%d -> R: %016b G: %016b B: %016b = %016b", R, G, B, Rr, Gr, Br, out))
	return out
}

func convertPixelsToRawBitmap(image image.Image, opacityPercentage int) []uint16 {
	imgHeight, imgWidth := image.Bounds().Max.Y, image.Bounds().Max.X
	bitmap := make([]uint16, imgWidth*imgHeight)

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			vectorEyes := sdk_wrapper.GetEyeColor()
			vR := uint32(vectorEyes.R) * 255
			vG := uint32(vectorEyes.G) * 255
			vB := uint32(vectorEyes.B) * 255

			r = r * vR / 0xffff
			g = g * vG / 0xffff
			b = b * vB / 0xffff
			bitmap[(y)*imgWidth+(x)] = convertPixesTo16BitRGB(r, g, b, a, uint16(opacityPercentage))
		}
	}
	return bitmap
}