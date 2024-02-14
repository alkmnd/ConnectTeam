package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
	"strings"
)

type TopicService struct {
	repo repository.Topic
} 

func NewTopicService(repo repository.Topic) *TopicService {
	return &TopicService{repo: repo}
}

func (s *TopicService) CreateTopic(topic connectteam.Topic) (int, error) {
	if strings.ReplaceAll(topic.Title, " ", "") == "" {
		return 0, errors.New("Incorrect title")
	}

	return s.repo.CreateTopic(topic)
}