package repository

import (
	connectteam "ConnectTeam/models"
	// connectteam "ConnectTeam"
	"fmt"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type TopicPostgres struct {
	db *sqlx.DB
}

func NewTopicPostgres(db *sqlx.DB) *TopicPostgres {
	return &TopicPostgres{db: db}
}

func (r *TopicPostgres) CreateTopic(topic connectteam.Topic) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (title) values ($1) RETURNING id", topicsTable)
	row := r.db.QueryRow(query, topic.Title)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *TopicPostgres) GetAll() ([]connectteam.Topic, error) {
	var topics []connectteam.Topic

	query := fmt.Sprintf("SELECT id, title FROM %s", topicsTable)
	err := r.db.Select(&topics, query)
	return topics, err
}

func (r *TopicPostgres) GetTopic(topicId uuid.UUID) (topic connectteam.Topic, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", topicsTable)
	err = r.db.Get(&topic, query, topicId)
	return topic, err
}

func (r *TopicPostgres) DeleteTopic(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", topicsTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *TopicPostgres) UpdateTopic(id uuid.UUID, title string) error {

	query := fmt.Sprintf(`UPDATE %s SET title = $1 WHERE id = $2`, topicsTable)

	_, err := r.db.Exec(query, title, id)

	return err
}

func (r *TopicPostgres) GetRandWithLimit(limit int) (topics []connectteam.Topic, err error) {

	query := fmt.Sprintf("SELECT id, title FROM %s ORDER BY RANDOM() LIMIT $1", topicsTable)
	err = r.db.Select(&topics, query, limit)
	return topics, err
}
