package repository

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GetUserWithEmail(email, password string) (connectteam.User, error)
	GetIdWithEmail(email string) (int, error)
	GetUserWithPhone(phoneNumber, password string) (connectteam.User, error)
	//VerifyUser(verifyUser connectteam.VerifyUser) error
	GetVerificationCode(email string) (string, error)
	CreateVerificationCode(email string, code string) error
	CheckIfExist(id int) (bool, error)
	DeleteVerificationCode(email string, code string) error
}

type User interface {
	GetUserById(id int) (connectteam.UserPublic, error)
	UpdateAccessWithId(id int, access string) error
	GetUsersList() ([]connectteam.UserPublic, error)
	GetPassword(id int) (string, error)
	UpdatePasswordWithId(newPassword string, id int) error
	GetVerificationCode(email string) (string, error)
	GetEmailWithId(id int) (string, error)
	UpdateEmail(email string, id int) error
	CreateVerificationCode(email string, code string) error
	DeleteVerificationCode(email string, code string) error
	CheckIfExist(email string) (bool, error)
	UpdateUserFirstName(id int, firstName string) error
	UpdateUserSecondName(id int, secondName string) error
	UpdateUserDescription(id int, secondName string) error
	UpdateCompanyName(id int, companyName string) error
	UpdateCompanyInfo(id int, info string) error
	UpdateCompanyURL(id int, companyURL string) error
	GetUserPlan(userId int) (connectteam.UserPlan, error)
	CreatePlanRequest(request connectteam.PlanRequest) (int, error)
	GetUserCredentials(id int) (connectteam.UserCredentials, error)
	UpdatePasswordWithEmail(newPassword string, email string) error
}

type Plan interface {
	GetUserActivePlan(userId int) (connectteam.UserPlan, error)
	CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error)
	GetUsersPlans() ([]connectteam.UserPlan, error)
	SetConfirmed(id int) error
	DeletePlan(id int) error
	SetExpiredStatus(id int) error
	DeleteOnConfirmPlan(userId int) error
	GetUserSubscriptions(userId int) ([]connectteam.UserPlan, error)
	GetHolderWithInvitationCode(code string) (id int, err error)
	SetExpiredStatusWithUserId(userId int) error
	GetMembers(code string) (users []connectteam.UserPublic, err error)
	DeleteUserFromSub(id int) error
}

type Topic interface {
	CreateTopic(topic connectteam.Topic) (int, error)
	GetAll() ([]connectteam.Topic, error)
	DeleteTopic(id int) error
	UpdateTopic(id int, title string) error
	GetTopic(topicId int) (topic connectteam.Topic, err error)
}

type Question interface {
	CreateQuestion(content string, topicId int) (int, error)
	DeleteQuestion(id int) error
	GetAll(topicId int) ([]connectteam.Question, error)
	UpdateQuestion(content string, id int) (connectteam.Question, error)
	GetRandWithLimit(topicId int, limit int) ([]connectteam.Question, error)
	GetQuestionTags(questionId int) ([]models.Tag, error)
	UpdateQuestionTags(questionId int, tags []models.Tag) ([]models.Tag, error)
	GetAllTags() ([]models.Tag, error)
}

type Game interface {
	CreateGame(game connectteam.Game) (connectteam.Game, error)
	GetCreatedGames(page int, userId int) (games []connectteam.Game, err error)
	CreateParticipant(userId int, gameId int) error
	GetGame(gameId int) (game connectteam.Game, err error)
	DeleteGame(gameId int) error
	GetGameWithInvitationCode(code string) (game connectteam.Game, err error)
	GetGames(page int, userId int) (games []connectteam.Game, err error)
	StartGame(gameId int) error
	SaveResults(gameId int, userId int, rate int) error
	GetResults(gameId int) (results []connectteam.UserResult, err error)
	EndGame(gameId int) error
}

type Repository struct {
	Authorization
	User
	Plan
	Topic
	Question
	Game
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
		Plan:          NewPlanPostgres(db),
		Topic:         NewTopicPostgres(db),
		Question:      NewQuestionPostgres(db),
		Game:          NewGamePostgres(db),
	}
}
