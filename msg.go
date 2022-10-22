package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trangmaiq/msg/internal/handler/conversation"
)

func main() {
	r := gin.Default()

	r.GET("/conversation", conversation.Join)
	r.Run(":9000")
}
