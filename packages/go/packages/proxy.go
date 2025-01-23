package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type XProtocolProxyService struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	ProxyMode bool   `json:"proxy_mode"`
}

type XProtocolProxyServiceClient struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type XProtocolProxyChannel struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewXProtocolProxyServiceClient(host string, port int) *XProtocolProxyServiceClient {
	return &XProtocolProxyServiceClient{
		Host: host,
		Port: port,
	}
}

func (s *XProtocolProxyServiceClient) Call(name string, payload XProtocolCallRequest) (XProtocolCallResponse, error) {
	if payload.FromProxyChannel != nil {
		payload.FromProxyChannel.Name = name
		payload.FromProxyChannel.Host = s.Host
		payload.FromProxyChannel.Port = s.Port
	}

	bodyTextJsonBytes, err := json.Marshal(payload)
	if err != nil {
		return XProtocolCallResponse{}, err
	}

	res, err := http.Post(fmt.Sprintf("http://%s:%d/call/%s", s.Host, s.Port, name), "application/json", bytes.NewBuffer(bodyTextJsonBytes))
	if err != nil {
		return XProtocolCallResponse{}, err
	}
	defer res.Body.Close()

	var response XProtocolCallResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return XProtocolCallResponse{}, err
	}

	return response, nil
}
