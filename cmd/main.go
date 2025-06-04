package main

import (
	"fmt"
	"net/http"
	"sync"

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

func wsHandler(w http.ResponseWriter,r*http.Request){

	//upgrade the upcoming request from HTTP to a websocket
	conn,err:=upgrader.Upgrade(w,r,nil)
	if err!=nil{
		fmt.Printf("error in upgrading :%v",err)
		return
	}
	// using go routine to manage multiple websocket simultaneousl
	go handleConnection(conn)

}
func handleConnection(conn *websocket.Conn){
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
        fmt.Printf("Message Received: %s\n", string(message))
        boradcast <- message  // Send to broadcast channel
	}
}

func handleMessages(){
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
func main() {
	http.HandleFunc("/ws",wsHandler)
	go handleMessages()
	fmt.Println("WebSocket server started on :8000")
	if err:=http.ListenAndServe(":8000",nil);err!=nil{
		fmt.Println("Error in Starting Server:",err)
	}
}