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
	"strings"
	vim_client "vectorx/pkg/vim-client"
)

var VIM_SERVER_URL = "localhost:8070"
var VIMEnabled = (os.Getenv("VIM_ENABLED") == "true")

var VIMDebug = true

func VIM_Register(intentList *[]IntentDef) error {
	registerSignUpToChat(intentList)
	registerLoginToChat(intentList)
	registerLogoutChat(intentList)
	registerSetChatTarget(intentList)
	registerQueryChatTarget(intentList)
	registerSendMessageToChat(intentList)

	addLocalizedString("STR_VIM_SIGN_UP_SUCCESSFUL", []string{"Signed up as %s1", "Registrato come %s1", "Registrado como %s1", "Enregistré comme %s1", "Aufgezeichnet wie %s1"})
	addLocalizedString("STR_VIM_ERROR_ALREADY_REGISTERED", []string{"Username %s1 is already registered", "Il nome %s1 è già in uso", "El nombre %s1 ya está registrado", "Le nom %s1 est déjà enregistré", "Benutzername %s1 ist bereits registriert"})
	addLocalizedString("STR_VIM_ERROR", []string{"Error", "Errore", "Error", "Erreur", "Fehler"})
	addLocalizedString("STR_VIM_LOGIN_SUCCESSFUL", []string{"Logged into chat service as %s1", "", "Acceso al servicio de chat como %s1", "Connecté au service de chat comme %s1", "Zugriff auf den Chat Service wie %s1"})
	addLocalizedString("STR_VIM_LOGOUT_SUCCESSFUL", []string{"Logout successful", "", "Desconexión realizada", "Déconnexion réussie", "Erfolgreich abmelden"})
	addLocalizedString("STR_VIM_MESSAGE_SENT", []string{"Message to %s1 sent", "Messaggio inviato a %s1", "Mensaje a %s1 enviado", "Message à %s1 envoyé", "Nachricht an %s1 gesendet"})
	addLocalizedString("STR_VIM_SEND_MESSAGE", []string{"say ", "invia ", "decir ", "dire ", "sagen "})
	addLocalizedString("STR_USER_SAYS_MESSAGE", []string{"%s1 says: %s2", "%s1 dice: %s2", "%s1 dice: %s2", "%s1 dit: %s2", "%s1 sagt: %s2"})
	addLocalizedString("STR_CHAT_TARGET_SET", []string{"chatting with %s1", "parliamo con %s1", "Chateando con %s1", "Discuter avec %s1", "Chatten mit %s1"})
	addLocalizedString("STR_CHAT_TARGET_UNKNOWN", []string{"not chatting with anyone", "non sto parlando con nessuno", "No estoy chateando con nadie", "Je ne parle à personne", "Nicht mit jemandem plaudern"})

	return nil
}

/**********************************************************************************************************************/
/*                                            SIGN UP TO CHAT                                                         */
/**********************************************************************************************************************/

func registerSignUpToChat(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"sign up to chat service"}
	utterances[LOCALE_ITALIAN] = []string{"registrati alla chat"}
	utterances[LOCALE_SPANISH] = []string{"Regístrese en el servicio de chat"}
	utterances[LOCALE_FRENCH] = []string{"Inscrivez-vous au service de chat"}
	utterances[LOCALE_GERMAN] = []string{"Registrieren Sie sich im Chat"}

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
	utterances[LOCALE_ITALIAN] = []string{"accedi alla chat"}
	utterances[LOCALE_SPANISH] = []string{"Conéctese al chat"}
	utterances[LOCALE_FRENCH] = []string{"Connectez-vous au chat"}
	utterances[LOCALE_GERMAN] = []string{"Verbindung zum Chat herstellen"}

	var intent = IntentDef{
		IntentName:            "extended_intent_vim_login",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               loginToChat,
		OSKRTriggersUserInput: nil,
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
	utterances[LOCALE_ITALIAN] = []string{"esci dalla chat"}
	utterances[LOCALE_SPANISH] = []string{"Salir del chat"}
	utterances[LOCALE_FRENCH] = []string{"Sortez du chat"}
	utterances[LOCALE_GERMAN] = []string{"Aus dem Chat rauskommen"}

	var intent = IntentDef{
		IntentName:            "extended_intent_vim_logout",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               logoutChat,
		OSKRTriggersUserInput: nil,
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
	utterances[LOCALE_SPANISH] = []string{"con quién estás hablando"}
	utterances[LOCALE_FRENCH] = []string{"A qui parles-tu"}
	utterances[LOCALE_GERMAN] = []string{"mit wem sprichst Du"}

	var intent = IntentDef{
		IntentName:            "extended_intent_vim_set_chat_target",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               queryChatTarget,
		OSKRTriggersUserInput: nil,
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
	utterances[LOCALE_SPANISH] = []string{"Habla con"}
	utterances[LOCALE_FRENCH] = []string{"Parler à"}
	utterances[LOCALE_GERMAN] = []string{"Sprechen Sie mit"}

	var intent = IntentDef{
		IntentName:            "extended_intent_vim_set_chat_target",
		Utterances:            utterances,
		Parameters:            []string{PARAMETER_CHAT_TARGET},
		Handler:               setChatTarget,
		OSKRTriggersUserInput: nil,
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
	utterances[LOCALE_SPANISH] = []string{"enviar"}
	utterances[LOCALE_FRENCH] = []string{"soumettre"}
	utterances[LOCALE_GERMAN] = []string{"senden"}

	var intent = IntentDef{
		IntentName:            "extended_intent_vim_message",
		Utterances:            utterances,
		Parameters:            []string{},
		Handler:               sendMessageToChat,
		OSKRTriggersUserInput: nil,
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
	Language  string `json:"language"`
}

type VIMUserInfoData struct {
	DisplayName string `json:"display_name"`
	UserId      string `json:"user_id"`
	IsHuman     bool   `json:"is_human"`
}

// Deprecated
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

// Deprecated
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

// Deprecated
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
			err := vim_client.SendMessageAndGo(VIM_SERVER_URL,
				myself.DisplayName,
				myself.UserId,
				other.DisplayName,
				other.UserId,
				sdk_wrapper.GetLanguageAndCountry(),
				botMessage)

			if err != nil {
				log.Fatal(err)
				println("FATAL: " + err.Error())
			}
			return err
		}
	}
	return errors.New("Unknown user")
}

func VIMAPIGetUserInfo(userName string) (VIMUserInfoData, error) {
	var info []VIMUserInfoData
	var dummy = VIMUserInfoData{}
	theUrl := "http://" + VIM_SERVER_URL + "/api/vim_list_targets"
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
			return dummy, err
		}
		for i := range info {
			if strings.ToLower(info[i].DisplayName) == strings.ToLower(userName) {
				return info[i], nil
			}
		}
	}
	return dummy, errors.New("Unknown user")
}

// Deprecated
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
	currentLanguage := sdk_wrapper.GetLanguage()
	currentLocale := sdk_wrapper.GetLocale()
	messageLanguage := strings.ToLower(strings.Split(msg.Language, "-")[0])
	messageLocale := msg.Language

	sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/vim/messageIn.wav"))

	if strings.HasPrefix(msg.Message, "*") && strings.HasSuffix(msg.Message, "*") {
		// Emoticon
		specialFileName := trimLeftChar(msg.Message)
		specialFileName = trimSuffix(specialFileName, "*")
		specialFileName = "images/vim/" + specialFileName + ".png"
		specialFileName = sdk_wrapper.GetDataPath(specialFileName)
		_, err := os.Stat(specialFileName)
		if err == nil {
			sdk_wrapper.MoveHead(3.0)
			sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
			sdk_wrapper.UseVectorEyeColorInImages(true)
			sdk_wrapper.SayText(getTextEx("STR_USER_SAYS_MESSAGE", []string{msg.From, ""}))
			go func() {
				sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/vim/specialMessageIn.wav"))
			}()
			sdk_wrapper.DisplayImage(specialFileName, 5000, true)
		} else {
			sdk_wrapper.SayText(getTextEx("STR_USER_SAYS_MESSAGE", []string{msg.From, msg.Message}))
		}
	} else {
		sdk_wrapper.SayText(getTextEx("STR_USER_SAYS_MESSAGE", []string{msg.From, ""}))
		if currentLanguage != messageLanguage {
			sdk_wrapper.SetLocale(messageLocale)
			sdk_wrapper.SetLanguage(messageLanguage)
		}
		sdk_wrapper.SayText(msg.Message)
		if currentLanguage != messageLanguage {
			sdk_wrapper.SetLocale(currentLocale)
			sdk_wrapper.SetLanguage(currentLanguage)
		}
	}
	sdk_wrapper.GetCustomSettings().LastChatMessageRead = msg.Id
	sdk_wrapper.SaveCustomSettings()
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

type ChatMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	FromId  string `json:"fromId"`
	ToId    string `json:"toId"`
	Message string `json:"msg"`
}
