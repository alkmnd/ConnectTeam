package service

import (
	connectteam "ConnectTeam/models"
	"ConnectTeam/pkg/repository"
	"errors"
	"github.com/google/uuid"
	"time"
)

type GameService struct {
	gameRepo         repository.Game
	notificationRepo repository.Notification
	planRepo         repository.Plan
}

func NewGameService(gameRepo repository.Game, notificationRepo repository.Notification, planRepo repository.Plan) *GameService {
	return &GameService{gameRepo: gameRepo, notificationRepo: notificationRepo, planRepo: planRepo}
}

func (s *GameService) SaveResults(gameId uuid.UUID, userId uuid.UUID, rate int) error {
	return s.gameRepo.SaveResults(gameId, userId, rate)
}
func (s *GameService) ChangeStartDate(gameId uuid.UUID, dateString string) error {
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return err
	}
	startDate := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, date.Location())
	if !startDate.After(time.Now()) {
		return errors.New("incorrect start date")
	}

	return s.gameRepo.ChangeStartDate(gameId, startDate)
}

func (s *GameService) ChangeGameName(gameId uuid.UUID, name string) error {
	if len(name) == 0 {
		return errors.New("incorrect game name")
	}

	return s.gameRepo.ChangeGameName(gameId, name)
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

	return s.gameRepo.CreateGame(game)
}

func (s *GameService) CreateParticipant(userId uuid.UUID, gameId uuid.UUID) error {
	return s.gameRepo.CreateParticipant(userId, gameId)
}

func (s *GameService) StartGame(gameId uuid.UUID) error {
	game, err := s.gameRepo.GetGame(gameId)
	if err != nil {
		return err
	}

	err = s.gameRepo.StartGame(game.Id)
	if err != nil {
		return err
	}

	//users, err := s.gameRepo.GetGameParticipants(gameId)
	//if err != nil {
	//	return nil
	//}

	//for i := range users {
	//	if users[i].Id != game.CreatorId {
	//		_ = s.notificationRepo.CreateGameStartNotification(gameId, users[i].Id)
	//	}
	//}
	return nil

}

func (s *GameService) GetGameParticipants(gameId uuid.UUID) (users []connectteam.UserPublic, err error) {
	return s.gameRepo.GetGameParticipants(gameId)
}

func (s *GameService) EndGame(gameId uuid.UUID) error {
	return s.gameRepo.EndGame(gameId)
}

func (s *GameService) GetResults(gameId uuid.UUID) (results []connectteam.UserResult, err error) {
	return s.gameRepo.GetResults(gameId)
}

func (s *GameService) GetCreatedGames(page int, userId uuid.UUID) ([]connectteam.Game, error) {
	return s.gameRepo.GetCreatedGames(page, userId)
}

func (s *GameService) GetGame(gameId uuid.UUID) (connectteam.Game, error) {
	return s.gameRepo.GetGame(gameId)
}

func (s *GameService) DeleteGameFromGameList(gameId uuid.UUID, userId uuid.UUID) error {
	return s.gameRepo.DeleteGameFromGameList(gameId, userId)
}

func (s *GameService) GetGameWithInvitationCode(code string) (connectteam.Game, error) {
	return s.gameRepo.GetGameWithInvitationCode(code)
}

func (s *GameService) GetGames(page int, userId uuid.UUID) ([]connectteam.Game, error) {
	return s.gameRepo.GetGames(page, userId)
}

func (s *GameService) CancelGame(gameId uuid.UUID, userId uuid.UUID) error {
	game, err := s.gameRepo.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.CreatorId != userId {
		return errors.New("permission denied")
	}

	if game.Status != "not_started" {
		return errors.New("permission denied")
	}
	err = s.gameRepo.CancelGame(gameId)
	if err != nil {
		return err
	}

	users, err := s.gameRepo.GetGameParticipants(gameId)
	if err != nil {
		return nil
	}

	for i := range users {
		if users[i].Id != game.CreatorId {
			err = s.notificationRepo.CreateGameCancelNotification(gameId, users[i].Id)
		}
	}
	return nil
}

func (s *GameService) InviteUserToGame(gameId uuid.UUID, userId uuid.UUID, creatorId uuid.UUID) error {
	game, err := s.gameRepo.GetGame(gameId)
	if err != nil {
		return err
	}
	if game.CreatorId == userId {
		return errors.New("permission denied")
	}
	if game.CreatorId != creatorId {
		return errors.New("permission denied")
	}
	return s.notificationRepo.CreateGameInviteNotification(gameId, userId)
}
