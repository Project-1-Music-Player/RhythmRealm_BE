USE rhythm_keyspace;

-- Insert sample users
INSERT INTO users (user_id, username, email, role) VALUES ('00000000-0000-0000-0000-000000000001', 'listener1', 'listener1@example.com', 'listener');
INSERT INTO users (user_id, username, email, role) VALUES ('00000000-0000-0000-0000-000000000002', 'artist1', 'artist1@example.com', 'admin');
INSERT INTO users (user_id, username, email, role) VALUES ('00000000-0000-0000-0000-000000000003', 'artist2', 'artist2@example.com', 'admin');

-- Insert sample artists
INSERT INTO artists (user_id, name, bio, profile_pic_url) VALUES ('00000000-0000-0000-0000-000000000002', 'Artist One', 'Bio for Artist One.', 'http://example.com/profiles/artist1.jpg');
INSERT INTO artists (user_id, name, bio, profile_pic_url) VALUES ('00000000-0000-0000-0000-000000000003', 'Artist Two', 'Bio for Artist Two.', 'http://example.com/profiles/artist2.jpg');

-- Insert sample playlists
INSERT INTO playlists (playlist_id, user_id, name, description) VALUES (33333333-3333-3333-3333-333333333333, '00000000-0000-0000-0000-000000000001', 'My Playlist', 'A great selection of songs.');

-- Insert songs into the playlist using the playlist_id and song_id
INSERT INTO playlist_songs (playlist_id, added_at, song_id) VALUES (33333333-3333-3333-3333-333333333333, toUnixTimestamp(now()), 11111111-1111-1111-1111-111111111111);
INSERT INTO playlist_songs (playlist_id, added_at, song_id) VALUES (33333333-3333-3333-3333-333333333333, toUnixTimestamp(now()), 22222222-2222-2222-2222-222222222222);
