package server

import (
	"encoding/json"
	"errors"
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

func getTopicId(r *http.Request) (string, error) {
	routeMatch := mux.RouteMatch{}
	if !mainRouter.Match(r, &routeMatch) {
		err := errors.New("Error: Route does not match any registered")
		return "", err
	}
	topicId := routeMatch.Vars["topicId"]
	return topicId, nil
}

func getSerializedMsg(topicId string, body string) ([]byte, error) {
	connectionMsg := message{
		Timestamp: time.Now().UTC().Format(time.RFC1123),
		Body:      body,
	}
	rawMsg, err := json.Marshal(connectionMsg)
	return rawMsg, err
}

func appendTopicConnection(topicId string, connection *websocket.Conn) {
	mainServer.connections[topicId] = append(mainServer.connections[topicId], connection)
}

func listenConnections(w http.ResponseWriter, r *http.Request) {
	topicId, err := getTopicId(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	payload, err := getSerializedMsg(topicId, fmt.Sprintf("Subscribed to topic: %s", topicId))
	if err != nil {
		fmt.Println(err)
		return
	}
	connection.WriteMessage(websocket.TextMessage, payload)
	appendTopicConnection(topicId, connection)
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			break
		}

		connection.WriteMessage(websocket.TextMessage, message)
		go messageHandler(topicId, []byte("foo"))
	}
}

func messageHandler(channelId string, message []byte) {
	fmt.Printf("%s: %s", channelId, message)
}

func InitServer() *server {
	mainServer = server{
		connections: map[string][]*websocket.Conn{},
		handler:     listenConnections,
		router:      mainRouter,
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
