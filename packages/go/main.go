package main

import (
	"encoding/json"

	"github.com/hasirciogli/x-protocol/packages/go/packages"
)

type HelloPayload struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

func main() {
	server := packages.NewXProtocolServer("localhost", 8080)
	server.RegisterCall("hello", func(payload json.RawMessage) packages.XProtocolCallResponse {
		var p HelloPayload
		p.Message = "hello"
		p.Name = "world"

		json.Unmarshal(payload, &p)

		return packages.XProtocolCallResponse{
			Success: true,
			Data:    p,
		}
	})
	server.Start()
}
