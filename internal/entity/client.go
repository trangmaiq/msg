package entity

import (
	"bytes"
	"log"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	ID           string
	Conversation *Conversation
	Conn         *websocket.Conn
	Send         chan []byte
}

func (c *Client) Read() {
	defer func() {
		c.Conversation.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message unexpeted error: %v", err)
			}

			log.Printf("read message failed: %v", err)
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Conversation.Broadcast <- message
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("write message failed: %v", err)
				return
			}
		}
	}
}
