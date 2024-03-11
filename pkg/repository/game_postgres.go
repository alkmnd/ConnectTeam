package repository

import (
	connectteam "ConnectTeam"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type GamePostgres struct {
	db *sqlx.DB
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{db: db}
}

func (r *GamePostgres) CreateGame(game connectteam.Game) (connectteam.Game, error) {
	query := fmt.Sprintf("INSERT INTO %s (creator_id, name, start_date, invitation_code) values ($1, $2, $3, $4) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, game.CreatorId, game.Name, game.StartDate, game.InvitationCode)
	if err := row.Scan(&game.Id); err != nil {
		return game, err
	}
	return game, nil
}

func (r *GamePostgres) GetCreatedGames(limit int, offset int, userId int) (games []connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE creator_id=$1 ORDER BY start_date DESC LIMIT $2 OFFSET $3`, gamesTable)
	err = r.db.Select(&games, query, userId, limit, offset)
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

func (r *GamePostgres) DeleteGame(gameId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", gamesTable)
	_, err := r.db.Exec(query, gameId)
	return err
}

func (r *GamePostgres) GetGameWithInvitationCode(code string) (game connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE and invitation_code=$1 LIMIT 1`, gamesTable)
	err = r.db.Get(&game, query, code)
	return game, err
}

func (r *GamePostgres) GetGames(limit int, offset int, userId int) (games []connectteam.Game, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s g
	JOIN %s p ON p.game_id = g.id WHERE user_id =$1 ORDER BY start_date DESC LIMIT $2 OFFSET $3`, gamesTable, gamesUsersTable)
	err = r.db.Select(&games, query, userId, limit, offset)
	return games, err
}
