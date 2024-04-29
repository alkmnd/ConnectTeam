package repository

import (
	"ConnectTeam"
	"fmt"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type Role string

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) CreateUser(user connectteam.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (email, first_name, second_name, description, password_hash, access, profile_image, company_name, company_info, company_url, company_logo) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.FirstName, user.SecondName, "", user.Password, "user", "", "", "", "", "")
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateVerificationCode(email string, code string) error {
	query := fmt.Sprintf("INSERT INTO %s (email, code) values ($1, $2) ON CONFLICT (email) DO UPDATE SET code = $2", codesTable)
	_, err := r.db.Exec(query, email, code)
	return err
}

func (r *AuthPostgres) GetUserWithEmail(email, password string) (connectteam.User, error) {
	var user connectteam.User
	query := fmt.Sprintf("SELECT id, email, first_name, second_name, access FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	if err := r.db.Get(&user, query, email, password); err != nil {
		print(err.Error())
		return user, err
	}
	return user, nil
}

func (r *AuthPostgres) GetIdWithEmail(email string) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("SELECT id  FROM %s WHERE email=$1", usersTable)
	if err := r.db.Get(&id, query, email); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) GetVerificationCode(email string) (string, error) {
	var code string
	query := fmt.Sprintf("SELECT code from %s WHERE email = $1", codesTable)
	err := r.db.Get(&code, query, email)

	return code, err
}

func (r *AuthPostgres) CheckIfExist(id uuid.UUID) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&count, query, id)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthPostgres) DeleteVerificationCode(email string, code string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email = $1 AND code=$2", codesTable)
	_, err := r.db.Exec(query, email, code)
	return err
}
