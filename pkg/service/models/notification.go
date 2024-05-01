package models

import "time"

type Notification struct {
	Type    string    `json:"type"`
	Payload string    `json:"payload"`
	Date    time.Time `json:"date"`
}
