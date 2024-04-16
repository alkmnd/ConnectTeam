package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
	"time"
)

type GameService struct {
	repo repository.Game
}

func NewGameService(repo repository.Game) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) SaveResults(gameId int, userId int, rate int) error {
	return s.repo.SaveResults(gameId, userId, rate)
}

func (s *GameService) CreateGame(creatorId int, startDateString string, name string) (game connectteam.Game, err error) {

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

func (s *GameService) CreateParticipant(userId int, gameId int) error {
	return s.repo.CreateParticipant(userId, gameId)
}

func (s *GameService) StartGame(gameId int) error {
	return s.repo.StartGame(gameId)
}

func (s *GameService) EndGame(gameId int) error {
	return s.repo.EndGame(gameId)
}

func (s *GameService) GetResults(gameId int) (results []connectteam.UserResult, err error) {
	return s.repo.GetResults(gameId)
}

func (s *GameService) GetCreatedGames(page int, userId int) ([]connectteam.Game, error) {
	return s.repo.GetCreatedGames(page, userId)
}

func (s *GameService) GetGame(gameId int) (connectteam.Game, error) {
	return s.repo.GetGame(gameId)
}

func (s *GameService) DeleteGame(gameId int) error {
	return s.repo.DeleteGame(gameId)
}

func (s *GameService) GetGameWithInvitationCode(code string) (connectteam.Game, error) {
	return s.repo.GetGameWithInvitationCode(code)
}

func (s *GameService) GetGames(page int, userId int) ([]connectteam.Game, error) {
	return s.repo.GetGames(page, userId)
}
