package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type connectionsHandler func(w http.ResponseWriter, r *http.Request)

type server struct {
	connections map[string][]*websocket.Conn
	handler     connectionsHandler
	router      *mux.Router
	setup       *http.Server
}

type message struct {
	Timestamp string `json:"timestamp"`
	Body      string `json:"body"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var mainServer = server{}
var mainRouter = mux.NewRouter()

func listenConnections(w http.ResponseWriter, r *http.Request) {
	routeMatch := mux.RouteMatch{}
	if !mainRouter.Match(r, &routeMatch) {
		fmt.Println("Error: Route does not match any registered")
		return
	}
	topicId := routeMatch.Vars["topicId"]
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	connectionMsg := message{
		Timestamp: "2024-01-01 00:00:00",
		Body:      fmt.Sprintf("Connected to topic: %s", topicId),
	}
	rawMsg, err := json.Marshal(connectionMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	connection.WriteMessage(websocket.TextMessage, rawMsg)
	// for {
	// 	_, message, err := connection.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}

	// 	connection.WriteMessage(websocket.TextMessage, message)
	// 	go messageHandler(topicId, []byte("foo"))
	// }
}

func messageHandler(channelId string, message []byte) {
	fmt.Printf("%s: %s", channelId, message)
}

func InitServer() *server {
	mainServer = server{
		handler: listenConnections,
		router:  mainRouter,
		setup: &http.Server{
			Addr:         ":8000",
			Handler:      mainRouter,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}
	return &mainServer
}

func RegisterURIs() {
	mainRouter.HandleFunc("/tropico/{topicId}", mainServer.handler)
}

func (server *server) HandleConnections() {
	RegisterURIs()
	mainServer.setup.ListenAndServe()
}
