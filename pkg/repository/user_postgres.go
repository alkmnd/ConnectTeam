package repository

import (
	connectteam "ConnectTeam"
	"fmt"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUserById(id uuid.UUID) (connectteam.UserPublic, error) {
	var user connectteam.UserPublic
	query := fmt.Sprintf("SELECT id, email, first_name, second_name, description, access, company_name, company_info, company_url, company_logo, profile_image FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	println(err)
	return user, err
}

func (r *UserPostgres) GetUserCredentials(id uuid.UUID) (connectteam.UserCredentials, error) {
	var userCred connectteam.UserCredentials
	query := fmt.Sprintf("SELECT email, password_hash FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&userCred, query, id)
	return userCred, err
}

func (r *UserPostgres) UpdateAccessWithId(id uuid.UUID, access string) error {
	query := fmt.Sprintf("UPDATE %s SET access = $1 WHERE id = %1", usersTable)

	_, err := r.db.Exec(query, access, id)

	return err
}

func (r *UserPostgres) GetUsersList() ([]connectteam.UserPublic, error) {
	var usersList []connectteam.UserPublic

	query := fmt.Sprintf("SELECT id, email, first_name, second_name, access, company_name, profile_image FROM %s", usersTable)
	err := r.db.Select(&usersList, query)
	return usersList, err
}

func (r *UserPostgres) GetPassword(id uuid.UUID) (string, error) {
	var dbPassword string
	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&dbPassword, query, id)
	if err != nil {
		return "", err
	}
	return dbPassword, nil
}

func (r *UserPostgres) UpdatePasswordWithId(newPassword string, id uuid.UUID) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, newPassword, id)

	return err
}

func (r *UserPostgres) UpdatePasswordWithEmail(newPassword string, email string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash = $1 WHERE email = $2", usersTable)

	_, err := r.db.Exec(query, newPassword, email)

	return err
}

func (r *UserPostgres) UpdateEmail(email string, id uuid.UUID) error {
	query := fmt.Sprintf("UPDATE %s SET email = $1 WHERE id = $1", usersTable)
	_, err := r.db.Exec(query, email, id)

	return err
}

func (r *UserPostgres) GetVerificationCode(email string) (string, error) {
	var code string
	query := fmt.Sprintf("SELECT code from %s WHERE email = $1", codesTable)
	err := r.db.Get(&code, query, email)

	return code, err
}

func (r *UserPostgres) CreateVerificationCode(email string, code string) error {
	query := fmt.Sprintf("INSERT INTO %s (email, code) values ($1, $2) ON CONFLICT (email) DO UPDATE SET code = $2", codesTable)
	_, err := r.db.Exec(query, email, code)
	return err
}

func (r *UserPostgres) DeleteVerificationCode(email string, code string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email = $1 AND code=$2", codesTable)
	_, err := r.db.Exec(query, email, code)
	return err
}

func (r *UserPostgres) GetEmailWithId(id uuid.UUID) (string, error) {
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

// func (r *AuthPostgres) CheckIfExist(id int) (bool, error) {
// 	var count int
// 	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1", usersTable)
// 	err := r.db.Get(&count, query, id)

// 	if err != nil {
// 		return false, err
// 	}

// 	return count > 0, nil
// }

func (r *UserPostgres) UpdateUserFirstName(id uuid.UUID, firstName string) error {
	query := fmt.Sprintf("UPDATE %s SET first_name = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, firstName, id)

	return err
}

func (r *UserPostgres) UpdateUserSecondName(id uuid.UUID, secondName string) error {
	query := fmt.Sprintf("UPDATE %s SET second_name = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, secondName, id)

	return err
}

func (r *UserPostgres) UpdateUserDescription(id uuid.UUID, description string) error {
	query := fmt.Sprintf("UPDATE %s SET description = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, description, id)

	return err
}

func (r *UserPostgres) UpdateCompanyName(id uuid.UUID, companyName string) error {
	query := fmt.Sprintf("UPDATE %s SET company_name = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, companyName, id)

	return err
}

func (r *UserPostgres) UpdateCompanyInfo(id uuid.UUID, info string) error {
	query := fmt.Sprintf("UPDATE %s SET company_info = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, info, id)

	return err
}

func (r *UserPostgres) UpdateCompanyURL(id uuid.UUID, url string) error {
	query := fmt.Sprintf("UPDATE %s SET company_url = $1 WHERE id = $1", usersTable)

	_, err := r.db.Exec(query, url, id)

	return err
}

func (r *UserPostgres) GetUserPlan(userId uuid.UUID) (connectteam.UserPlan, error) {
	var userPlan connectteam.UserPlan
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", subscriptionsTable)
	err := r.db.Get(&userPlan, query, userId)

	return userPlan, err
}

func (r *UserPostgres) CreatePlanRequest(request connectteam.PlanRequest) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (user_id, duration, request_date, plan_type) VALUES ($1, $2, $3, $4) RETURNING id", planRequestsTable)

	row := r.db.QueryRow(query, request.UserId, request.Duration, request.RequestDate, request.PlanType)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
