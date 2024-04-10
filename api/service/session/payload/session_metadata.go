package payload

type SessionMetadata struct {
	MetadataID   int          `db:"metadata_id"`
	CreatedAt    string       `db:"created_at"`
	GameID       int          `json:"game_id" db:"game_id"`
	GameLocation GameLocation `json:"game_location"`
	SessionID    []byte       `json:"-" db:"session_id"`
}

type Server struct {
	IP   string `json:"ip"`
	Port int    `json:"port,string"`
	Host string `json:"-" db:"game_location_server_host"`
}

type GameLocation struct {
	Protocol string `json:"protocol" db:"game_location_protocol"`
	Path     string `json:"path" db:"game_location_path"`
	Server   Server `json:"server"`
}
