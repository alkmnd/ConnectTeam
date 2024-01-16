package repository

import (
	connectteam "ConnectTeam"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Role string

const (
	Admin Role = "Admin"
	User = "User"
)
type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) CreateUser(user connectteam.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, phone_number, first_name, second_name, password_hash, role) values ($1, $2, $3, $4, $5, $6) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.PhoneNumber, user.FirstName, user.SecondName, user.Password, "user")
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (r *AuthPostgres) GetUserWithEmail(email, password string) (connectteam.User, error) {
	var user connectteam.User
	query := fmt.Sprintf("SELECT id, email, phone_number, first_name, second_name, role, is_verified  FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	if err := r.db.Select(&user, query, email, password); err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthPostgres) GetUserWithPhone(phoneNumber, password string) (connectteam.User, error) {
	var user connectteam.User
	println(password)
	// query := fmt.Sprintf("SELECT id, email, phone_number, first_name, second_name, role, is_verified FROM %s WHERE phone_number=$1 AND password_hash=$2 LIMIT 1", usersTable)
	// if err := r.db.Select(&user, query, phoneNumber, password); err != nil {
	// 	println(err.Error())
	// 	return user, err
	// 
	query := fmt.Sprintf("SELECT id, email, phone_number, first_name, second_name, is_verified, role FROM %s WHERE phone_number=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, phoneNumber, password)
	println(user.Id)
	println(user.Email)
	return user, err
}

func (r *AuthPostgres) VerifyUser(verifyUser connectteam.VerifyUser) error {
	query := fmt.Sprintf("UPDATE %s SET is_verified = true WHERE id = %d", usersTable, verifyUser.Id)

	_, err := r.db.Exec(query)
	
	return err

}
func (r *AuthPostgres) Verify(verifyUser connectteam.VerifyUser) error {
	query := fmt.Sprintf("UPDATE %s SET is_verified = true WHERE id = %d", usersTable, verifyUser.Id)

	_, err := r.db.Exec(query)
	
	return err
}