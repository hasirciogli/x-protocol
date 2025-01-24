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

	res, err := http.Post(fmt.Sprintf("http://%s:%d/calls", c.Host, c.Port), "application/json", bytes.NewBuffer([]byte(parsedBody)))
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

	// if response.Data is a string, convert it to json.RawMessage
	// if reflect.TypeOf(response.Data).Kind() == reflect.String {
	// 	response.Data = json.RawMessage(response.Data.(string))
	// }

	if appMode == "development" {
		fmt.Println("Call yanıtı alındı -> " + string(response.Data))
	}

	return XProtocolClientCallResponse{
		Success: true,
		Data:    json.RawMessage(response.Data),
		Error:   response.Error,
	}
}

func NewXProtocolClient(host string, port int) *XProtocolClient {
	return &XProtocolClient{
		Host: host,
		Port: port,
	}
}
