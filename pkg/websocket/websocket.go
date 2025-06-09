package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections
var upgrader=websocket.Upgrader{
	CheckOrigin: func (r *http.Request)bool  {
		return true
	},
}

var clients=make(map[*websocket.Conn]bool) //connected clients
var boradcast=make(chan []byte) //Broadcast channel
var mutex=&sync.Mutex{} //protected clients map

func WebScoketHandler(c *gin.Context){

	//upgrade the upcoming request from HTTP to a websocket
	conn,err:=upgrader.Upgrade(c.Writer,c.Request,nil)
	if err!=nil{
		fmt.Printf("error in upgrading :%v",err)
		return
	}
	// using go routine to manage multiple websocket simultaneousl
	go handleConnection(conn)

}
func handleConnection(conn *websocket.Conn){
	remoteAddr:=conn.RemoteAddr().String()
	mutex.Lock()
	clients[conn]=true
	mutex.Unlock()

	defer func(){
		mutex.Lock()
        delete(clients, conn)
        mutex.Unlock()
        conn.Close()
	}()
	for{
     _, message, err := conn.ReadMessage()
    	if err != nil {
            fmt.Printf("Error reading message: %v\n", err)
            break
        }
        fmt.Printf("Message Received: %s\nFrom Connection :%v\n", string(message),remoteAddr)
        boradcast <- message  // Send to broadcast channel
	}
}

func HandleMessages(){
	for{

		message:= <-boradcast
		mutex.Lock() 
		for client:=range clients{
			if err:= client.WriteMessage(websocket.TextMessage,message);err!=nil{
				fmt.Printf("Error broadcasting to client: %v\n", err)
				client.Close()
				delete(clients,client)
				continue
			}
		}

		mutex.Unlock()
	}
}