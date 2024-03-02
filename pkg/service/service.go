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
	VerifyUser(verifyUser connectteam.VerifyUser) error
	VerifyEmail(verifyEmail connectteam.VerifyEmail) (int, error)
	DeleteVerificationCode(id int, code string) error
}

type User interface {
	GetUserById(id int) (connectteam.UserPublic, error)
	UpdateAccessWithId(id int, access connectteam.AccessLevel) error
	GetUsersList() ([]connectteam.UserPublic, error)
	UpdatePassword(oldPassword string, newPassword string, id int) error
	UpdateEmail(id int, newEmail string, code string) error
	DeleteVerificationCode(id int, code string) error
	CheckEmailOnChange(id int, email string, password string) error
	UpdatePersonalData(id int, user connectteam.UserPersonalInfo) error
	UpdateCompanyData(id int, company connectteam.UserCompanyData) error
	GetUserPlan(userId int) (connectteam.UserPlan, error)
	RestorePasswordAuthorized(id int) error
	RestorePassword(email string) error
}

type Plan interface {
	GetUserPlan(userId int) (connectteam.UserPlan, error)
	CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error)
	GetUsersPlans() ([]connectteam.UserPlan, error)
	ConfirmPlan(id int) error
	SetPlanByAdmin(userId int, planType string, expiryDateString string) error
	DeletePlan(id int) error
}

type Topic interface {
	CreateTopic(topic connectteam.Topic) (int, error)
	GetAll() ([]connectteam.Topic, error)
	DeleteTopic(id int) error
	UpdateTopic(id int, title string) error
}

type Question interface {
	CreateQuestion(content string, topicId int) (int, error)
	DeleteQuestion(id int) error
	GetAll(topicId int) ([]connectteam.Question, error)
	UpdateQuestion(content string, id int) (connectteam.Question, error)
}

type Service struct {
	Authorization
	User
	Plan
	Topic
	Question
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
		Plan:          NewPlanService(repos.Plan),
		Topic:         NewTopicService(repos.Topic),
		Question:      NewQuestionService(repos.Question),
	}
}
