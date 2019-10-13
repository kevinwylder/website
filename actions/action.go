package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Event is the interface an Action must satisfy to be logged
type Event interface {
	// Read loads a validated request and is given a convenience timestamp
	Read([]byte, time.Time) error
	// Store puts the event in the database
	Store(*sql.DB) int
}

// Action is a generic action which should be included in the message to dispatch
type Action struct {
	Timestamp int64  `json:"time"`
	Type      string `json:"action"`
}

// GetEvent will get the event of the action
func GetEvent(data []byte) (Event, error) {
	action := &Action{}
	err := json.Unmarshal(data, action)
	if err != nil {
		return nil, fmt.Errorf("Could not decode request: %s", err.Error())
	}
	t := time.Unix(action.Timestamp, 0)
	if t.Add(time.Minute).Before(time.Now()) {
		return nil, fmt.Errorf("Cannot accept request, message timestamp expired")
	}

	var event Event
	switch action.Type {
	case "climb":
		event = &climb{}
	case "work":
		event = &work{}
	default:
		return nil, fmt.Errorf("Unknown action: %s", action.Type)
	}
	return event, event.Read(data, t)
}
