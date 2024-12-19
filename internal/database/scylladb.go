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
	GetUserByID(userID string) (*models.User, error)
	UpdateUserRole(userID, role string) error

	InsertSong(songID gocql.UUID, title, userID, album string, releaseDate time.Time, genre, songURL, thumbnailURL string) error
	RemoveSong(songID gocql.UUID) error
	GetSongsByUserID(userID string) ([]models.Song, error)
	GetAllSongs() ([]models.Song, error)
	GetObjectNameBySongID(songID string) (string, error)
	GetSongThumbnailBySongID(songID string) (string, error)
	SearchSongs(query string, limit, offset int) ([]models.Song, error)

	AddPlaylist(playlistID gocql.UUID, userID, name, description string) error
	UpdatePlaylist(playlistID gocql.UUID, name, description string) error
	AddSongToPlaylist(playlistID gocql.UUID, userID string, songID gocql.UUID, addedAt time.Time) error
	RemoveSongFromPlaylist(playlistID gocql.UUID, songID gocql.UUID) error
	GetSongsInPlaylist(playlistID gocql.UUID) ([]models.Song, error)
	FetchPlaylists(userID string) ([]models.Playlist, error)
	RemovePlaylist(playlistID gocql.UUID) error

	LikeSong(userID string, songID gocql.UUID) error
	UnlikeSong(userID string, songID gocql.UUID) error
	GetLikedSongsByUser(userID string) ([]models.Song, error)

	GetSongUserID(songID gocql.UUID) (string, error)

	GetAllArtists() ([]models.Artist, error)
	GetArtistWithSongs(artistID string) (*models.ArtistWithSongs, error)
	FollowArtist(artistID string, followerID string) error
	UnfollowArtist(artistID string, followerID string) error
	GetFollowedArtists(userID string) ([]models.Artist, error)
	GetArtistFollowersCount(artistID string) (int, error)
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
	query := `INSERT INTO songs (song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url, play_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)`
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
	query := `SELECT song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url, play_count FROM songs WHERE user_id = ?`
	iter := s.session.Query(query, userID).Iter()

	var songs []models.Song
	var song models.Song
	for iter.Scan(&song.SongID, &song.Title, &song.UserID, &song.Album, &song.ReleaseDate, &song.Genre, &song.SongURL, &song.ThumbnailURL, &song.PlayCount) {
		songs = append(songs, song)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *scyllaService) GetAllSongs() ([]models.Song, error) {
	query := `SELECT song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url, play_count FROM songs`
	iter := s.session.Query(query).Iter()

	var songs []models.Song
	var song models.Song
	for iter.Scan(&song.SongID, &song.Title, &song.UserID, &song.Album, &song.ReleaseDate, &song.Genre, &song.SongURL, &song.ThumbnailURL, &song.PlayCount) {
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

func (s *scyllaService) GetSongThumbnailBySongID(songID string) (string, error) {
	var thumbnailURL string
	query := `SELECT thumbnail_url FROM songs WHERE song_id = ? LIMIT 1`
	if err := s.session.Query(query, songID).Scan(&thumbnailURL); err != nil {
		return "", err
	}
	return thumbnailURL, nil
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

func (s *scyllaService) UpdatePlaylist(playlistID gocql.UUID, name, description string) error {
	query := `UPDATE playlists SET name = ?, description = ? WHERE playlist_id = ?`
	if err := s.session.Query(query, name, description, playlistID).Exec(); err != nil {
		log.Printf("Failed to update playlist: %v", err)
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

func (s *scyllaService) GetSongsInPlaylist(playlistID gocql.UUID) ([]models.Song, error) {
	query := `SELECT song_id FROM playlist_songs WHERE playlist_id = ?`
	iter := s.session.Query(query, playlistID).Iter()

	var songIDs []gocql.UUID
	var songID gocql.UUID
	for iter.Scan(&songID) {
		songIDs = append(songIDs, songID)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	if len(songIDs) == 0 {
		return []models.Song{}, nil
	}

	query = `SELECT song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url, play_count FROM songs WHERE song_id IN ?`
	iter = s.session.Query(query, songIDs).Iter()

	var songs []models.Song
	var song models.Song
	for iter.Scan(&song.SongID, &song.Title, &song.UserID, &song.Album, &song.ReleaseDate, &song.Genre, &song.SongURL, &song.ThumbnailURL, &song.PlayCount) {
		songs = append(songs, song)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return songs, nil
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

func (s *scyllaService) LikeSong(userID string, songID gocql.UUID) error {
	query := `INSERT INTO song_likes (user_id, song_id, liked_at) VALUES (?, ?, ?)`
	if err := s.session.Query(query, userID, songID, time.Now()).Exec(); err != nil {
		log.Printf("Failed to like song: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) UnlikeSong(userID string, songID gocql.UUID) error {
	query := `DELETE FROM song_likes WHERE user_id = ? AND song_id = ?`
	if err := s.session.Query(query, userID, songID).Exec(); err != nil {
		log.Printf("Failed to unlike song: %v", err)
		return err
	}
	return nil
}

func (s *scyllaService) GetLikedSongsByUser(userID string) ([]models.Song, error) {
	var songIDs []gocql.UUID
	var songID gocql.UUID
	query := `SELECT song_id FROM song_likes WHERE user_id = ?`
	iter := s.session.Query(query, userID).Iter()
	for iter.Scan(&songID) {
		songIDs = append(songIDs, songID)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	if len(songIDs) == 0 {
		return []models.Song{}, nil
	}

	var likedSongs []models.Song
	var song models.Song
	query = `SELECT song_id, title, user_id, album, release_date, genre, song_url, thumbnail_url, play_count FROM songs WHERE song_id IN ?`
	iter = s.session.Query(query, songIDs).Iter()
	for iter.Scan(&song.SongID, &song.Title, &song.UserID, &song.Album, &song.ReleaseDate, &song.Genre, &song.SongURL, &song.ThumbnailURL, &song.PlayCount) {
		likedSongs = append(likedSongs, song)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return likedSongs, nil
}

func (s *scyllaService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	query := `SELECT user_id, username, email, role FROM users WHERE user_id = ? LIMIT 1`
	if err := s.session.Query(query, userID).Scan(&user.UserID, &user.Username, &user.Email, &user.Role); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *scyllaService) GetSongUserID(songID gocql.UUID) (string, error) {
	var userID string
	query := `SELECT user_id FROM songs WHERE song_id = ? LIMIT 1`
	if err := s.session.Query(query, songID).Scan(&userID); err != nil {
		if err == gocql.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return userID, nil
}

func (s *scyllaService) UpdateUserRole(userID, role string) error {
	query := "UPDATE users SET role = ? WHERE user_id = ?"
	if err := s.session.Query(query, role, userID).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *scyllaService) GetAllArtists() ([]models.Artist, error) {
	query := `SELECT user_id, username, email, role FROM users WHERE role = 'artist' ALLOW FILTERING`
	iter := s.session.Query(query).Iter()

	var artists []models.Artist
	var artist models.Artist
	for iter.Scan(&artist.UserID, &artist.Username, &artist.Email, &artist.Role) {
		// Get followers count for each artist
		followers, err := s.GetArtistFollowersCount(artist.UserID)
		if err != nil {
			return nil, err
		}
		artist.Followers = followers
		artists = append(artists, artist)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (s *scyllaService) GetArtistWithSongs(artistID string) (*models.ArtistWithSongs, error) {
	// Get artist info
	var artist models.Artist
	query := `SELECT user_id, username, email, role FROM users WHERE user_id = ? LIMIT 1`
	if err := s.session.Query(query, artistID).Scan(&artist.UserID, &artist.Username, &artist.Email, &artist.Role); err != nil {
		return nil, err
	}

	// Get followers count
	followers, err := s.GetArtistFollowersCount(artistID)
	if err != nil {
		return nil, err
	}
	artist.Followers = followers

	// Get artist's songs
	songs, err := s.GetSongsByUserID(artistID)
	if err != nil {
		return nil, err
	}

	return &models.ArtistWithSongs{
		Artist: artist,
		Songs:  songs,
	}, nil
}

func (s *scyllaService) FollowArtist(artistID string, followerID string) error {
	query := `INSERT INTO artist_followers (artist_id, follower_id, followed_at) VALUES (?, ?, ?)`
	return s.session.Query(query, artistID, followerID, time.Now()).Exec()
}

func (s *scyllaService) UnfollowArtist(artistID string, followerID string) error {
	query := `DELETE FROM artist_followers WHERE artist_id = ? AND follower_id = ?`
	return s.session.Query(query, artistID, followerID).Exec()
}

func (s *scyllaService) GetFollowedArtists(userID string) ([]models.Artist, error) {
	// First get all artist IDs that the user follows
	query := `SELECT artist_id FROM artist_followers WHERE follower_id = ?`
	iter := s.session.Query(query, userID).Iter()

	var artistIDs []string
	var artistID string
	for iter.Scan(&artistID) {
		artistIDs = append(artistIDs, artistID)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	if len(artistIDs) == 0 {
		return []models.Artist{}, nil
	}

	// Then get artist details for each ID
	var artists []models.Artist
	for _, id := range artistIDs {
		var artist models.Artist
		query := `SELECT user_id, username, email, role FROM users WHERE user_id = ? LIMIT 1`
		if err := s.session.Query(query, id).Scan(&artist.UserID, &artist.Username, &artist.Email, &artist.Role); err != nil {
			continue
		}

		// Get followers count
		followers, err := s.GetArtistFollowersCount(id)
		if err != nil {
			continue
		}
		artist.Followers = followers

		artists = append(artists, artist)
	}

	return artists, nil
}

func (s *scyllaService) GetArtistFollowersCount(artistID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM artist_followers WHERE artist_id = ?`
	if err := s.session.Query(query, artistID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
