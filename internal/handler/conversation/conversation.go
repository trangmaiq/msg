package conversation

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/trangmaiq/msg/internal/entity"
)

const (
	DefaultHandshakeTimeout = 5 * time.Second
	DefaultReadBufferSize   = 8 * 1024
)

var conversation = map[string]entity.Conversation{
	"1": {
		Forwarder: make(chan []byte),
	},
}

type JoinConversationInput struct {
	ConversationID string `uri:"id"`
}

func init() {
	go func() {
		for {
			select {
			case payload := <-conversation["1"].Forwarder:
				for _, sub := range conversation["1"].Subscribers {
					sub.Send <- payload
					log.Println("send to ", sub.ID)
				}
			}
		}
	}()
}

func Join(c *gin.Context) {
	var input JoinConversationInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

	upgrader := websocket.Upgrader{
		HandshakeTimeout: DefaultHandshakeTimeout,
		ReadBufferSize:   DefaultReadBufferSize,
		WriteBufferSize:  DefaultReadBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "join conversation failed",
		})
		return
	}

	subscriber := entity.Subscriber{
		ID:           uuid.NewString(),
		Conn:         conn,
		Conversation: conversation["1"].Forwarder,
		Send:         make(chan []byte),
	}
	subscribers := append(conversation["1"].Subscribers, subscriber)

	con := conversation["1"]
	con.Subscribers = subscribers
	conversation["1"] = con
	log.Println(conversation["1"])
	go subscriber.Read()
	subscriber.Write()
}
