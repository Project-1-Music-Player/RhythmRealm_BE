CREATE KEYSPACE rhythm_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'} AND durable_writes = true;

USE rhythm_keyspace;

CREATE TABLE user_sessions (
    session_id UUID PRIMARY KEY,
    user_id UUID,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    data TEXT
);