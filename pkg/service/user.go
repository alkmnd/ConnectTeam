package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
)

type UserService struct {
	repo repository.UserInterface

}
func NewUserService(repo repository.UserInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(id int) (connectteam.UserPublic, error) {
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

func (s *UserService) ChangePassword(old_password string, new_password string, id int) (error) {
	db_password, err := s.repo.GetPassword(id)
	if err != nil {
		return err
	}

	if db_password != generatePasswordHash(old_password) {
		return errors.New("Wrong old password")
	}

	return s.repo.ChangePassword(generatePasswordHash(new_password), id)
}