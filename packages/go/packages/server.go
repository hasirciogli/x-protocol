package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type XProtocolServer struct {
	Host          string                           `json:"host"`
	Port          int                              `json:"port"`
	ProxyMode     bool                             `json:"proxy_mode"`
	Calls         map[string]XProtocolServerCall   `json:"calls"`
	ProxyChannels map[string]XProtocolProxyChannel `json:"proxy_channels"`
}

type XProtocolServerCall struct {
	Name    string                                              `json:"name"`
	Handler func(payload json.RawMessage) XProtocolCallResponse `json:"handler"`
}

type XProtocolCallRequest struct {
	ProxyChannelName *string                `json:"proxy_channel_name"`
	Name             string                 `json:"name"`
	Payload          json.RawMessage        `json:"payload"`
	FromProxyChannel *XProtocolProxyChannel `json:"from_proxy_channel"`
}

type XProtocolCallResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *string         `json:"error"`
}

func NewXProtocolServer(host string, port int) *XProtocolServer {
	return &XProtocolServer{
		Host:          host,
		Port:          port,
		ProxyMode:     false,
		Calls:         map[string]XProtocolServerCall{},
		ProxyChannels: map[string]XProtocolProxyChannel{},
	}
}

func (s *XProtocolServer) Start() {
	fmt.Println("Starting XProtocolServer on http://" + s.Host + ":" + strconv.Itoa(s.Port))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/calls" && r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			bodyAsString := string(body)
			var callRequest XProtocolCallRequest
			err = json.Unmarshal([]byte(bodyAsString), &callRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if callRequest.ProxyChannelName != nil && *callRequest.ProxyChannelName != "" {
				proxyChannel, ok := s.ProxyChannels[*callRequest.ProxyChannelName]
				if !ok {
					http.Error(w, "Proxy channel not found", http.StatusNotFound)
					return
				}

				response := proxyChannel.Call(callRequest.Name, callRequest)
				if response.Error != nil {
					if response.ProxyServerError {
						w.Header().Set("Content-Type", "plain/text")
						w.WriteHeader(*response.ProxyStatus)
						w.Write([]byte(*response.ProxyError))
					} else {
						w.Header().Set("Content-Type", "plain/text")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(*response.Error))
					}
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(XProtocolCallResponse{
					Success: true,
					Data:    json.RawMessage(response.Data),
					Error:   nil,
				})
				return
			}

			call, ok := s.Calls[callRequest.Name]
			if !ok {
				http.Error(w, "Call not found", http.StatusNotFound)
				return
			}
			fmt.Println("Call isteği alındı -> "+callRequest.Name+" | from proxy:", callRequest.FromProxyChannel != nil)

			appMode := os.Getenv("APP_MODE")
			if appMode == "development" {
				fmt.Println("Call isteği alındı -> "+callRequest.Name+" | from proxy:", callRequest.FromProxyChannel != nil)
			}

			response := call.Handler(callRequest.Payload)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(XProtocolCallResponse{
				Success: true,
				Data:    response.Data,
				Error:   nil,
			})
			return
		} else {
			w.Write([]byte("unknown"))
		}
	})
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), mux)
}

func (s *XProtocolServer) RegisterCall(name string, handler func(payload json.RawMessage) XProtocolCallResponse) {
	s.Calls[name] = XProtocolServerCall{
		Name:    name,
		Handler: handler,
	}

	fmt.Println("Call registered -> " + name + " ✔️")
}

// proxy channel start

func (s *XProtocolServer) RegisterProxyChannel(name string, host string, port int) {
	s.ProxyChannels[name] = XProtocolProxyChannel{
		Name: name,
		Host: host,
		Port: port,
	}

	fmt.Println("proxy channel registered -> " + name + " | " + host + ":" + strconv.Itoa(port) + " ✔️")
}

func (s *XProtocolServer) GetProxyChannel(name string) XProtocolProxyChannel {
	return s.ProxyChannels[name]
}

func (s *XProtocolServer) GetProxyChannelHost(name string) string {
	return s.ProxyChannels[name].Host
}

func (s *XProtocolServer) GetProxyChannelPort(name string) int {
	return s.ProxyChannels[name].Port
}

func (s *XProtocolServer) UpdateProxyChannel(name string, host string, port int) {
	s.ProxyChannels[name] = XProtocolProxyChannel{
		Host: host,
		Port: port,
	}
}
