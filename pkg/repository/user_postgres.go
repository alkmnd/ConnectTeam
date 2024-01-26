package repository

import (
	connectteam "ConnectTeam"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}


func (r *UserPostgres) GetUserById(id int) (connectteam.UserPublic, error) {
	var user connectteam.UserPublic
	query := fmt.Sprintf("SELECT id, email, phone_number, first_name, second_name, access, company_name, profile_image FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *UserPostgres) ChangeAccessById(id int, access string) (error) {
	query := fmt.Sprintf("UPDATE %s SET access = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, access)
	
	return err
}

func (r *UserPostgres) GetUsersList() ([]connectteam.UserPublic, error) {
	var usersList []connectteam.UserPublic

	query := fmt.Sprintf("SELECT id, email, phone_number, first_name, second_name, access, company_name, profile_image FROM %s", usersTable)
	err := r.db.Select(&usersList, query)
	return usersList, err
}

func (r *UserPostgres) GetPassword(id int) (string, error) {
	var db_password string 
	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&db_password, query, id)
	if err != nil {
		return "", err
	}
	return db_password, nil
}

func (r *UserPostgres) ChangePassword(new_password string, id int) (error) {
	query := fmt.Sprintf("UPDATE %s SET password_hash = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, new_password)
	
	return err
}