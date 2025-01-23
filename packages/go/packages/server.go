package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

// Generic type T
type XProtocolServer struct {
	Host          string                           `json:"host"`
	Port          int                              `json:"port"`
	ProxyMode     bool                             `json:"proxy_mode"`
	Calls         map[string]XProtocolServerCall   `json:"calls"`
	ProxyChannels map[string]XProtocolProxyChannel `json:"proxy_channels"`
}

type XProtocolServerCall struct {
	Name             string                                              `json:"name"`
	Handler          func(payload json.RawMessage) XProtocolCallResponse `json:"handler"`
	FromProxyChannel *XProtocolProxyChannel                              `json:"from_proxy_channel"`
}

type XProtocolCallRequest struct {
	Name             string                 `json:"name"`
	Payload          json.RawMessage        `json:"payload"`
	FromProxyChannel *XProtocolProxyChannel `json:"from_proxy_channel"`
}

type XProtocolCallResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

func NewXProtocolServer(host string, port int) *XProtocolServer {
	return &XProtocolServer{
		Host:      host,
		Port:      port,
		ProxyMode: false,
		Calls:     map[string]XProtocolServerCall{},
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
			call, ok := s.Calls[callRequest.Name]
			if !ok {
				http.Error(w, "Call not found", http.StatusNotFound)
				return
			}

			response := call.Handler(callRequest.Payload)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
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

	fmt.Println("call registered -> " + name + " ✔️")
}

func (s *XProtocolServer) RegisterCallWithPayload(name string, handler interface{}) {
	s.Calls[name] = XProtocolServerCall{
		Name: name,
		Handler: func(payload json.RawMessage) XProtocolCallResponse {
			p := reflect.New(reflect.TypeOf(handler).In(0)).Interface()
			json.Unmarshal(payload, &p)
			return handler.(func(p interface{}) XProtocolCallResponse)(p)
		},
	}

	fmt.Println("call registered -> " + name + " ✔️")
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
