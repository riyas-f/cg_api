ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = ON;
ALTER SYSTEM SET logging_collector = ON;
ALTER SYSTEM SET log_directory = 'pg_log';
ALTER SYSTEM SET log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log';
ALTER SYSTEM SET log_connections = ON;
ALTER SYSTEM SET log_disconnections = ON;

SELECT pg_reload_conf();

CREATE DATABASE session_db;

\c session_db;

SET TIME ZONE 'UTC';

CREATE TYPE status_enum AS ENUM ('Provisioning', 'WaitingForConnection', 'Pairing', 'Running', 'Failed', 'Terminated');

CREATE TABLE gpu_list {
    gpu_id SERIAL PRIMARY KEY,
    gpu_name VARCHAR(128) NOT NULL,
    gpu_alt_name VARCHAR(128) NOT NULL,
    n_available INTEGER NOT NULL,
    version INTEGER NOT NULL, 
}

CREATE TABLE user_session (
    session_id bytea PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    request_status status_enum NOT NULL,
    last_update TIMESTAMPTZ NOT NULL,
    marked_for_deletion VARCHAR(5) NOT NULL
);

CREATE TABLE session_host(
    host_id SERIAL PRIMARY KEY, 
    webhook_host VARCHAR(15) NOT NULL,
    webhook_port INT NOT NULL,
    network_id VARCHAR(64) NOT NULL,
    session_id bytea NOT NULL,
    FOREIGN KEY (session_id) REFERENCES user_session(session_id)
);

CREATE TABLE session_metadata(
    metadata_id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL, 
    game_id INT NOT NULL,
    game_location_protocol VARCHAR(64) NOT NULL,
    game_location_server_host VARCHAR(64) NOT NULL,
    game_location_path VARCHAR(64) NOT NULL,
    session_id bytea NOT NULL,
    FOREIGN KEY (session_id) REFERENCES user_session(session_id)
);

