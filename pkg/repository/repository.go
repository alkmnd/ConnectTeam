package repository

import (
	connectteam "ConnectTeam"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GetUserWithEmail(email, password string) (connectteam.User, error)
	GetIdWithEmail(email string) (int, error)
	GetUserWithPhone(phoneNumber,  password string) (connectteam.User, error)
	VerifyUser(verifyUser connectteam.VerifyUser)  error	
	GetVerificationCode(id int) (string, error)
	CreateVerificationCode(id int, code string) (error)
	CheckIfExist(id int) (bool, error)
	DeleteVerificationCode(id int, code string) (error)
}

type User interface {
	GetUserById(id int) (connectteam.UserPublic, error)
	UpdateAccessWithId(id int, access string) (error)
	GetUsersList() ([]connectteam.UserPublic, error)
	GetPassword(id int) (string, error)
	UpdatePassword(new_password string, id int) (error)
	GetVerificationCode(id int) (string, error)
	GetEmailWithId(id int) (string, error)
	UpdateEmail(email string, id int) (error)
	CreateVerificationCode(id int, code string) (error)
	DeleteVerificationCode(id int, code string) (error)
	CheckIfExist(email string) (bool, error)
	UpdateUserFirstName(id int, firstName string) (error)
	UpdateUserSecondName(id int, secondName string) (error)
	UpdateUserDescription(id int, secondName string) (error)
	UpdateCompanyName(id int, companyName string) (error)
	UpdateCompanyInfo(id int, info string) (error)
	UpdateCompanyURL(id int, companyURL string) (error)
	GetUserPlan(user_id int) (connectteam.UserPlan, error)
	CreatePlanRequest(request connectteam.PlanRequest) (int, error)
	GetUserCredentials(id int) (connectteam.UserCredentials, error)
}

type Plan interface {
	GetUserPlan(userId int) (connectteam.UserPlan, error)
	CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error)
	GetUsersPlans() ([] connectteam.UserPlan, error)
	SetConfirmed(id int) (error)
	DeletePlan(id int) (error)

}

type Repository struct {
	Authorization
	User
	Plan
}
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User: NewUserPostgres(db),
		Plan: NewPlanPostgres(db),
	}
}