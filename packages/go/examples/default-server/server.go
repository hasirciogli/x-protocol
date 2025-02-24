package main

import (
	"encoding/json"

	"github.com/hasirciogli/x-protocol/packages/go/src/packages"
)

func Hello(payload json.RawMessage) packages.XProtocolCallResponse {
	return packages.XProtocolCallResponse{
		Success: true,
		Data:    json.RawMessage(`{"message": "Hello, World!"}`),
		Error:   nil,
	}
}

func main() {
	server := packages.NewXProtocolServer("localhost", 8080)
	// server.RegisterAuthCallback(func(authHeader string) bool {
	// 	return authHeader == "Bearer 123456"
	// })

	server.RegisterCall("hello", Hello)
	server.Start()
}
