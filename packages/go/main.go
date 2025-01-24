package main

import (
	"encoding/json"

	"github.com/hasirciogli/x-protocol/packages/go/packages"
)

type HelloPayload struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

func Hello(payload json.RawMessage) packages.XProtocolCallResponse {
	var p HelloPayload
	p.Message = "hello"
	p.Name = "world"

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		errString := err.Error()
		return packages.XProtocolCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}

	return packages.XProtocolCallResponse{
		Success: true,
		Data:    json.RawMessage(jsonBytes),
		Error:   nil,
	}
}

func main() {
	server := packages.NewXProtocolServer("localhost", 8080)
	server.RegisterCall("hello2", Hello)
	server.RegisterCall("hello", Hello)
	server.RegisterProxyChannel("test", "localhost", 8090)

	server.Start()
}
