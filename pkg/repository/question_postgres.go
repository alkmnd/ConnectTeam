package repository

import (
	connectteam "ConnectTeam/models"
	"ConnectTeam/pkg/repository/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type QuestionPostgres struct {
	db *sqlx.DB
}

func (r *QuestionPostgres) GetQuestionTags(questionId uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag

	query := fmt.Sprintf("SELECT id, name FROM %s t JOIN %s tq ON t.id = tq.tag_id WHERE tq.question_id = $1", tagsTable, tagsQuestionsTable)
	err := r.db.Select(&tags, query, questionId)
	return tags, err
}

func NewQuestionPostgres(db *sqlx.DB) *QuestionPostgres {
	return &QuestionPostgres{db: db}
}

func (r *QuestionPostgres) CreateQuestion(content string, topicId uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (content, topic_id) values ($1, $2) RETURNING id", questionsTable)
	row := r.db.QueryRow(query, content, topicId)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *QuestionPostgres) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	query := fmt.Sprintf("SELECT id, name FROM %s", tagsTable)
	err := r.db.Select(&tags, query)
	return tags, err

}

func (r *QuestionPostgres) UpdateQuestionTags(questionId uuid.UUID, tags []models.Tag) ([]models.Tag, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return tags, err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
		err := tx.Commit()
		if err != nil {
			return
		}
	}()

	query := fmt.Sprintf("DELETE FROM %s WHERE question_id = $1", tagsQuestionsTable)
	_, err = tx.Exec(query, questionId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	for i := range tags {
		query := fmt.Sprintf("INSERT INTO %s (tag_id, question_id) VALUES ($1, $2) ON CONFLICT (tag_id, question_id) DO NOTHING", tagsQuestionsTable)
		_, err := tx.Exec(query, tags[i].Id, questionId)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	_ = tx.Commit()

	query = fmt.Sprintf("SELECT id, name FROM %s t JOIN %s tq ON t.id = tq.tag_id WHERE tq.question_id = $1", tagsTable, tagsQuestionsTable)
	err = r.db.Select(&tags, query, questionId)

	return tags, err

}

func (r *QuestionPostgres) DeleteQuestion(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", questionsTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *QuestionPostgres) GetAll(topicId uuid.UUID) ([]connectteam.Question, error) {
	var questions []connectteam.Question

	query := fmt.Sprintf("SELECT id, topic_id, content FROM %s WHERE topic_id = $1", questionsTable)
	err := r.db.Select(&questions, query, topicId)
	return questions, err
}

func (r *QuestionPostgres) GetRandWithLimit(topicId uuid.UUID, limit int) ([]connectteam.Question, error) {
	var questions []connectteam.Question

	query := fmt.Sprintf("SELECT id, topic_id, content FROM %s WHERE topic_id = $1 ORDER BY RANDOM() LIMIT $2", questionsTable)
	err := r.db.Select(&questions, query, topicId, limit)
	return questions, err
}

func (r *QuestionPostgres) GetTagsUsers(userId uuid.UUID, gameId uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	query := fmt.Sprintf("SELECT t.id, t.name FROM %s t JOIN %s tu ON t.id = tu.tag_id WHERE game_id = $1 AND user_id = $2", tagsTable, tagsUsersTable)
	err := r.db.Select(&tags, query, gameId, userId)
	return tags, err
}

func (r *QuestionPostgres) CreateTagsUsers(userId uuid.UUID, gameId uuid.UUID, tagId uuid.UUID) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, game_id, tag_id) values ($1, $2, $3)  ON CONFLICT (user_id, tag_id, game_id) DO NOTHING", tagsUsersTable)
	_, err := r.db.Exec(query, userId, gameId, tagId)
	return err

}

func (r *QuestionPostgres) UpdateQuestion(content string, id uuid.UUID) (connectteam.Question, error) {
	var question connectteam.Question
	query := fmt.Sprintf(`UPDATE %s SET content = $1 WHERE id = $2 RETURNING id, topic_id, content`, questionsTable)
	row := r.db.QueryRow(query, content, id)
	err := row.Scan(&question.Id, &question.TopicId, &question.Content)

	return question, err
}
