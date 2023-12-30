package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GenerateToken(email, password string) (string, error)
}

type Service struct {
	Authorization
}
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}