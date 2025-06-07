package main

import (
	"fmt"

	ws "chat-app/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	router:=gin.Default()
	router.GET("/ws",func(ctx *gin.Context) {
		ws.WebScoketHandler(ctx.Writer,ctx.Request)
	})
	go ws.HandleMessages()
	fmt.Println("server started on :8000")
	if err:=router.Run(":8000");err!=nil{
		fmt.Println("Error in Starting Server:",err)
	}
}