package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
	"github.com/google/uuid"
	"time"
)

type GameService struct {
	repo repository.Game
}

func NewGameService(repo repository.Game) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) SaveResults(gameId uuid.UUID, userId uuid.UUID, rate int) error {
	return s.repo.SaveResults(gameId, userId, rate)
}

func (s *GameService) CreateGame(creatorId uuid.UUID, startDateString string, name string) (game connectteam.Game, err error) {

	game.InvitationCode, err = generateInviteCode()
	if err != nil {
		return game, err
	}

	date, err := time.Parse(time.RFC3339, startDateString)
	if err != nil {
		return game, err
	}
	startDate := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, date.Location())

	if !startDate.After(time.Now()) {
		return game, errors.New("incorrect start date")
	}
	game.StartDate = startDate
	game.CreatorId = creatorId
	game.Status = "not_started"
	if len(name) == 0 {
		return game, errors.New("incorrect game name")
	}
	game.Name = name

	return s.repo.CreateGame(game)
}

func (s *GameService) CreateParticipant(userId uuid.UUID, gameId uuid.UUID) error {
	return s.repo.CreateParticipant(userId, gameId)
}

func (s *GameService) StartGame(gameId uuid.UUID) error {
	return s.repo.StartGame(gameId)
}

func (s *GameService) EndGame(gameId uuid.UUID) error {
	return s.repo.EndGame(gameId)
}

func (s *GameService) GetResults(gameId uuid.UUID) (results []connectteam.UserResult, err error) {
	return s.repo.GetResults(gameId)
}

func (s *GameService) GetCreatedGames(page int, userId uuid.UUID) ([]connectteam.Game, error) {
	return s.repo.GetCreatedGames(page, userId)
}

func (s *GameService) GetGame(gameId uuid.UUID) (connectteam.Game, error) {
	return s.repo.GetGame(gameId)
}

func (s *GameService) DeleteGame(gameId uuid.UUID) error {
	return s.repo.DeleteGame(gameId)
}

func (s *GameService) GetGameWithInvitationCode(code string) (connectteam.Game, error) {
	return s.repo.GetGameWithInvitationCode(code)
}

func (s *GameService) GetGames(page int, userId uuid.UUID) ([]connectteam.Game, error) {
	return s.repo.GetGames(page, userId)
}

func (s *GameService) CancelGame(gameId uuid.UUID, userId uuid.UUID) error {
	game, err := s.repo.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.CreatorId != userId {
		return errors.New("permission denied")
	}
	return s.repo.CancelGame(gameId)
}
