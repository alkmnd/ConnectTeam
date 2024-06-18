package models

import "github.com/google/uuid"

type Question struct {
	Id      uuid.UUID `json:"id"`
	Content string    `json:"content"`
	TopicId uuid.UUID `json:"topic_id"`
	Tags    []Tag     `json:"tags"`
}
