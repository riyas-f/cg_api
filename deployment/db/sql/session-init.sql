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

CREATE TYPE status_enum AS ENUM ('Provisioning', 'WaitingForConnection', 'Pairing
', 'Running', 'Failed', 'TERMINATED');

CREATE TABLE session_host(
    host_id SERIAL PRIMARY KEY, 
    host_ip VARCHAR(15) NOT NULL,
    network_id VARCHAR(64) NOT NULL, 
)

CREATE TABLE user_session (
    session_id bytea PRIMARY KEY,
    request_status status_enum NOT NULL,
    host_id INT, 
    FOREIGN KEY(host_id) REFERENCES session_host(host_id),
    last_update TIMESTAMPTZ NOT NULL,
    marked_for_deletion BOOLEAN NOT NULL, 
);

CREATE TABLE session_metadata(
    metadata_id SERIAL PRIMARY KEY,
    created_at TIMESTAMPZ NOT NULL, 
    game_id INT NOT NULL,
    username VARCHAR(64) NOT NULL,
    session_id bytea NOT NULL,
    FOREIGN KEY (session_id) REFERENCES user_session(session_id),
)

