package models

import "github.com/google/uuid"

type Tag struct {
	Id   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}
