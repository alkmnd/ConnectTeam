package repository

import (
	connectteam "ConnectTeam"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type QuestionPostgres struct {
	db *sqlx.DB
}

func NewQuestionPostgres(db *sqlx.DB) *QuestionPostgres {
	return &QuestionPostgres{db: db}
}

func (r *QuestionPostgres) CreateQuestion(content string, topicId int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (content, topic_id) values ($1, $2) RETURNING id", questionsTable)
	row := r.db.QueryRow(query, content, topicId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *QuestionPostgres) DeleteQuestion(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", questionsTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *QuestionPostgres) GetAll(topicId int) ([]connectteam.Question, error) {
	var questions []connectteam.Question

	query := fmt.Sprintf("SELECT id, topic_id, content FROM %s WHERE topic_id = $1", questionsTable)
	err := r.db.Select(&questions, query, topicId)
	return questions, err
}

func (r *QuestionPostgres) UpdateQuestion(content string, id int) (connectteam.Question, error) {
	var question connectteam.Question
	query := fmt.Sprintf(`UPDATE %s SET content = $1 WHERE id = %d RETURNING id, topic_id, content`, questionsTable, id)
	row := r.db.QueryRow(query, content)
	err := row.Scan(&question.Id, &question.TopicId, &question.Content)

	return question, err
}
