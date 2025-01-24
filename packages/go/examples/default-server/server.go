package main

import (
	"encoding/json"

	"github.com/hasirciogli/x-protocol/packages/go/packages"
)

func Hello(payload json.RawMessage) packages.XProtocolCallResponse {
	return packages.XProtocolCallResponse{
		Success: true,
		Data:    json.RawMessage(payload),
		Error:   nil,
	}
}

func main() {
	server := packages.NewXProtocolServer("localhost", 8080)
	server.RegisterCall("hello", Hello)
	server.RegisterProxyChannel("p1", "localhost", 8090)

	server.Start()
}
