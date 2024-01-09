package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func assertEqual(t *testing.T, a any, b any) {
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	if a != b || typeA != typeB {
		t.Fatalf("Error: Not equal:\n'%s(%s)'\n'%s(%s)'", a, typeA, b, typeB)
	}
}

func initTestServer() (*server, *httptest.Server) {
	serverSetup := InitServer()
	testServer := httptest.NewServer(http.HandlerFunc(serverSetup.handler))
	RegisterURIs()

	return serverSetup, testServer
}

func connectToTopic(testCtx *testing.T, testServer *httptest.Server, topic string) *websocket.Conn {
	uri := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/tropico/" + topic
	testClient, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		testCtx.Fatalf("%v", err)
	}
	return testClient
}

func receiveMessage(t *testing.T, ws *websocket.Conn) []byte {
	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return p
}

func TestServer_ClientReceivesMessageWhenConnectingToValidTopic(t *testing.T) {
	serverSetup, testServer := initTestServer()
	defer testServer.Close()
	ws := connectToTopic(t, testServer, "test-topic")
	defer ws.Close()

	recvMessage := receiveMessage(t, ws)
	response := message{}
	json.Unmarshal(recvMessage, &response)
	assertEqual(t, response.Body, "Subscribed to topic: test-topic")
	assertEqual(t, len(serverSetup.connections), 1)
}
