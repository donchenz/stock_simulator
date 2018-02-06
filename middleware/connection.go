package main

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type connection struct {
	ws *websocket.Conn

	send chan []byte
}

type Message struct {
	Seq int
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			break
		}
		h.heartbeatRefresh <- c
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}





