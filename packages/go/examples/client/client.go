package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hasirciogli/x-protocol/packages/go/src/packages"
)

func main() {
	var testRequest = struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "admin",
		Password: "123456",
	}

	// var proxyChannelName = "p1" // blank or nil means no proxy channel is used

	jsonBytes, err := json.Marshal(testRequest)
	if err != nil {
		fmt.Println(err)
	}

	var client = packages.NewXProtocolClient("localhost", 8080)

	var response = client.Call(packages.XProtocolClientCallRequest{
		Name:    "hello", // function name (default: login for an example)
		Payload: json.RawMessage(jsonBytes),
		// ProxyChannelName: &proxyChannelName,
	})

	if response.Error != nil {
		fmt.Printf(*response.Error)
		return
	}

	text, err := json.Marshal(response.Data)
	if err != nil {
		fmt.Println(err)
	}

	// beautify json
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, text, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(prettyJSON.String())
}
