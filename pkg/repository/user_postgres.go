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
	query := fmt.Sprintf("SELECT id, email, first_name, second_name, access, company_name, profile_image FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *UserPostgres) ChangeAccessWithId(id int, access string) (error) {
	query := fmt.Sprintf("UPDATE %s SET access = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, access)
	
	return err
}

func (r *UserPostgres) GetUsersList() ([]connectteam.UserPublic, error) {
	var usersList []connectteam.UserPublic

	query := fmt.Sprintf("SELECT id, email, first_name, second_name, access, company_name, profile_image FROM %s", usersTable)
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

func (r *UserPostgres) ChangeEmail(email string, id int) (error) {
	query := fmt.Sprintf("UPDATE %s SET email = $1 WHERE id = %d", usersTable, id)
	_, err := r.db.Exec(query, email)

	return err
}

func (r *UserPostgres) GetVerificationCode(id int) (string, error) {
	var code string
	query := fmt.Sprintf("SELECT code from %s WHERE user_id = $1", codesTable)
	err := r.db.Get(&code, query, id)
	println("meow " + code)

	return code, err
}

func (r *UserPostgres) CreateVerificationCode(user_id int, code string) (error){
	query := fmt.Sprintf("INSERT INTO %s (user_id, code) values ($1, $2) ON CONFLICT (user_id) DO UPDATE SET code = $2", codesTable)
	_, err := r.db.Exec(query, user_id, code)
	return err
}

func (r *UserPostgres) DeleteVerificationCode(id int, code string) (error) {
	println(id)
	println("postgres: "+code)
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND code=$2", codesTable)
	_, err := r.db.Exec(query, id, code)
	return err
}

func (r *UserPostgres) GetEmailWithId(id int) (string, error) {
	var email string
	query := fmt.Sprintf("SELECT email from %s WHERE user_id = $1", usersTable)
	err := r.db.Get(&email, query, id)

	return email, err
}

func (r *UserPostgres) CheckIfExist(email string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", usersTable)
	err := r.db.Get(&count, query, email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserPostgres) ChangeUserFirstName(id int, firstName string) (error) {
	query := fmt.Sprintf("UPDATE %s SET first_name = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, firstName)
	
	return err
}

func (r *UserPostgres) ChangeUserSecondName(id int, secondName string) (error) {
	query := fmt.Sprintf("UPDATE %s SET second_name = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, secondName)
	
	return err
}

func (r *UserPostgres) ChangeUserDescription(id int, description string) (error) {
	query := fmt.Sprintf("UPDATE %s SET description = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, description)
	
	return err
}
