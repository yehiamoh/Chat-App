package main

import (
	"fmt"

	"chat-app/pkg/database"
	"chat-app/pkg/routes"
	ws "chat-app/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize database connection
	db := database.DBConnection()

	routes.Init()
	AppRouter := routes.NewRouter(db)


	router:=gin.Default()

	router.GET("/auth", AppRouter.BeginAuthRoute)
	router.GET("/auth/callback", AppRouter.CallBackAuthRoute)
	router.GET("/logout", AppRouter.Logout)
	router.GET("/user",AppRouter.GetUserFromSession)

	router.GET("/ws",ws.WebScoketHandler)
	go ws.HandleMessages()




	fmt.Println("server started on :8000")
	if err:=router.Run(":8000");err!=nil{
		fmt.Println("Error in Starting Server:",err)
	}
}