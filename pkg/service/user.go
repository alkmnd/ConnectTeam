package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
)

type UserService struct {
	repo repository.UserInterface
}
func NewUserService(repo repository.UserInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(id int) (connectteam.User, error) {
	user, err := s.repo.GetUserById(id)
	return user, err 
}

func (s *UserService) ChangeAccessById(id int, access string) (error) {
	if err := s.repo.ChangeAccessById(id, access); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUsersList() ([]connectteam.UserPublic, error) {
	return s.repo.GetUsersList()
}