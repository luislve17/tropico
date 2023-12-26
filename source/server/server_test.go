package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func initTestServer() *httptest.Server {
	server := InitServer()
	testServer := httptest.NewServer(http.HandlerFunc(server.handler))

	return testServer
}

func connectToTopic(testCtx *testing.T, testServer *httptest.Server, topic string) *websocket.Conn {
	uri := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/tropico/" + topic
	testClient, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		testCtx.Fatalf("%v", err)
	}
	return testClient
}

func TestServer_ClientSubscribesToTopic(t *testing.T) {
	testServer := initTestServer()
	defer testServer.Close()
	ws := connectToTopic(t, testServer, "test-topic")
	defer ws.Close()

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	response := message{}
	json.Unmarshal(p, &response)
	if response.Body != "Connected to topic: test-topic" {
		t.Fatalf("Unexpected message: '%s'", response.Body)
	}
}
