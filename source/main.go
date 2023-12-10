package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	baseUri     string
	connections map[string][]*websocket.Conn
	router      *mux.Router
	setup       *http.Server
}

func initServer(baseUri string) *Server {
	router := mux.NewRouter()
	return &Server{
		baseUri: baseUri,
		router:  router,
		setup: &http.Server{
			Handler:      router,
			Addr:         ":8000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}
}

func (server *Server) handleConnections() {
	server.router.HandleFunc("/ws/{channelId}", server.listenConnections)
	server.setup.ListenAndServe()
}

func (server *Server) listenConnections(w http.ResponseWriter, r *http.Request) {
	channelId := mux.Vars(r)["channelId"]
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			break
		}

		connection.WriteMessage(websocket.TextMessage, message)
		go messageHandler(channelId, message)
	}
	fmt.Println("Connection closed")
}

func messageHandler(channelId string, message []byte) {
	fmt.Printf("%s: %s", channelId, message)
}

func main() {
	server := initServer("/ws")
	server.handleConnections()
}
