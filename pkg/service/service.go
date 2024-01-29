package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GenerateToken(login, password string, isEmail bool) (string, error)
	ParseToken(token string) (int, string, error)
	VerifyPhone(verifyPhone connectteam.VerifyPhone) (string, error)
	VerifyUser(verifyUser connectteam.VerifyUser)  error
	VerifyEmail(verifyEmail connectteam.VerifyEmail) (int, error)
	DeleteVerificationCode(id int, code string) (error)
	
}

type UserInterface interface {
	GetUserById(id int) (connectteam.UserPublic, error)
	ChangeAccessWithId(id int, access string) (error)
	GetUsersList() ([] connectteam.UserPublic, error)
	ChangePassword(old_password string, new_password string, id int) (error)
	ChangeEmail(id int, newEmail string, code string) (error) 
	DeleteVerificationCode(id int, code string) (error)
	CheckEmailOnChange(id int, email string, password string) (error)
	ChangePersonalData(id int, user connectteam.UserPersonalInfo) (error)
}


type Service struct {
	Authorization
	UserInterface
}
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		UserInterface: NewUserService(repos.UserInterface),
	}
}