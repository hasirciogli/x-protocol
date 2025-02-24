package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type XProtocolClient struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type XProtocolClientCallResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *string         `json:"error"`
}

type XProtocolClientCallRequest struct {
	ProxyChannelName *string         `json:"proxy_channel_name"`
	Name             string          `json:"name"`
	Payload          json.RawMessage `json:"payload"`
	Token            *string         `json:"token"`
}

func (c *XProtocolClient) Call(xprotoCallRequest XProtocolClientCallRequest) XProtocolClientCallResponse {
	bodyTextJsonBytes, err := json.Marshal(xprotoCallRequest)
	if err != nil {
		errString := err.Error()
		return XProtocolClientCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}

	appMode := os.Getenv("APP_MODE")
	if appMode == "development" {
		fmt.Println("Call isteği gönderildi -> " + xprotoCallRequest.Name)
	}

	parsedBody := string(bodyTextJsonBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/calls", c.Host, c.Port), bytes.NewBuffer([]byte(parsedBody)))
	if err != nil {
		errString := err.Error()
		return XProtocolClientCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}

	if xprotoCallRequest.Token != nil {
		req.Header.Set("Authorization", "Bearer "+*xprotoCallRequest.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		errString := err.Error()
		return XProtocolClientCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errString := err.Error()
		return XProtocolClientCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}

	if res.StatusCode != http.StatusOK {
		errString := string(body)
		return XProtocolClientCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}

	var response XProtocolCallResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		errString := err.Error()
		return XProtocolClientCallResponse{
			Success: false,
			Data:    nil,
			Error:   &errString,
		}
	}

	if appMode == "development" {
		fmt.Println("Call yanıtı alındı -> " + string(response.Data))
	}

	return XProtocolClientCallResponse(response)
}

func NewXProtocolClient(host string, port int) *XProtocolClient {
	return &XProtocolClient{
		Host: host,
		Port: port,
	}
}
