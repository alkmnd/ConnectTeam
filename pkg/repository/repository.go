package repository

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	"time"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user connectteam.User) (uuid.UUID, error)
	GetUserWithEmail(email, password string) (connectteam.User, error)
	GetIdWithEmail(email string) (uuid.UUID, error)
	//VerifyUser(verifyUser connectteam.VerifyUser) error
	GetVerificationCode(email string) (string, error)
	CreateVerificationCode(email string, code string) error
	CheckIfExist(id uuid.UUID) (bool, error)
	DeleteVerificationCode(email string, code string) error
}

type User interface {
	GetUserById(id uuid.UUID) (connectteam.UserPublic, error)
	UpdateAccessWithId(id uuid.UUID, access string) error
	GetUsersList() ([]connectteam.UserPublic, error)
	GetPassword(id uuid.UUID) (string, error)
	UpdatePasswordWithId(newPassword string, id uuid.UUID) error
	GetVerificationCode(email string) (string, error)
	GetEmailWithId(id uuid.UUID) (string, error)
	UpdateEmail(email string, id uuid.UUID) error
	CreateVerificationCode(email string, code string) error
	DeleteVerificationCode(email string, code string) error
	CheckIfExist(email string) (bool, error)
	UpdateUserFirstName(id uuid.UUID, firstName string) error
	UpdateUserSecondName(id uuid.UUID, secondName string) error
	UpdateUserDescription(id uuid.UUID, secondName string) error
	UpdateCompanyName(id uuid.UUID, companyName string) error
	UpdateCompanyInfo(id uuid.UUID, info string) error
	UpdateCompanyURL(id uuid.UUID, companyURL string) error
	GetUserPlan(userId uuid.UUID) (connectteam.UserPlan, error)
	CreatePlanRequest(request connectteam.PlanRequest) (uuid.UUID, error)
	GetUserCredentials(id uuid.UUID) (connectteam.UserCredentials, error)
	UpdatePasswordWithEmail(newPassword string, email string) error
}

type Plan interface {
	GetUserActivePlan(userId uuid.UUID) (connectteam.UserPlan, error)
	CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error)
	GetUsersPlans() ([]connectteam.UserPlan, error)
	SetConfirmed(id uuid.UUID) error
	DeletePlan(id uuid.UUID) error
	SetExpiredStatus(id uuid.UUID) error
	DeleteOnConfirmPlan(userId uuid.UUID) error
	GetUserSubscriptions(userId uuid.UUID) ([]connectteam.UserPlan, error)
	GetHolderWithInvitationCode(code string) (id uuid.UUID, err error)
	SetExpiredStatusWithUserId(userId uuid.UUID) error
	GetMembers(id uuid.UUID) (users []connectteam.UserPublic, err error)
	DeleteUserFromSub(userId uuid.UUID, planId uuid.UUID) error
	UpgradePlan(planId uuid.UUID, planType string, invitationCode string) error
	AddUserToSubscription(userId uuid.UUID, planId uuid.UUID, access string) error
	GetPlan(planId uuid.UUID) (sub connectteam.Subscription, err error)
}

type Topic interface {
	CreateTopic(topic connectteam.Topic) (uuid.UUID, error)
	GetAll() ([]connectteam.Topic, error)
	DeleteTopic(id uuid.UUID) error
	UpdateTopic(id uuid.UUID, title string) error
	GetTopic(topicId uuid.UUID) (topic connectteam.Topic, err error)
}

type Question interface {
	CreateQuestion(content string, topicId uuid.UUID) (uuid.UUID, error)
	DeleteQuestion(id uuid.UUID) error
	GetAll(topicId uuid.UUID) ([]connectteam.Question, error)
	UpdateQuestion(content string, id uuid.UUID) (connectteam.Question, error)
	GetRandWithLimit(topicId uuid.UUID, limit int) ([]connectteam.Question, error)
	GetQuestionTags(questionId uuid.UUID) ([]models.Tag, error)
	UpdateQuestionTags(questionId uuid.UUID, tags []models.Tag) ([]models.Tag, error)
	GetAllTags() ([]models.Tag, error)
}

type Game interface {
	CreateGame(game connectteam.Game) (connectteam.Game, error)
	GetCreatedGames(page int, userId uuid.UUID) (games []connectteam.Game, err error)
	CreateParticipant(userId uuid.UUID, gameId uuid.UUID) error
	GetGame(gameId uuid.UUID) (game connectteam.Game, err error)
	DeleteGame(gameId uuid.UUID) error
	GetGameWithInvitationCode(code string) (game connectteam.Game, err error)
	GetGames(page int, userId uuid.UUID) (games []connectteam.Game, err error)
	StartGame(gameId uuid.UUID) error
	SaveResults(gameId uuid.UUID, userId uuid.UUID, rate int) error
	GetResults(gameId uuid.UUID) (results []connectteam.UserResult, err error)
	EndGame(gameId uuid.UUID) error
	CancelGame(gameId uuid.UUID) error
	GetGameParticipant(gameId uuid.UUID) (users []connectteam.UserPublic, err error)
	ChangeStartDate(gameId uuid.UUID, date time.Time) error
	ChangeGameName(gameId uuid.UUID, name string) error
}

type Notification interface {
	GetNotifications(userId uuid.UUID) ([]models.Notification, error)
	CreateGameCancelNotification(gameId uuid.UUID, userId uuid.UUID) error
	CreateGameStartNotification(gameId uuid.UUID, userId uuid.UUID) error
	CreateGameInviteNotification(gameId uuid.UUID, userId uuid.UUID) error
	CreateSubInviteNotification(holderId uuid.UUID, userId uuid.UUID) error
	CreateDeleteFromSubNotification(holderId uuid.UUID, userId uuid.UUID) error
}

type Payment interface {
	CreatePayment(paymentRequest models.PaymentRequest) (models.PaymentResponse, error)
	GetPayment(orderId string) (models.PaymentResponse, error)
}

type Repository struct {
	Authorization
	User
	Plan
	Topic
	Question
	Game
	Notification
	Payment
}

func NewRepository(db *sqlx.DB, rdb *redis.Client, yooClient *yookassa.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
		Plan:          NewPlanPostgres(db),
		Topic:         NewTopicPostgres(db),
		Question:      NewQuestionPostgres(db),
		Game:          NewGamePostgres(db),
		Notification:  NewNotification(rdb),
		Payment:       NewYookassaClient(yooClient),
	}
}
