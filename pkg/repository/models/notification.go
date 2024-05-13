package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	Id      uuid.UUID `json:"id" redis:"id"`
	Type    string    `json:"type" redis:"type"`
	Payload string    `json:"payload" redis:"payload"`
	Date    time.Time `json:"date" redis:"date"`
	IsRead  bool      `json:"is_read" redis:"is_read"`
}

func (n Notification) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

func (n Notification) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &n)
}
