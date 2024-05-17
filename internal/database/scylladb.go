package database

import (
	"log"
	"os"
	"rr-backend/internal/models"
	"time"

	"github.com/gocql/gocql"
)

type ScyllaService interface {
	Health() map[string]string
	UpsertUser(userID, username, email, role string) error

	InsertSong(songID gocql.UUID, title, userID, album string, releaseDate time.Time, genre, songURL, thumbnailURL string) error
	RemoveSong(songID gocql.UUID) error
	GetSongsByUserID(userID string) ([]models.Song, error)
	GetObjectNameBySongID(songID string) (string, error)
	SearchSongs(query string, limit, offset int) ([]models.Song, error)

	AddPlaylist(playlistID gocql.UUID, userID, name, description string) error
	AddSongToPlaylist(playlistID gocql.UUID, userID string, songID gocql.UUID, addedAt time.Time) error
	RemoveSongFromPlaylist(playlistID gocql.UUID, songID gocql.UUID) error
	FetchPlaylists(userID string) ([]models.Playlist, error)
	RemovePlaylist(playlistID gocql.UUID) error

	IncrementPlayCount(songID string) error
}

type scyllaService struct {
	session *gocql.Session
}

func NewScylla() ScyllaService {
	cluster := gocql.NewCluster(os.Getenv("DB_HOST"))
	cluster.Port = 9042
	cluster.Keyspace = os.Getenv("DB_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Cannot connect to ScyllaDB:", err)
	}

	s := &scyllaService{
		session: session,
	}
	return s
}

func (s *scyllaService) Health() map[string]string {
	if err := s.session.Query(`SELECT now() FROM system.local`).Exec(); err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *scyllaService) UpsertUser(userID, username, email, role string) error {
	query := `INSERT INTO users (user_id, username, email, role) VALUES (?, ?, ?, ?)`
	if err := s.session.Query(query, userID, username, email, role).Exec(); err != nil {
		log.Printf("Failed to upsert user: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) InsertSong(songID gocql.UUID, title, userID, album string, releaseDate time.Time, genre, songURL, thumbnailURL string) error {
	query := `INSERT INTO songs (song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if err := s.session.Query(query, songID, title, userID, album, releaseDate, genre, songURL, thumbnailURL).Exec(); err != nil {
		log.Printf("Failed to insert song: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) RemoveSong(songID gocql.UUID) error {
	query := `DELETE FROM songs WHERE song_id = ?`
	if err := s.session.Query(query, songID).Exec(); err != nil {
		log.Printf("Failed to remove song: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) GetSongsByUserID(userID string) ([]models.Song, error) {
	query := `SELECT song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url FROM songs WHERE user_id = ?`
	iter := s.session.Query(query, userID).Iter()

	var songs []models.Song
	var song models.Song
	for iter.Scan(&song.SongID, &song.Title, &song.UserID, &song.Album, &song.ReleaseDate, &song.Genre, &song.SongURL, &song.ThumbnailURL) {
		songs = append(songs, song)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *scyllaService) GetObjectNameBySongID(songID string) (string, error) {
	var objectName string
	query := `SELECT song_url FROM songs WHERE song_id = ? LIMIT 1`
	if err := s.session.Query(query, songID).Scan(&objectName); err != nil {
		return "", err
	}
	return objectName, nil
}

func (s *scyllaService) SearchSongs(query string, limit, offset int) ([]models.Song, error) {
	var songs []models.Song
	cqlQuery := "SELECT song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url FROM songs WHERE title LIKE ? ALLOW FILTERING"

	iter := s.session.Query(cqlQuery, "%"+query+"%").PageSize(limit).PageState(nil).Iter() // create an iterator that go through the select results

	var song models.Song
	for iter.Scan(&song.SongID, &song.Title, &song.UserID, &song.Album, &song.ReleaseDate, &song.Genre, &song.SongURL, &song.ThumbnailURL) {
		songs = append(songs, song)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *scyllaService) AddPlaylist(playlistID gocql.UUID, userID, name, description string) error {
	query := `INSERT INTO playlists (playlist_id, user_id, name, description) VALUES (?, ?, ?, ?)`
	if err := s.session.Query(query, playlistID, userID, name, description).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *scyllaService) AddSongToPlaylist(playlistID gocql.UUID, userID string, songID gocql.UUID, addedAt time.Time) error {
	query := `INSERT INTO playlist_songs (playlist_id, added_at, song_id) VALUES (?, ?, ?)`
	if err := s.session.Query(query, playlistID, addedAt, songID).Exec(); err != nil {
		log.Printf("Failed to add song to playlist: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) RemoveSongFromPlaylist(playlistID gocql.UUID, songID gocql.UUID) error {
	query := `DELETE FROM playlist_songs WHERE playlist_id = ? AND song_id = ?`
	if err := s.session.Query(query, playlistID, songID).Exec(); err != nil {
		log.Printf("Failed to remove song from playlist: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) FetchPlaylists(userID string) ([]models.Playlist, error) {
	query := `SELECT playlist_id, name, description FROM playlists WHERE user_id = ?`
	iter := s.session.Query(query, userID).Iter()

	var playlists []models.Playlist
	var playlistID gocql.UUID
	var name string
	var description string

	for iter.Scan(&playlistID, &name, &description) {
		playlists = append(playlists, models.Playlist{
			PlaylistID:  playlistID,
			UserID:      userID,
			Name:        name,
			Description: description,
		})
	}

	if err := iter.Close(); err != nil {
		log.Printf("Failed to fetch playlists: %v", err)
		return nil, err
	}

	return playlists, nil
}

func (s scyllaService) RemovePlaylist(playlistID gocql.UUID) error {
	batch := s.session.NewBatch(gocql.LoggedBatch)
	batch.Query(`DELETE FROM playlists WHERE playlist_id = ?`, playlistID)
	batch.Query(`DELETE FROM playlist_songs WHERE playlist_id = ?`, playlistID)
	if err := s.session.ExecuteBatch(batch); err != nil {
		log.Printf("Failed to remove playlists: %v", err)
		return err
	}
	return nil
}

func (s scyllaService) IncrementPlayCount(songID string) error {
	query := `UPDATE songs SET play_count = play_count + 1 WHERE song_id = ?`
	if err := s.session.Query(query, songID).Exec(); err != nil {
		log.Printf("Failed to increment play count: %v", err)
		return err
	}
	return nil
}
