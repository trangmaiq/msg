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

var conversations = map[string]*entity.Conversation{
	"1": entity.NewConversation(),
}

type JoinConversationInput struct {
	ConversationID string `uri:"id"`
}

func init() {
	for _, c := range conversations {
		go c.Start()
	}
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

	client := entity.Client{
		ID:           uuid.NewString(),
		Conn:         conn,
		Send:         make(chan []byte),
		Conversation: conversations["1"],
	}
	client.Conversation.Register <- &client
	go client.Read()
	go client.Write()
}
