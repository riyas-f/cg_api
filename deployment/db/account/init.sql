ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = ON;
ALTER SYSTEM SET logging_collector = ON;
ALTER SYSTEM SET log_directory = 'pg_log';
ALTER SYSTEM SET log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log';
ALTER SYSTEM SET log_connections = ON;
ALTER SYSTEM SET log_disconnections = ON;

SELECT pg_reload_conf();


CREATE DATABASE account_db;

\c account_db

SET TIME ZONE 'UTC';

CREATE TABLE account(
	account_id SERIAL PRIMARY KEY,
	username VARCHAR(64) NOT NULL,
	name VARCHAR(64) NOT NULL,
	password VARCHAR(80) NOT NULL,
    password_salt VARCHAR(64) NOT NULL,
	email VARCHAR(64) NOT NULL,
	steamid VARCHAR(64) NOT NULL,
	is_active VARCHAR(6) NOT NULL
);

CREATE TABLE user_otp(
    otp_id SERIAL PRIMARY KEY,
	email  VARCHAR(64) NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL,
	otp CHAR(6) NOT NULL,
	last_resend TIMESTAMPTZ NOT NULL,
	marked_for_deletion VARCHAR(6) NOT NULL
);