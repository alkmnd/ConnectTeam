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
		return 0, errors.New("incorrect title")
	}

	return s.repo.CreateTopic(topic)
}

func (s *TopicService) GetAll() ([] connectteam.Topic, error) {
	return s.repo.GetAll()
}

func (s *TopicService) DeleteTopic(id int) (error) {
	return s.repo.DeleteTopic(id)
}

func (s *TopicService) UpdateTopic(id int, title string) (error) {
	return s.repo.UpdateTopic(id, title)
}