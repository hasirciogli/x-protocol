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
	server.RegisterCall("hello", func(payload json.RawMessage) json.RawMessage {
		var p HelloPayload
		p.Message = "hello"
		p.Name = "world"

		str, err := json.Marshal(p)
		if err != nil {
			return json.RawMessage(`{"error": "` + err.Error() + `"}`)
		}
		return json.RawMessage(str)
	})
	server.Start()
}
