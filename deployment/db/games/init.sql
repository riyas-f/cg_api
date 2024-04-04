ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = ON;
ALTER SYSTEM SET logging_collector = ON;
ALTER SYSTEM SET log_directory = 'pg_log';
ALTER SYSTEM SET log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log';
ALTER SYSTEM SET log_connections = ON;
ALTER SYSTEM SET log_disconnections = ON;

SELECT pg_reload_conf();

CREATE DATABASE games_db;

\c games_db

SET TIME ZONE 'UTC';

CREATE TYPE protocol AS ENUM ('nas', 'http', 'gsutil');

CREATE TABLE storage_location(
	storage_id SERIAL PRIMARY KEY,
	protocol protocol NOT NULL,
	host VARCHAR(1024) NOT NULL,
	port INT, 
	location VARCHAR(1024) NOT NULL
);

CREATE TABLE games (
	game_id INT PRIMARY KEY,
	name VARCHAR(256) NOT NULL,
	display_picture_url VARCHAR(1024) NOT NULL,
	icon_url VARCHAR(1024) NOT NULL,
	storage_id INT NOT NULL, 
	FOREIGN KEY (storage_id) REFERENCES storage_location(storage_id)
);

CREATE TABLE user_games(
	collections_id SERIAL PRIMARY KEY, 
	game_id INT NOT NULL,
	FOREIGN KEY (game_id) REFERENCES games(game_id),
	username VARCHAR(64) NOT NULL
);

INSERT INTO storage_location(protocol,host,port,location) 
VALUES 
	('nas', '10.1.11.169', NULL, 'data/games/Granblue Fantasy: Relink'),
	('nas', '10.1.11.169', NULL,  'data/games/Red Dead Redemption 2'),
	('nas', '10.1.11.169', NULL, 'data/games/Persona 3 Reload');
INSERT INTO games(game_id,name,display_picture_url,icon_url,storage_id)
VALUES 
	(881020, 'Granblue Fantasy: Relink', 'https://cdn.akamai.steamstatic.com/steam/apps/881020/header.jpg', 'http://media.steampowered.com/steamcommunity/public/images/apps/881020/1b9bcf2c07e1c760349fd4ecbc3b8690ef915e11.jpg', 1),
	(1174180, 'Red Dead Redemption 2', 'https://cdn.akamai.steamstatic.com/steam/apps/1174180/header.jpg', 'http://media.steampowered.com/steamcommunity/public/images/apps/1174180/5106abd9c1187a97f23295a0ba9470c94804ec6c.jpg', 2),
	(2161700, 'Persona 3 Reload', 'https://cdn.akamai.steamstatic.com/steam/apps/2161700/header.jpg', 'http://media.steampowered.com/steamcommunity/public/images/apps/2161700/2aeb189b126f766bb5930f725fbdcdd171a93c56.jpg', 3);

