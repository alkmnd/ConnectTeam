package repository

import (
	"fmt"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable         = "users"
	codesTable         = "verification_codes"
	subscriptionsTable = "subscriptions"
	usersSubsTable     = "subs_holders"
	planRequestsTable  = "plan_requests"
	topicsTable        = "topics"
	questionsTable     = "questions"
	gamesTable         = "games"
	gamesUsersTable    = "games_users"
	resultsTable       = "results"
	tagsTable          = "tags"
	tagsQuestionsTable = "tags_questions"
	tagsUsersTable     = "tags_users"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
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
