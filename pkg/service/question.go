package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
)

type QuestionService struct {
	repo repository.Question
}

func NewQuestionService(repo repository.Question) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) CreateQuestion(content string, topicId int) (int, error) {
	if len(content) == 0 {
		return 0, errors.New("empty string")
	}

	return s.repo.CreateQuestion(content, topicId)
}

func (s *QuestionService) DeleteQuestion(id int) error {
	return s.repo.DeleteQuestion(id)
}

func (s *QuestionService) GetAll(topicId int) ([]connectteam.Question, error) {
	return s.repo.GetAll(topicId)
}

func (s *QuestionService) UpdateQuestion(content string, id int) (connectteam.Question, error) {
	var q connectteam.Question
	if len(content) == 0 {
		return q, errors.New("empty string")
	}

	return s.repo.UpdateQuestion(content, id)
}
