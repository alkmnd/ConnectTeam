package repository

import (
	_ "github.com/lib/pq"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

const (
	codesTable = "verification_codes"
	plansUsersTable = "plans_users"
	planRequestsTable = "plan_requests"
	topicsTable = "topics"
)

type Config struct {
	Host string 
	Port string 
	Username string 
	Password string 
	DBName string 
	SSLMode string 
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", 
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}