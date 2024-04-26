package models

type Question struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	TopicId int    `json:"topic_id"`
	Tags    []Tag  `json:"tags"`
}
