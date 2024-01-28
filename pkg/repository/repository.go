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
}

type UserInterface interface {
	GetUserById(id int) (connectteam.UserPublic, error)
	ChangeAccessById(id int, access string) (error)
	GetUsersList() ([]connectteam.UserPublic, error)
	GetPassword(id int) (string, error)
	ChangePassword(new_password string, id int) (error)

}


type Repository struct {
	Authorization
	UserInterface
}
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		UserInterface: NewUserPostgres(db),
	}
}