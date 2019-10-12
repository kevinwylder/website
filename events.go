package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
)

// Event is a structure that holds signed data
type Event struct {
	Signature string          `json:"signature"`
	Data      json.RawMessage `json:"data"`
}

// Verify checks a signature based on the secret
func (event *Event) Verify(secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(event.Data)
	expectedMAC := mac.Sum(nil)

	recievedMAC, err := base64.StdEncoding.DecodeString(event.Signature)
	if err != nil {
		return false
	}

	return hmac.Equal(recievedMAC, expectedMAC)
}

// Store puts the event in the database
func (event *Event) Store(db *sql.DB) {
	// TODO
}
