package intents

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var VIM_SERVER_URL = os.Getenv("VIM_SERVER") // "https://www.wondergarden.app/VIM" //"http://192.168.43.65/VIM"
var VIMEnabled = (os.Getenv("VIM_ENABLED") == "true")

var VIMDebug = true

func VIM_Register(intentList *[]IntentDef) error {
	if VIMEnabled && VIM_SERVER_URL != "" {
		registerSignUpToChat(intentList)
		registerLoginToChat(intentList)
		registerLogoutChat(intentList)
		registerSetChatTarget(intentList)
		registerQueryChatTarget(intentList)
		registerSendMessageToChat(intentList)

		addLocalizedString("STR_VIM_SIGN_UP_SUCCESSFUL", []string{"Sign up successful as %s1", "", "", "", ""})
		addLocalizedString("STR_VIM_ERROR_ALREADY_REGISTERED", []string{"Username %s1 is already registered", "", "", "", ""})
		addLocalizedString("STR_VIM_ERROR", []string{"Error", "", "", "", ""})
		addLocalizedString("STR_VIM_LOGIN_SUCCESSFUL", []string{"Logged into chat service as %s1", "", "", "", ""})
		addLocalizedString("STR_VIM_LOGOUT_SUCCESSFUL", []string{"Logout successful", "", "", "", ""})
		addLocalizedString("STR_VIM_MESSAGE_SENT", []string{"Message to %s1 sent", "", "", "", ""})
		addLocalizedString("STR_VIM_SEND_MESSAGE", []string{"say ", "invia ", "", "", ""})
		addLocalizedString("STR_USER_SAYS_MESSAGE", []string{"%s1 says: %s2", "%s1 dice: %s2", "", "", ""})
		addLocalizedString("STR_CHAT_TARGET_SET", []string{"chatting with %s1", "parliamo con %s1", "", "", ""})
		addLocalizedString("STR_CHAT_TARGET_UNKNOWN", []string{"not chatting with anyone", "non sto parlando con nessuno", "", "", ""})
	}
	return nil
}

/**********************************************************************************************************************/
/*                                            SIGN UP TO CHAT                                                         */
/**********************************************************************************************************************/

func registerSignUpToChat(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"sign up to chat service"}
	utterances[LOCALE_ITALIAN] = []string{"registrati al servizio di chat"}

	var intent = IntentDef{
		IntentName: "extended_intent_vim_signup",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    signUpToChat,
	}
	*intentList = append(*intentList, intent)

	return nil
}

func signUpToChat(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	robotName := sdk_wrapper.GetRobotName()
	errorMessage := getText("STR_VIM_ERROR")

	if len(robotName) > 0 {
		err := VIMAPISignup(sdk_wrapper.GetRobotName(), sdk_wrapper.GetRobotSerial())
		if err == nil {
			sdk_wrapper.SayText(getTextEx("STR_VIM_SIGN_UP_SUCCESSFUL", []string{robotName}))
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
		} else if err.Error() == "Already registered" {
			sdk_wrapper.SayText(getTextEx("STR_VIM_ERROR_ALREADY_REGISTERED", []string{robotName}))
		} else {
			println(err.Error())
		}
	}
	if returnIntent == STANDARD_INTENT_IMPERATIVE_NEGATIVE {
		sdk_wrapper.SayText(errorMessage)
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                            LOGIN TO CHAT                                                           */
/**********************************************************************************************************************/

func registerLoginToChat(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"login to chat service"}
	utterances[LOCALE_ITALIAN] = []string{"attiva il servizio di chat"}

	var intent = IntentDef{
		IntentName: "extended_intent_vim_login",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    loginToChat,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func loginToChat(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	robotName := sdk_wrapper.GetRobotName()

	if len(robotName) > 0 {
		if VIMAPILogin(robotName, sdk_wrapper.GetRobotSerial()) == nil {
			sdk_wrapper.SayText(getTextEx("STR_VIM_LOGIN_SUCCESSFUL", []string{robotName}))
			sdk_wrapper.GetCustomSettings().LoggedInToChat = true
			sdk_wrapper.SaveCustomSettings()
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
		}
	}
	if returnIntent == STANDARD_INTENT_IMPERATIVE_NEGATIVE {
		sdk_wrapper.SayText(getText("STR_VIM_ERROR"))
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                            LOGOUT TO CHAT                                                           */
/**********************************************************************************************************************/

func registerLogoutChat(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"quit chat service"}
	utterances[LOCALE_ITALIAN] = []string{"esci dal servizio di chat"}

	var intent = IntentDef{
		IntentName: "extended_intent_vim_logout",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    logoutChat,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func logoutChat(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	robotName := sdk_wrapper.GetRobotName()

	if len(robotName) > 0 {
		if VIMAPILogout(robotName) == nil {
			sdk_wrapper.SayText(getTextEx("STR_VIM_LOGOUT_SUCCESSFUL", []string{robotName}))
			sdk_wrapper.GetCustomSettings().LoggedInToChat = false
			sdk_wrapper.SaveCustomSettings()
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
		}
	}
	if returnIntent == STANDARD_INTENT_IMPERATIVE_NEGATIVE {
		sdk_wrapper.SayText(getText("STR_VIM_ERROR"))
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                         QUERY CHAT TARGET                                                           */
/**********************************************************************************************************************/

func registerQueryChatTarget(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"who are you chatting with"}
	utterances[LOCALE_ITALIAN] = []string{"con chi stai parlando"}

	var intent = IntentDef{
		IntentName: "extended_intent_vim_set_chat_target",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    queryChatTarget,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func queryChatTarget(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	chatTargetName := sdk_wrapper.GetChatTarget()
	if len(chatTargetName) > 0 {
		sdk_wrapper.SayText(getTextEx("STR_CHAT_TARGET_SET", []string{chatTargetName}))
		returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	} else {
		sdk_wrapper.SayText(getText("STR_CHAT_TARGET_UNKNOWN"))
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                         SET CHAT TARGET                                                          */
/**********************************************************************************************************************/

func registerSetChatTarget(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"chat with"}
	utterances[LOCALE_ITALIAN] = []string{"parla con"}

	var intent = IntentDef{
		IntentName: "extended_intent_vim_set_chat_target",
		Utterances: utterances,
		Parameters: []string{PARAMETER_CHAT_TARGET},
		Handler:    setChatTarget,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func setChatTarget(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	if len(params.ChatTargetName) > 0 {
		sdk_wrapper.SetChatTarget(params.ChatTargetName)
		sdk_wrapper.SayText(getTextEx("STR_CHAT_TARGET_SET", []string{params.ChatTargetName}))
		returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                            SEND CHAT MESSAGE                                                       */
/**********************************************************************************************************************/

func registerSendMessageToChat(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"say"}
	utterances[LOCALE_ITALIAN] = []string{"invia"}

	var intent = IntentDef{
		IntentName: "extended_intent_vim_message",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    sendMessageToChat,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func sendMessageToChat(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_IMPERATIVE_NEGATIVE
	robotName := sdk_wrapper.GetRobotName()
	targetName := sdk_wrapper.GetChatTarget()

	if len(robotName) > 0 && len(targetName) > 0 {
		message := speechText[len(getText("STR_VIM_SEND_MESSAGE")):]
		if VIMAPISendMessageTo(targetName, message) == nil {
			sdk_wrapper.SayText(getTextEx("STR_VIM_MESSAGE_SENT", []string{targetName}))
			returnIntent = STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE
		}
	}
	if returnIntent == STANDARD_INTENT_IMPERATIVE_NEGATIVE {
		sdk_wrapper.SayText(getText("STR_VIM_ERROR"))
	}
	return returnIntent
}

/**********************************************************************************************************************/
/*                                            CHAT API                                                                */
/**********************************************************************************************************************/

type VIMChatMessage struct {
	Id        int32  `json:"id"`
	From      string `json:"from"`
	FromId    string `json:"from_id"`
	Message   string `json:"message"`
	Read      bool   `json:"read"`
	Timestamp int    `json:"timestamp"`
}

type VIMUserInfoData struct {
	DisplayName string `json:"display_name"`
	UserId      string `json:"user_id"`
	IsHuman     bool   `json:"is_human"`
}

func VIMAPISignup(robotName string, serialNo string) error {
	data := url.Values{
		"fname":     {"Vector"},
		"lname":     {robotName},
		"email":     {robotName + "@vectorx.org"},
		"password":  {serialNo},
		"serial_no": {serialNo},
	}

	resp, err := http.PostForm(VIM_SERVER_URL+"/php/signup.php", data)

	if err != nil {
		log.Fatal(err)
		println("FATAL: " + err.Error())
	} else {
		var responseText []byte
		responseText, err = ioutil.ReadAll(resp.Body)
		println("RESPONSE: " + string(responseText))
		if string(responseText) != "success" {
			err = errors.New(string(responseText))
		} else {
			err = nil
		}
	}
	return err
}

func VIMAPILogin(robotName string, serialNo string) error {
	data := url.Values{
		"email":    {robotName + "@vectorx.org"},
		"password": {serialNo},
	}

	resp, err := http.PostForm(VIM_SERVER_URL+"/php/login.php", data)

	if err != nil {
		log.Fatal(err)
	} else {
		var responseText []byte
		responseText, err = ioutil.ReadAll(resp.Body)
		if string(responseText) != "success" {
			err = errors.New(string(responseText))
		} else {
			err = nil
		}
	}
	return err
}

func VIMAPILogout(robotName string) error {
	data := url.Values{
		"email":    {robotName + "@vectorx.org"},
		"password": {robotName},
	}
	//TODO:
	resp, err := http.PostForm(VIM_SERVER_URL+"/php/logout-vector.php", data)

	if err != nil {
		log.Fatal(err)
	} else {
		var responseText []byte
		responseText, err = ioutil.ReadAll(resp.Body)
		if string(responseText) != "success" {
			err = errors.New(string(responseText))
		} else {
			err = nil
		}
	}
	return err
}

func VIMAPISetTarget(name string) {
	sdk_wrapper.SetChatTarget(name)
}

func VIMAPISendMessage(botMessage string) error {
	return VIMAPISendMessageTo(sdk_wrapper.GetChatTarget(), botMessage)
}

func VIMAPISendMessageTo(botTo string, botMessage string) error {
	robotName := sdk_wrapper.GetRobotName()
	if len(robotName) > 0 {
		myself, e1 := VIMAPIGetUserInfo(robotName)
		other, e2 := VIMAPIGetUserInfo(botTo)
		if e1 == nil && e2 == nil {
			println(fmt.Sprintf("Sending message '%s' from %s to %s", botMessage, myself.UserId, other.UserId))
			data := url.Values{
				"incoming_id": {other.UserId},
				"unique_id":   {myself.UserId},
				"message":     {botMessage},
			}

			resp, err := http.PostForm(VIM_SERVER_URL+"/php/insert-chat.php", data)

			if err != nil {
				log.Fatal(err)
				println("FATAL: " + err.Error())
			} else {
				var responseText []byte
				responseText, err = ioutil.ReadAll(resp.Body)
				println("RESPONSE: " + string(responseText))
				if string(responseText) != "" {
					err = errors.New("Cannot send message")
				}
			}
			return err
		}
	}
	return errors.New("Unknown user")
}

func VIMAPIGetUserInfo(userName string) (VIMUserInfoData, error) {
	var info []VIMUserInfoData
	theUrl := VIM_SERVER_URL + "/php/userInfo.php?displayName=" + userName
	resp, err := http.Get(theUrl)

	if err != nil {
		if VIMDebug {
			println("FATAL: " + err.Error())
		}
	} else {
		var responseText []byte
		responseText, err = ioutil.ReadAll(resp.Body)
		//println("RESPONSE: " + string(responseText))
		err = json.Unmarshal([]byte(responseText), &info)
		if err != nil {
			println(err.Error())
		}
		return info[0], err
	}
	return info[0], err

	return info[0], errors.New("Unknown user")
}

func VIMAPICheckMessages(robotSerialNo string, lastReadMessageId int32) ([]VIMChatMessage, error) {
	var arr []VIMChatMessage

	if len(robotSerialNo) > 0 {
		theUrl := fmt.Sprintf(VIM_SERVER_URL+"/php/get-chat-vector.php?unique_id=%s&last_message_id=%d",
			robotSerialNo,
			lastReadMessageId)
		resp, err := http.Get(theUrl)

		if err != nil {
			if VIMDebug {
				println("FATAL: " + err.Error())
			}
		} else {
			var responseText []byte
			responseText, err = ioutil.ReadAll(resp.Body)
			//println("RESPONSE: " + string(responseText))
			err = json.Unmarshal([]byte(responseText), &arr)
			if err != nil && VIMDebug {
				println(err.Error())
			}
			return arr, err
		}
		return arr, err
	}
	return arr, errors.New("Unknown user")
}

// <a href="https://www.freepik.com/free-vector/mixed-emoji-set_4159931.htm#query=emoticon&position=0&from_view=keyword">Image by rawpixel.com</a> on Freepik

func VIMAPIPlayMessage(msg VIMChatMessage) {
	sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/vim/messageIn.wav"))
	if msg.Message == ":-)" || msg.Message == ":)" {
		sdk_wrapper.MoveHead(3.0)
		sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
		sdk_wrapper.UseVectorEyeColorInImages(true)
		sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/vim/specialMessageIn.wav"))
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/vim/smile.png"), 5000, true)
	} else {
		sdk_wrapper.SayText(getTextEx("STR_USER_SAYS_MESSAGE", []string{msg.From, msg.Message}))
	}
	sdk_wrapper.GetCustomSettings().LastChatMessageRead = msg.Id
	sdk_wrapper.SaveCustomSettings()
}
