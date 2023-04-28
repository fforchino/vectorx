package intents

import (
	"context"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"image/color"
	"math"
	"time"
)

/**********************************************************************************************************************/
/*                                          FOLLOW CUBE                                                               */
/**********************************************************************************************************************/

func FollowCube_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"follow the cube"}
	utterances[LOCALE_ITALIAN] = []string{"segui il cubo"}
	utterances[LOCALE_SPANISH] = []string{"Sigue el cubo"}
	utterances[LOCALE_FRENCH] = []string{"Suivez le cube"}
	utterances[LOCALE_GERMAN] = []string{"Folgen Sie dem Würfel"}

	var intent = IntentDef{
		IntentName:            "extended_intent_follow_the_cube",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               followCube,
		OSKRTriggersUserInput: false,
	}
	*intentList = append(*intentList, intent)

	return nil
}

func followCube(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	doFollow(true)
	return returnIntent
}

func doFollow(useFx bool) {
	//s1 := rand.NewSource(time.Now().UnixNano())
	//rnd := rand.New(s1)

	sdk_wrapper.MoveHead(-3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)

	const WIDTH = sdk_wrapper.VECTOR_SCREEN_WIDTH
	const HEIGHT = sdk_wrapper.VECTOR_SCREEN_HEIGHT
	const MAX_SPEED = 50
	var oldSize float64 = 0

	sdk_wrapper.Robot.Conn.EnableMarkerDetection(
		context.Background(),
		&vectorpb.EnableMarkerDetectionRequest{Enable: true},
	)
	var cubeSize float64 = 0

	// Read input asynchronously
	go func() {
		for {
			evt := sdk_wrapper.WaitForEvent()
			if evt != nil {
				evtUserIntent := evt.GetUserIntent()
				evtObject := evt.GetObjectEvent()
				if evtUserIntent != nil {
					println(fmt.Sprintf("Received intent %d", evtUserIntent.IntentId))
					println(evtUserIntent.JsonData)
					println(evtUserIntent.String())
				}
				if evtObject != nil {
					appearedObject := evtObject.GetObjectAvailable()
					if appearedObject != nil {
						println("An object is available")
					}
					observerdObject := evtObject.GetRobotObservedObject()
					if observerdObject != nil && observerdObject.GetObjectType() == vectorpb.ObjectType_BLOCK_LIGHTCUBE1 {
						cubey := (observerdObject.ImgRect.YTopLeft + observerdObject.ImgRect.Height) / 2
						cubex := (observerdObject.ImgRect.XTopLeft + observerdObject.ImgRect.Width) / 2
						cubeWidth := observerdObject.ImgRect.Width
						cubeHeight := observerdObject.ImgRect.Height
						cubeCenterX := cubex + cubeWidth/2
						delta := WIDTH/2 - cubeCenterX
						// delta : w/2 = 1 : MAX_SPEED
						speed := delta * MAX_SPEED / (WIDTH / 2)
						cubeSize = math.Sqrt(float64(cubeWidth*cubeWidth + cubeHeight + cubeHeight))
						println(fmt.Sprintf("Spotted cube at %f,%f size: %f => Speed : %f", cubex, cubey, cubeSize, speed))

						if cubeSize < oldSize {
							if speed < 0 {
								sdk_wrapper.DriveWheelsForward(-1*speed, 0, -1*speed, 0)
							} else {
								sdk_wrapper.DriveWheelsForward(0, speed, 0, speed)
							}
							println("FORWARD")
							time.Sleep(time.Duration(500) * time.Millisecond)
							sdk_wrapper.DriveWheelsForward(0, 0, 0, 0)
						}
						oldSize = cubeSize
					}
				}
			}
		}
	}()

	for true {
	}
}
