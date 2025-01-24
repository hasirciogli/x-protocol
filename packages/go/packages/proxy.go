package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type XProtocolProxyService struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	ProxyMode bool   `json:"proxy_mode"`
}

type XProtocolProxyChannel struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type XProtocolProxyServiceClient struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type XProtocolProxyCallResponse struct {
	Success          bool    `json:"success"`
	Data             string  `json:"data"`
	Error            *string `json:"error"`
	ProxyStatus      *int    `json:"proxy_status"`
	ProxyError       *string `json:"proxy_error"`
	ProxyServerError bool    `json:"proxy_server_error"`
}

func NewXProtocolProxyServiceClient(host string, port int) *XProtocolProxyServiceClient {
	return &XProtocolProxyServiceClient{
		Host: host,
		Port: port,
	}
}

func (s *XProtocolProxyChannel) Call(name string, xprotoCallRequest XProtocolCallRequest) XProtocolProxyCallResponse {
	bodyTextJsonBytes, err := json.Marshal(XProtocolCallRequest{
		Name:    name,
		Payload: xprotoCallRequest.Payload,
		FromProxyChannel: &XProtocolProxyChannel{
			Name: s.Name,
			Host: s.Host,
			Port: s.Port,
		},
	})
	if err != nil {
		errString := err.Error()
		return XProtocolProxyCallResponse{
			Success:          false,
			Data:             "",
			Error:            &errString,
			ProxyStatus:      nil,
			ProxyError:       nil,
			ProxyServerError: false,
		}
	}

	parsedBody := string(bodyTextJsonBytes)

	res, err := http.Post(fmt.Sprintf("http://%s:%d/calls", s.Host, s.Port), "application/json", bytes.NewBuffer([]byte(parsedBody)))
	if err != nil {
		errString := err.Error()
		return XProtocolProxyCallResponse{
			Success:          false,
			Data:             "",
			Error:            &errString,
			ProxyServerError: false,
			ProxyStatus:      nil,
			ProxyError:       &errString,
		}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		errString := err.Error()
		return XProtocolProxyCallResponse{
			Success:          false,
			Data:             "",
			Error:            &errString,
			ProxyStatus:      nil,
			ProxyError:       &errString,
			ProxyServerError: false,
		}
	}

	if res.StatusCode != http.StatusOK {
		errString := string(body)
		return XProtocolProxyCallResponse{
			Success:          false,
			Data:             "",
			Error:            &errString,
			ProxyStatus:      &res.StatusCode,
			ProxyError:       &errString,
			ProxyServerError: true,
		}
	}

	var response XProtocolCallResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		errString := err.Error()
		return XProtocolProxyCallResponse{
			Success:          false,
			Data:             "",
			Error:            &errString,
			ProxyStatus:      nil,
			ProxyError:       &errString,
			ProxyServerError: false,
		}
	}

	// if response.Data is a string, convert it to json.RawMessage
	// if reflect.TypeOf(response.Data).Kind() == reflect.String {
	// 	response.Data = json.RawMessage(response.Data.(string))
	// }

	err = json.Unmarshal([]byte(response.Data), &response.Data)
	if err != nil {
		errString := err.Error()
		return XProtocolProxyCallResponse{
			Success:          false,
			Data:             "",
			Error:            &errString,
			ProxyStatus:      nil,
			ProxyError:       &errString,
			ProxyServerError: false,
		}
	}

	return XProtocolProxyCallResponse{
		Success:          true,
		Data:             string(response.Data),
		Error:            response.Error,
		ProxyStatus:      nil,
		ProxyError:       nil,
		ProxyServerError: false,
	}
}
