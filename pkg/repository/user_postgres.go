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


func (r *UserPostgres) GetUserById(id int) (connectteam.User, error) {
	var user connectteam.User
	query := fmt.Sprintf("SELECT id, email, phone_number, first_name, second_name, is_verified, access FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *UserPostgres) ChangeAccessById(id int, access string) (error) {
	query := fmt.Sprintf("UPDATE %s SET access = $1 WHERE id = %d", usersTable, id)

	_, err := r.db.Exec(query, access)
	
	return err
	
}