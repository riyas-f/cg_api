ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = ON;
ALTER SYSTEM SET logging_collector = ON;
ALTER SYSTEM SET log_directory = 'pg_log';
ALTER SYSTEM SET log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log';
ALTER SYSTEM SET log_connections = ON;
ALTER SYSTEM SET log_disconnections = ON;

SELECT pg_reload_conf();

CREATE DATABASE auth_db;

\c auth_db

SET TIME ZONE 'UTC';

-- public : all user can access the endpoint 
-- restricted : only the owner can access the endpoint
-- private : only admin user that can access the endpoint
CREATE TYPE access_mode AS ENUM ('public', 'restricted', 'private');

CREATE TABLE revoked_token(
	token_id SERIAL PRIMARY KEY,
	username VARCHAR(64) NOT NULL,
	token VARCHAR(1024) NOT NULL,
	type CHAR(16) NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE endpoint_access(
	endpoint_id SERIAL PRIMARY KEY,
	endpoint VARCHAR(256) NOT NULL,
	method VARCHAR(32) NOT NULL,
	access access_mode NOT NULL
);

