package models

type SongPlayRequest struct {
	UserID   string `json:"user_id"`
	GuildID  string `json:"guild_id"`
	SongName string `json:"song_name"`
}
