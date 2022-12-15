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
	opencv_ifc "vectorx/pkg/opencv-ifc"
)

/**********************************************************************************************************************/
/*                                          FOLLOW (INDEX) FINGER                                                     */
/**********************************************************************************************************************/

func FollowFinger_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"follow the finger"}
	utterances[LOCALE_ITALIAN] = []string{"segui il dito"}
	utterances[LOCALE_SPANISH] = []string{"sigue el dedo"}
	utterances[LOCALE_FRENCH] = []string{"suivre le doigt"}
	utterances[LOCALE_GERMAN] = []string{"Folgen Sie dem Finger"}

	var intent = IntentDef{
		IntentName: "extended_intent_follow_finger",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    followFinger,
	}
	*intentList = append(*intentList, intent)

	return nil
}

func followFinger(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	doFollow(100)
	return returnIntent
}

func doFollow(numSteps int) {
	opencv_ifc.CreateClient()

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)

	x := sdk_wrapper.VECTOR_SCREEN_WIDTH / 2
	y := sdk_wrapper.VECTOR_SCREEN_HEIGHT / 2

	for i := 0; i <= numSteps; i++ {
		println(fmt.Sprintf("Step %d/%d", i, numSteps))
		img, err := sdk_wrapper.GetHiResCameraPicture()
		if err == nil {
			var handInfo map[string]interface{}
			jsonData := opencv_ifc.SendImageToImageServer(&img)
			println("OpenCV server response: " + jsonData)
			json.Unmarshal([]byte(jsonData), &handInfo)
			index_x := int(handInfo["index_x"].(float64))
			index_y := int(handInfo["index_y"].(float64))
			if index_x != -1 {
				x = sdk_wrapper.VECTOR_SCREEN_WIDTH * index_x / img.Bounds().Dx()
			}
			if index_y != -1 {
				y = sdk_wrapper.VECTOR_SCREEN_HEIGHT * index_y / img.Bounds().Dy()
			}
			bgImage := image.NewRGBA(image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: sdk_wrapper.VECTOR_SCREEN_WIDTH, Y: sdk_wrapper.VECTOR_SCREEN_HEIGHT},
			})
			dc := gg.NewContext(sdk_wrapper.VECTOR_SCREEN_WIDTH, sdk_wrapper.VECTOR_SCREEN_HEIGHT)
			dc.DrawImage(bgImage, 0, 0)
			dc.SetColor(sdk_wrapper.GetEyeColor())
			dc.DrawCircle(float64(x), float64(y), 10)

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
		}
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
			bitmap[(y)*imgWidth+(x)] = convertPixesTo16BitRGB(r, g, b, a, uint16(opacityPercentage))
		}
	}
	return bitmap
}
