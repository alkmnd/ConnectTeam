package models

import (
	"encoding/json"
	"time"
)

type Notification struct {
	Type    string    `json:"type" redis:"type"`
	Payload string    `json:"payload" redis:"payload"`
	Date    time.Time `json:"date" redis:"date"`
}

func (n Notification) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

func (n Notification) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &n)
}