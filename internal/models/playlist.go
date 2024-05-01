package models

import "github.com/gocql/gocql"

type Playlist struct {
	PlaylistID  gocql.UUID `json:"playlist_id"`
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}
