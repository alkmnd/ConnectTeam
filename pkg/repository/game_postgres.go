package repository

import (
	connectteam "ConnectTeam/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type GamePostgres struct {
	db *sqlx.DB
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{db: db}
}

const (
	limit = 10
)

func (r *GamePostgres) CreateGame(game connectteam.Game) (connectteam.Game, error) {
	query := fmt.Sprintf("INSERT INTO %s (creator_id, name, start_date, invitation_code, status) values ($1, $2, $3, $4, $5) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, game.CreatorId, game.Name, game.StartDate, game.InvitationCode, game.Status)
	if err := row.Scan(&game.Id); err != nil {
		return game, err
	}
	return game, nil
}

func (r *GamePostgres) ChangeStartDate(gameId uuid.UUID, date time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET start_date = $1 WHERE id = $2`, gamesTable)

	_, err := r.db.Exec(query, date, gameId)

	return err
}

func (r *GamePostgres) ChangeGameName(gameId uuid.UUID, name string) error {
	query := fmt.Sprintf(`UPDATE %s SET name = $1 WHERE id = $2`, gamesTable)

	_, err := r.db.Exec(query, name, gameId)

	return err
}

func (r *GamePostgres) SaveResults(gameId uuid.UUID, userId uuid.UUID, rate int, name string) error {
	query := fmt.Sprintf("INSERT INTO %s (game_id, user_id, value, name) values ($1, $2, $3, $4)", resultsTable)
	_, err := r.db.Exec(query, gameId, userId, rate, name)
	return err
}

func (r *GamePostgres) GetCreatedGames(page int, userId uuid.UUID) (games []connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE creator_id=$1 ORDER BY start_date DESC LIMIT $2 OFFSET $3`, gamesTable)
	err = r.db.Select(&games, query, userId, limit, page*limit)
	return games, err
}

func (r *GamePostgres) CreateParticipant(userId uuid.UUID, gameId uuid.UUID) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, game_id) VALUES ($1, $2) ON CONFLICT (user_id, game_id) DO NOTHING
		RETURNING *`, gamesUsersTable)

	_, err := r.db.Exec(query, userId, gameId)
	return err
}

func (r *GamePostgres) GetGame(gameId uuid.UUID) (game connectteam.Game, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", gamesTable)
	err = r.db.Get(&game, query, gameId)
	return game, err
}

func (r *GamePostgres) GetResults(gameId uuid.UUID) (results []connectteam.UserResult, err error) {

	query := fmt.Sprintf(`SELECT user_id, value FROM %s WHERE game_id=$1`, resultsTable)
	err = r.db.Select(&results, query, gameId)
	return results, err
}

func (r *GamePostgres) StartGame(gameId uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'in_progress' WHERE id = $1`, gamesTable)

	_, err := r.db.Exec(query, gameId)

	return err
}

func (r *GamePostgres) EndGame(gameId uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'ended' WHERE id = $1`, gamesTable)

	_, err := r.db.Exec(query, gameId)

	return err
}

func (r *GamePostgres) DeleteGameFromGameList(gameId uuid.UUID, userId uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE game_id = $1 AND user_id=$2", gamesUsersTable)
	_, err := r.db.Exec(query, gameId, userId)
	return err
}

func (r *GamePostgres) GetGameWithInvitationCode(code string) (game connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE invitation_code=$1 LIMIT 1`, gamesTable)
	err = r.db.Get(&game, query, code)
	return game, err
}

func (r *GamePostgres) GetGames(page int, userId uuid.UUID) (games []connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT id, creator_id, invitation_code, name, start_date, status FROM %s g
	JOIN %s p ON p.game_id = g.id WHERE p.user_id=$1 ORDER BY start_date DESC LIMIT $2 OFFSET $3`, gamesTable, gamesUsersTable)
	err = r.db.Select(&games, query, userId, limit, limit*page)
	return games, err
}

func (r *GamePostgres) CancelGame(gameId uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'cancelled' WHERE id = $1`, gamesTable)

	_, err := r.db.Exec(query, gameId)

	return err
}

func (r *GamePostgres) GetGameParticipants(gameId uuid.UUID) (users []connectteam.UserPublic, err error) {
	query := fmt.Sprintf(`SELECT u.id FROM %s u 
    JOIN %s gu ON gu.user_id = u.id
	 WHERE gu.game_id=$1`, usersTable, gamesUsersTable)
	err = r.db.Select(&users, query, gameId)
	return users, err
}
