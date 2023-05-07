// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim_server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"vectorx/pkg/intents"
	"vectorx/pkg/vim-client"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
		forwardMessageToVector(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func forwardMessageToVector(message []byte) error {
	var chatMessage vim_client.ChatMessage
	err := json.Unmarshal(message, &chatMessage)
	if err != nil {
		return err
	}
	serial := chatMessage.ToId

	/*if intents.VIMEnabled*/
	{
		println("VIM Enabled")
		/*isChatty := isBotInChatMood(serial)
		if isChatty*/{
			var ctx = context.Background()
			var start = make(chan bool)
			var stop = make(chan bool)
			intents.RegisterIntents()
			// See if we can connect to a bot with this ESN
			err := sdk_wrapper.InitSDKForWirepod(serial)
			if err == nil {
				go func() {
					_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
				}()
				select {
				case <-start:
					var vcm intents.VIMChatMessage
					vcm.Message = chatMessage.Message
					vcm.From = chatMessage.From
					vcm.FromId = chatMessage.FromId
					vcm.Read = false
					vcm.Timestamp = int(time.Now().UnixMilli())
					println(fmt.Sprintf("[%d] New message from %s: %s", vcm.Timestamp, vcm.From, vcm.Message))
					intents.VIMAPIPlayMessage(vcm)
					stop <- true
				}
			}
		}
	}
	return nil
}

func isBotInChatMood(serial string) bool {
	serial = strings.ToLower(serial)
	// Peek into the given vector custom settings and read the value
	customSettingsPath := filepath.Join(os.Getenv("VECTORX_HOME"), "vectorfs")
	customSettingsPath = filepath.Join(customSettingsPath, "nvm")
	customSettingsPath = filepath.Join(customSettingsPath, serial)
	customSettingsPath = filepath.Join(customSettingsPath, "custom_settings.json")

	println(customSettingsPath)
	botCustomSettingsJSONFile, err := os.ReadFile(customSettingsPath)
	if err == nil {
		var botCustomSettings sdk_wrapper.CustomSettings
		err := json.Unmarshal(botCustomSettingsJSONFile, &botCustomSettings)
		if err == nil {
			//println("OK")
			return botCustomSettings.LoggedInToChat
		}
	}
	return false
}
