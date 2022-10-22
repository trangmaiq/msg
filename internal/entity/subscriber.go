package entity

import (
	"log"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	ID           string
	Conn         *websocket.Conn
	Conversation chan<- []byte
	Send         chan []byte
}

func (s *Subscriber) Read() {
	defer s.Conn.Close()
	for {
		log.Println("read message from ", s.ID)
		_, payload, err := s.Conn.ReadMessage()
		if err != nil {
			log.Printf("read message failed: %v", err)
			continue
		}

		s.Conversation <- payload
	}
}

func (s *Subscriber) Write() {
	defer s.Conn.Close()

	for msg := range s.Send {
		log.Println("write message to ", s.ID)
		err := s.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("write message failed: %w", err)
			continue
		}
	}
}
