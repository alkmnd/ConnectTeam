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