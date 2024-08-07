package models

import "github.com/google/uuid"

type UserResult struct {
	Id              int       `json:"id" db:"id"`
	UserId          uuid.UUID `json:"user_id" db:"user_id"`
	UserTemporaryId uuid.UUID `json:"user_temp_id" db:"user_temp_id"`
	Value           int       `json:"value" db:"value"`
	Name            string    `json:"name" db:"name"`
}
