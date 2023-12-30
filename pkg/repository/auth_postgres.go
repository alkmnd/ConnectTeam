package repository

import (
	connectteam "ConnectTeam"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user connectteam.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, phone_number, first_name, second_name, password_hash) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.PhoneNumber, user.FirstName, user.SecondName, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (r *AuthPostgres) GetUser(email, password string) (connectteam.User, error) {
	var user connectteam.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}