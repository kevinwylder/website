package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Climb has a grade and a type
type climb struct {
	Grade string `json:"grade"`
	Type  string `json:"type"`
}

func (climb *climb) Read(data []byte, timestamp time.Time) error {
	err := json.Unmarshal(data, climb)
	if err != nil {
		return fmt.Errorf("(climb) cannot read request: %s", err.Error())
	}
	return nil
}

func (climb *climb) Store(db *sql.DB) int {
	_, err := db.Exec(`
INSERT INTO climbs (type, grade) VALUES (?, ?);
`, climb.Type, climb.Grade)
	if err != nil {
		log.Println("(climb) SQL error " + err.Error())
		return 500
	}
	return 200
}
