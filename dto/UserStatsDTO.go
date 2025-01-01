package dto

type UserStatsDTO struct {
	TotalSongs int64                 `json:"total_songs"`
	TopSongs   []SongRequestCountDTO `json:"top_songs"`
}
