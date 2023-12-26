package main

import (
	server "github.com/luislve17/tropico/server"
)

func main() {
	server := server.InitServer()
	server.HandleConnections()
}
