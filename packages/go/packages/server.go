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
	Host      string                         `json:"host"`
	Port      int                            `json:"port"`
	ProxyMode bool                           `json:"proxy_mode"`
	Calls     map[string]XProtocolServerCall `json:"calls"`
}

type XProtocolServerCall struct {
	Name    string                                        `json:"name"`
	Handler func(payload json.RawMessage) json.RawMessage `json:"handler"`
}

type XProtocolCallRequest struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type XProtocolCallResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
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

			response := XProtocolCallResponse{
				Success: true,
				Data:    call.Handler(callRequest.Payload),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		} else {
			w.Write([]byte("unknown"))
		}
	})
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), mux)
}

func (s *XProtocolServer) RegisterCall(name string, handler func(payload json.RawMessage) json.RawMessage) {
	s.Calls[name] = XProtocolServerCall{
		Name:    name,
		Handler: handler,
	}
}

func (s *XProtocolServer) RegisterCallWithPayload(name string, handler interface{}) {
	s.Calls[name] = XProtocolServerCall{
		Name: name,
		Handler: func(payload json.RawMessage) json.RawMessage {
			p := reflect.New(reflect.TypeOf(handler).In(0)).Interface()
			json.Unmarshal(payload, &p)
			return handler.(func(p interface{}) json.RawMessage)(p)
		},
	}
}
