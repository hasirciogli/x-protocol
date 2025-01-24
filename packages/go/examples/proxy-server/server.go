package main

import (
	"encoding/json"

	"github.com/hasirciogli/x-protocol/packages/go/packages"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func Login(payload json.RawMessage) packages.XProtocolCallResponse {
	var lp LoginPayload
	err := json.Unmarshal(payload, &lp)
	if err != nil {
		errString := err.Error()
		return packages.XProtocolCallResponse{
			Success: false,
			Data:    json.RawMessage("{}"),
			Error:   &errString,
		}
	}

	if lp.Username == "admin" && lp.Password == "123456" {
		lr := LoginResponse{
			Message: "Login successful",
			Status:  true,
		}

		jsonBytes, err := json.Marshal(lr)
		if err != nil {
			errString := err.Error()
			return packages.XProtocolCallResponse{
				Success: false,
				Data:    json.RawMessage(""),
				Error:   &errString,
			}
		}

		return packages.XProtocolCallResponse{
			Success: true,
			Data:    json.RawMessage(jsonBytes),
			Error:   nil,
		}
	}

	lr := LoginResponse{
		Message: "Login failed",
		Status:  false,
	}

	jsonBytes, err := json.Marshal(lr)
	if err != nil {
		errString := err.Error()
		return packages.XProtocolCallResponse{
			Success: false,
			Data:    json.RawMessage(""),
			Error:   &errString,
		}
	}

	return packages.XProtocolCallResponse{
		Success: true,
		Data:    json.RawMessage(jsonBytes),
		Error:   nil,
	}
}

func main() {
	server := packages.NewXProtocolServer("localhost", 8090)
	server.RegisterCall("login", Login)
	server.Start()
}
