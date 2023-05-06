package vim

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
	Message string `json:"msg"`
}

var c *websocket.Conn

func SendMessageAndGo(strFrom string, strFromId string, strTo string, strToId string, strMsg string) {
	openConnection("localhost:8080")
	chatMessage := ChatMessage{
		From:    strFrom,
		To:      strTo,
		FromId:  strFromId,
		ToId:    strToId,
		Message: strMsg,
	}
	// Send a message
	sendMessage(chatMessage)
	closeConnection()
}

func openConnection(host string) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	var err error
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
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
	return err
}
