DROP KEYSPACE IF EXISTS rhythm_keyspace;
CREATE KEYSPACE rhythm_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'} AND durable_writes = true;

USE rhythm_keyspace;

CREATE TABLE IF NOT EXISTS artists (
    user_id TEXT PRIMARY KEY,
    name TEXT,
    bio TEXT,
    profile_pic_url TEXT
);

CREATE TABLE IF NOT EXISTS songs (
    song_id UUID PRIMARY KEY,
    title TEXT,
    user_id TEXT,
    album TEXT,
    release_date TIMESTAMP,
    genre TEXT,
    song_url TEXT,
    thumbnail_url TEXT,
    play_count INT
);


CREATE INDEX IF NOT EXISTS songs_genre_idx ON songs(genre);

-- Table for storing user information (listeners and admin)
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT,
    email TEXT,
    role TEXT -- 'listener', 'admin'
);

CREATE TABLE IF NOT EXISTS playlists (
    playlist_id UUID PRIMARY KEY,
    user_id TEXT,
    name TEXT,
    description TEXT
);

CREATE TABLE IF NOT EXISTS playlist_songs (
    playlist_id UUID,
    song_id UUID,
    added_at TIMESTAMP,
    PRIMARY KEY (playlist_id, song_id)
)

CREATE TABLE song_play_counts (
  song_id UUID,
  play_count COUNTER,
  PRIMARY KEY (artist_id, song_id)
);

CREATE TABLE song_likes (
  user_id TEXT,
  song_id UUID,
  liked_at TIMESTAMP,
  PRIMARY KEY (user_id, song_id)
);

CREATE INDEX song_likes_by_song ON song_likes (song_id);
CREATE INDEX IF NOT EXISTS playlists_user_id_idx ON playlists(user_id);

CREATE INDEX IF NOT EXISTS songs_user_id_idx ON songs(user_id);