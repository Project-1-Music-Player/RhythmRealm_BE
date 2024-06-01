package models

import "time"

type Song struct {
	SongID       string    `json:"song_id"`
	Title        string    `json:"title"`
	UserID       string    `json:"user_id"`
	Album        string    `json:"album"`
	ReleaseDate  time.Time `json:"release_date"`
	Genre        string    `json:"genre"`
	SongURL      string    `json:"song_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	PlayCount    int       `json:"play_count"`
}
type SongUpload struct {
	SongID       string    `json:"song_id"`
	Title        string    `json:"title"`
	UserID       string    `json:"user_id"`
	Album        string    `json:"album"`
	ReleaseDate  time.Time `json:"release_date"`
	Genre        string    `json:"genre"`
	SongUrl      string    `json:"song_url"`
	ThumbnailUrl string    `json:"thumbnail_url"`
}
