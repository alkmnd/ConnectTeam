package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	repoModels "ConnectTeam/pkg/repository/models"
	"ConnectTeam/pkg/service/models"
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

func (s *QuestionService) GetRandWithLimit(topicId int, limit int) ([]models.Question, error) {
	var questions []models.Question
	repoQuestions, err := s.repo.GetRandWithLimit(topicId, limit)
	if err != nil {
		return nil, err
	}
	for i := range repoQuestions {
		var tags []models.Tag
		questionTags, err := s.repo.GetQuestionTags(repoQuestions[i].Id)
		for j := range questionTags {
			tags = append(tags, models.Tag{
				Id:   questionTags[j].Id,
				Name: questionTags[j].Name,
			})
		}
		if err != nil {
			questions = append(questions, models.Question{
				Id:      repoQuestions[i].Id,
				TopicId: repoQuestions[i].TopicId,
				Content: repoQuestions[i].Content,
			})
			continue
		}
		questions = append(questions, models.Question{
			Id:      repoQuestions[i].Id,
			TopicId: repoQuestions[i].TopicId,
			Content: repoQuestions[i].Content,
			Tags:    tags,
		})
	}
	return questions, nil
}

func (s *QuestionService) DeleteQuestion(id int) error {
	return s.repo.DeleteQuestion(id)
}

func (s *QuestionService) GetAll(topicId int) ([]models.Question, error) {
	var questions []models.Question
	repoQuestions, err := s.repo.GetAll(topicId)
	if err != nil {
		return nil, err
	}
	for i := range repoQuestions {
		var tags []models.Tag
		questionTags, err := s.repo.GetQuestionTags(repoQuestions[i].Id)
		for j := range questionTags {
			tags = append(tags, models.Tag{
				Id:   questionTags[j].Id,
				Name: questionTags[j].Name,
			})
		}
		if err != nil {
			questions = append(questions, models.Question{
				Id:      repoQuestions[i].Id,
				TopicId: repoQuestions[i].TopicId,
				Content: repoQuestions[i].Content,
			})
			continue
		}
		questions = append(questions, models.Question{
			Id:      repoQuestions[i].Id,
			TopicId: repoQuestions[i].TopicId,
			Content: repoQuestions[i].Content,
			Tags:    tags,
		})
	}
	return questions, nil
}

func (s *QuestionService) UpdateQuestion(content string, id int) (connectteam.Question, error) {
	var q connectteam.Question
	if len(content) == 0 {
		return q, errors.New("empty string")
	}

	return s.repo.UpdateQuestion(content, id)
}

func (s *QuestionService) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	repoTags, err := s.repo.GetAllTags()
	if err != nil {
		return nil, err
	}
	for i := range repoTags {
		tags = append(tags, models.Tag{
			Id:   repoTags[i].Id,
			Name: repoTags[i].Name,
		})
	}

	return tags, nil
}

func (s *QuestionService) UpdateQuestionTags(questionId int, tags []models.Tag) ([]models.Tag, error) {
	var repoTags []repoModels.Tag
	for i := range tags {
		repoTags = append(repoTags, repoModels.Tag{
			Id: tags[i].Id,
		})
	}
	repoTags, err := s.repo.UpdateQuestionTags(questionId, repoTags)
	if err != nil {
		return nil, err
	}
	tags = make([]models.Tag, 0)
	for i := range repoTags {
		tags = append(tags, models.Tag{
			Id:   repoTags[i].Id,
			Name: repoTags[i].Name,
		})
	}

	return tags, nil
}
