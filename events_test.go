package main

import (
	"encoding/json"
	"testing"
)

func TestSignatureVerify(test *testing.T) {
	secret := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	passing := [][]byte{
		[]byte(`{"signature":"Xlu3ZnlOXY8Shti1YP0HyYTAmcRiYliC+DL8+UoYLFc=","data":{"action":"test"}}`),
	}
	failing := [][]byte{
		[]byte(`{"signature":"Xlu3ZnlOXY8Shti1YP0HyYTAmcRiYliC+DL8+UoYLFc","data":{"action":"test"}}`),
		[]byte(`{"signature":"","data":{"action":"test"}}`),
	}

	for _, message := range passing {
		event := &Event{}
		err := json.Unmarshal(message, event)
		if err != nil {
			test.Fatal(err)
		}
		if !event.Verify(secret) {
			test.Fatal("Should have matched the HMAC sum")
		}
	}

	for _, message := range failing {
		event := &Event{}
		err := json.Unmarshal(message, event)
		if err != nil {
			test.Fatal(err)
		}
		if event.Verify(secret) {
			test.Fatal("Should not have matched the HMAC sum")
		}
	}
}
