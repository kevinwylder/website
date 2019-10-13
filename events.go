package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kevinwylder/website/actions"
)

// EventRequest is a structure that holds signed data
type EventRequest struct {
	Signature string          `json:"signature"`
	Data      json.RawMessage `json:"data"`
}

// Verify checks a signature based on the secret
func (req *EventRequest) Verify(secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(req.Data)
	expectedMAC := mac.Sum(nil)

	recievedMAC, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		return false
	}

	return hmac.Equal(recievedMAC, expectedMAC)
}

func (server *serverState) OnEvent(body io.Reader) int {
	request := &EventRequest{}
	decoder := json.NewDecoder(body)
	decoder.Decode(request)
	if !request.Verify(server.secret) {
		return 401
	}
	event, err := actions.GetEvent(request.Data)
	if err != nil {
		fmt.Println("Couldn't parse event: " + err.Error())
		return 400
	}
	return event.Store(server.db)
}
