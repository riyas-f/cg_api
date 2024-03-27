package payload

type UserGames struct {
	Username string `db:"username"`
	GameID   string `db:"game_id"`
}
