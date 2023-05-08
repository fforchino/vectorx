package vim_client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type ChatMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	FromId  string `json:"fromId"`
	ToId    string `json:"toId"`
	Lang    string `json:"lang"`
	Message string `json:"msg"`
}

var c *websocket.Conn

func SendMessageAndGo(strServerUrl string, strFrom string, strFromId string, strTo string, strToId string, strLang string, strMsg string) error {
	err := openConnection(strServerUrl)
	if err == nil {
		chatMessage := ChatMessage{
			From:    strFrom,
			To:      strTo,
			FromId:  strFromId,
			ToId:    strToId,
			Lang:    strLang,
			Message: strMsg,
		}
		// Send a message
		err = sendMessage(chatMessage)
		err = closeConnection()
	}
	return err
}

func openConnection(host string) error {
	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	var err error
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return err
}

func sendMessage(msg ChatMessage) error {
	payload, err := json.Marshal(msg)
	if err == nil {
		err = c.WriteMessage(websocket.TextMessage, payload)
	}
	if err != nil {
		log.Println("write:", err)
	}
	return err
}

func closeConnection() error {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
	}
	defer c.Close()
	return err
}
