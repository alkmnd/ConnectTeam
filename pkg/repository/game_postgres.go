package repository

import (
	connectteam "ConnectTeam"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
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

func (r *GamePostgres) GetCreatedGames(page int, userId int) (games []connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE creator_id=$1 ORDER BY start_date DESC LIMIT $2 OFFSET $3`, gamesTable)
	err = r.db.Select(&games, query, userId, limit, page*limit)
	return games, err
}

func (r *GamePostgres) CreateParticipant(userId int, gameId int) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, game_id) VALUES ($1, $2)
		RETURNING *`, gamesUsersTable)

	_, err := r.db.Exec(query, userId, gameId)
	return err
}

func (r *GamePostgres) GetGame(gameId int) (game connectteam.Game, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", gamesTable)
	err = r.db.Get(&game, query, gameId)
	return game, err
}

func (r *GamePostgres) StartGame(gameId int) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'in_progress' WHERE id = %d`, gamesTable, gameId)

	_, err := r.db.Exec(query)

	return err
}

func (r *GamePostgres) DeleteGame(gameId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", gamesTable)
	_, err := r.db.Exec(query, gameId)
	return err
}

func (r *GamePostgres) GetGameWithInvitationCode(code string) (game connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE invitation_code=$1 LIMIT 1`, gamesTable)
	err = r.db.Get(&game, query, code)
	return game, err
}

func (r *GamePostgres) GetGames(page int, userId int) (games []connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT id, creator_id, invitation_code, name, start_date, status FROM %s g
	JOIN %s p ON p.game_id = g.id WHERE p.user_id=$1 ORDER BY start_date DESC LIMIT $2 OFFSET $3`, gamesTable, gamesUsersTable)
	log.Println(limit * page)
	err = r.db.Select(&games, query, userId, limit, limit*page)
	return games, err
}
