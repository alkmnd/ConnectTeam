package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GenerateToken(login, password string, isEmail bool) (string, string, error)
	ParseToken(token string) (int, string, error)
	VerifyPhone(verifyPhone connectteam.VerifyPhone) (string, error)
	VerifyUser(verifyUser connectteam.VerifyUser)  error
	VerifyEmail(verifyEmail connectteam.VerifyEmail) (int, error)
	DeleteVerificationCode(id int, code string) (error)
	
}

type User interface {
	GetUserById(id int) (connectteam.UserPublic, error)
	UpdateAccessWithId(id int, access string) (error)
	GetUsersList() ([] connectteam.UserPublic, error)
	UpdatePassword(old_password string, new_password string, id int) (error)
	UpdateEmail(id int, newEmail string, code string) (error) 
	DeleteVerificationCode(id int, code string) (error)
	CheckEmailOnChange(id int, email string, password string) (error)
	UpdatePersonalData(id int, user connectteam.UserPersonalInfo) (error)
	UpdateCompanyData(id int, company connectteam.UserCompanyData) (error)
	GetUserPlan(user_id int) (connectteam.UserPlan, error)
}

type Plan interface {
	GetUserPlan(user_id int) (connectteam.UserPlan, error)
	CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error)
}

type Service struct {
	Authorization
	User
	Plan
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User: NewUserService(repos.User),
		Plan: NewPlanService(repos.Plan),
	}
}