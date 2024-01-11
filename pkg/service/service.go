package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GenerateToken(login, password string, isEmail bool) (string, error)
	ParseToken(token string) (int, error)
	VerifyPhone(verifyPhone connectteam.VerifyPhone) (string, error)
	VerifyUser(verifyUser connectteam.VerifyUser)  error
}


type Service struct {
	Authorization
}
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}