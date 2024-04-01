package payload

type UserGames struct {
	Username string `db:"username"`
	GameID   int    `db:"game_id"`
}
