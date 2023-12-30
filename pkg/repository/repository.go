package repository

import (
	connectteam "ConnectTeam"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user connectteam.User) (int, error)
	GetUserWithEmail(email, password string) (connectteam.User, error)
	GetUserWithPhone(phoneNumber,  password string) (connectteam.User, error)
}

type Repository struct {
	Authorization
}
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}