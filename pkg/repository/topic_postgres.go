package repository

import (
	// connectteam "ConnectTeam"
	"fmt"

	connectteam "ConnectTeam"

	"github.com/jmoiron/sqlx"
)

type TopicPostgres struct {
	db *sqlx.DB
}

func NewTopicPostgres(db *sqlx.DB) *TopicPostgres {
	return &TopicPostgres{db: db}
}

func (r *TopicPostgres) CreateTopic(topic connectteam.Topic) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (title) values ($1) RETURNING id", topicsTable)
	row := r.db.QueryRow(query, topic.Title)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TopicPostgres) GetAll() ([]connectteam.Topic, error) {
	var topics []connectteam.Topic

	query := fmt.Sprintf("SELECT id, title FROM %s", topicsTable)
	err := r.db.Select(&topics, query)
	return topics, err
}

func (r *TopicPostgres) DeleteTopic(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", topicsTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *TopicPostgres) UpdateTopic(id int, title string) error {

	query := fmt.Sprintf(`UPDATE %s SET title = $1 WHERE id = %d`, topicsTable, id)

	_, err := r.db.Exec(query, title)

	return err
}
