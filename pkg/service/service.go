package service

import (
	connectteam "ConnectTeam/models"
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/repository/filestorage"
	"ConnectTeam/pkg/service/models"
	"ConnectTeam/pkg/service/uploader"
	"context"
	"github.com/google/uuid"
	"io"
	//	"io"
)

type Authorization interface {
	CreateUser(user connectteam.UserSignUpRequest) (uuid.UUID, error)
	GenerateToken(login, password string, isEmail bool) (string, string, error, uuid.UUID)
	ParseToken(token string) (uuid.UUID, string, error)
	//VerifyPhone(verifyPhone connectteam.VerifyPhone) (string, error)
	//VerifyUser(verifyUser connectteam.VerifyUser) error
	VerifyEmail(verifyEmail connectteam.VerifyEmail) error
	DeleteVerificationCode(email string, code string) error
}

type User interface {
	GetUserById(id uuid.UUID) (connectteam.UserPublic, error)
	UpdateAccessWithId(id uuid.UUID, access connectteam.AccessLevel) error
	GetUsersList() ([]connectteam.UserPublic, error)
	UpdatePassword(oldPassword string, newPassword string, id uuid.UUID) error
	UpdateEmail(id uuid.UUID, newEmail string, code string) error
	DeleteVerificationCode(email string, code string) error
	CheckEmailOnChange(id uuid.UUID, email string, password string) error
	UpdatePersonalData(id uuid.UUID, user connectteam.UserPersonalInfo) error
	UpdateCompanyData(id uuid.UUID, company connectteam.UserCompanyData) error
	GetUserPlan(userId uuid.UUID) (connectteam.UserPlan, error)
	RestorePasswordAuthorized(id uuid.UUID) error
	RestorePassword(email string) error
	CheckIfExist(email string) (bool, error)
}

type Plan interface {
	GetUserActivePlan(userId uuid.UUID) (connectteam.UserPlan, error)
	CreatePlan(orderId string, userId uuid.UUID) (connectteam.UserPlan, error)
	GetUsersPlans() ([]connectteam.UserPlan, error)
	ConfirmPlan(id uuid.UUID) error
	SetPlanByAdmin(userId uuid.UUID, planType string, expiryDateString string) error
	DeletePlan(id uuid.UUID) error
	CheckIfSubscriptionExists(userId uuid.UUID) (bool, error)
	CreateTrialPlan(userId uuid.UUID) (userPlan connectteam.UserPlan, err error)
	GetUserSubscriptions(userId uuid.UUID) ([]connectteam.UserPlan, error)
	GetHolderWithInvitationCode(code string) (id uuid.UUID, err error)
	AddUserToAdvanced(holderPlan connectteam.UserPlan, userId uuid.UUID) (userPlan connectteam.UserPlan, err error)
	GetMembers(id uuid.UUID) ([]connectteam.UserPublic, error)
	DeleteUserFromSub(userId uuid.UUID, planId uuid.UUID) error
	UpgradePlan(orderId string, planId uuid.UUID, userId uuid.UUID) error
	InviteUserToSub(planId uuid.UUID, userId uuid.UUID, holderId uuid.UUID) error
	GetPlan(planId uuid.UUID) (sub connectteam.Subscription, err error)
}

type Topic interface {
	CreateTopic(topic connectteam.Topic) (uuid.UUID, error)
	GetAll() ([]connectteam.Topic, error)
	DeleteTopic(id uuid.UUID) error
	UpdateTopic(id uuid.UUID, title string) error
	GetTopic(id uuid.UUID) (connectteam.Topic, error)
	GetRandWithLimit(limit int) (topics []connectteam.Topic, err error)
}

type Question interface {
	CreateQuestion(content string, topicId uuid.UUID) (uuid.UUID, error)
	DeleteQuestion(id uuid.UUID) error
	GetAll(topicId uuid.UUID) ([]models.Question, error)
	UpdateQuestion(content string, id uuid.UUID) (connectteam.Question, error)
	GetRandWithLimit(topicId uuid.UUID, limit int) ([]models.Question, error)
	GetAllTags() ([]models.Tag, error)
	UpdateQuestionTags(questionId uuid.UUID, tags []models.Tag) ([]models.Tag, error)
	GetTagsUsers(userId uuid.UUID, gameId uuid.UUID) ([]models.Tag, error)
	CreateTagsUsers(userId uuid.UUID, gameId uuid.UUID, tagId uuid.UUID) error
}

type Payment interface {
	CreatePayment(userId uuid.UUID, plan string, returnURL string) (models.PaymentResponse, error)
}

type Uploader interface {
	Upload(ctx context.Context, file io.Reader, size int64, contentType string) (string, error)
}

type Game interface {
	CreateGame(creatorId uuid.UUID, startDateString string, name string) (connectteam.Game, error)
	GetCreatedGames(page int, userId uuid.UUID) ([]connectteam.Game, error)
	CreateParticipant(userId uuid.UUID, gameId uuid.UUID) error
	GetGame(gameId uuid.UUID) (connectteam.Game, error)
	DeleteGameFromGameList(gameId uuid.UUID, userId uuid.UUID) error
	GetGameWithInvitationCode(code string) (connectteam.Game, error)
	GetGames(page int, userId uuid.UUID) ([]connectteam.Game, error)
	GetResults(gameId uuid.UUID) (results []connectteam.UserResult, err error)
	StartGame(gameId uuid.UUID) error
	EndGame(gameId uuid.UUID) error
	SaveResults(gameId uuid.UUID, userId uuid.UUID, rate int) error
	CancelGame(gameId uuid.UUID, userId uuid.UUID) error
	InviteUserToGame(gameId uuid.UUID, userId uuid.UUID, creatorId uuid.UUID) error
	ChangeStartDate(gameId uuid.UUID, dateString string) error
	ChangeGameName(gameId uuid.UUID, name string) error
	GetGameParticipants(gameId uuid.UUID) (users []connectteam.UserPublic, err error)
}

type Notification interface {
	GetUserNotifications(userId uuid.UUID) (notifications []models.Notification, err error)
}

type Service struct {
	Authorization
	User
	Plan
	Topic
	Question
	Uploader
	Game
	Payment
	Notification
}

func NewService(repos *repository.Repository, fileStorage *filestorage.FileStorage) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
		Plan:          NewPlanService(repos.Plan, repos.Payment, repos.Notification),
		Topic:         NewTopicService(repos.Topic),
		Question:      NewQuestionService(repos.Question),
		Uploader:      uploader.NewUploader(fileStorage),
		Game:          NewGameService(repos.Game, repos.Notification, repos.Plan),
		Payment:       NewPaymentService(repos.Payment),
		Notification:  NewNotificationService(repos.Notification),
	}
}
