package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type TopicService struct {
	repo repository.Topic
}

func NewTopicService(repo repository.Topic) *TopicService {
	return &TopicService{repo: repo}
}

func (s *TopicService) CreateTopic(topic connectteam.Topic) (uuid.UUID, error) {
	if strings.ReplaceAll(topic.Title, " ", "") == "" {
		return uuid.Nil, errors.New("incorrect title")
	}

	return s.repo.CreateTopic(topic)
}

func (s *TopicService) GetAll() ([]connectteam.Topic, error) {
	return s.repo.GetAll()
}

func (s *TopicService) DeleteTopic(id uuid.UUID) error {
	return s.repo.DeleteTopic(id)
}

func (s *TopicService) GetTopic(id uuid.UUID) (connectteam.Topic, error) {
	return s.repo.GetTopic(id)
}

func (s *TopicService) UpdateTopic(id uuid.UUID, title string) error {
	return s.repo.UpdateTopic(id, title)
}
