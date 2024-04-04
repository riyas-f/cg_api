ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = ON;
ALTER SYSTEM SET logging_collector = ON;
ALTER SYSTEM SET log_directory = 'pg_log';
ALTER SYSTEM SET log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log';
ALTER SYSTEM SET log_connections = ON;
ALTER SYSTEM SET log_disconnections = ON;

SELECT pg_reload_conf();

CREATE DATABASE session_db;

\c session_db

SET TIME ZONE 'UTC';

CREATE TABLE supported_games (
	games_id SERIAL PRIMARY KEY,
	icon VARCHAR(1024) NOT NULL,
	name VARCHAR(256) NOT NULL,
)

CREATE TABLE user_games(
	username VARCHAR(64) NOT NULL,
	games_id VARCHAR(32) NOT NULL,
	playtime INT NOT NULL,
	last_played TIMESTAMPTZ NOT NULL,
);

