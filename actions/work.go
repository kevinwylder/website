package actions

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Work is a requst to notify arrival or leaving work
type work struct {
	time    time.Time
	Arrived bool `json:"arrived"`
}

func (work *work) Read(data []byte, t time.Time) error {
	err := json.Unmarshal(data, work)
	if err != nil {
		return err
	}
	work.time = t
	return nil
}

func (work *work) Store(db *sql.DB) int {
	return 300
}
