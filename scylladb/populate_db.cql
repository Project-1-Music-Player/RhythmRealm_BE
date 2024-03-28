USE rhythm_keyspace;

INSERT INTO user_sessions (session_id, user_id, start_time, end_time, data)
VALUES (uuid(), uuid(), toTimestamp(now()), toTimestamp(now()), 'Sample data');