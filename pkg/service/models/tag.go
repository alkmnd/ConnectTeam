package models

import "github.com/google/uuid"

type Tag struct {
	Id   uuid.UUID
	Name string
}
