-- ==========================================================================
-- Examples:
-- psql -h 12.13.14.15 -p 5432 -U username
-- psql -h 12.13.14.15 -p 5432 -d DB_NAME -U DB_USER -a -f scripts/create_database.sql -W
-- ==========================================================================

BEGIN TRANSACTION;

DROP TABLE IF EXISTS rs_game_players;
DROP TABLE IF EXISTS rs_session;
DROP TABLE IF EXISTS rs_game;
DROP TABLE IF EXISTS rs_sport;
DROP TABLE IF EXISTS rs_user;
DROP TABLE IF EXISTS rs_audit_trail;

CREATE TABLE IF NOT EXISTS "rs_user" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	"user_id" VARCHAR(256) NOT NULL UNIQUE,
	"rf_id" VARCHAR(128) UNIQUE,
	"password" VARCHAR(256),
	"first_name" VARCHAR(128) NOT NULL,
	"last_name" VARCHAR(128) NOT NULL,
	"email" VARCHAR(256) NOT NULL UNIQUE,
	"about" VARCHAR(512),
	"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sport/Game like Chess, Ping Pong
CREATE TABLE IF NOT EXISTS "rs_sport" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	-- TODO: maybe add category
	"sport_name" VARCHAR(128) NOT NULL UNIQUE
);

-- ==========================================================================
-- Represents one played game that can contain X players. Those players are
-- stored into rs_game_players table
-- ==========================================================================
CREATE TABLE IF NOT EXISTS "rs_game" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	"sport_name" VARCHAR(128) NOT NULL,
	"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- ==========================================================================
-- Represents all players for specific game and those player
-- statuses (Is player a Winner, Loser, Tie) for each game
-- ==========================================================================

CREATE TABLE IF NOT EXISTS "rs_game_players" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	"game_id" INTEGER REFERENCES rs_game(id) NOT NULL,
	-- Player must exist in database
	"user_id" VARCHAR(128) REFERENCES rs_user(user_id) NOT NULL,
	-- We have to code types for this, 0 loser, 1 winner, 2 tie
	"result" INTEGER NOT NULL
);

-- Authentication
CREATE TABLE IF NOT EXISTS "rs_session" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	-- Affects to session related code
	"user_id" VARCHAR(128),
	"token" VARCHAR(128),
	"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	"modified_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	-- Github Authentication
	-- FOREIGN KEY(user_id) REFERENCES rs_user(user_id)
);

-- Config
CREATE TABLE IF NOT EXISTS "rs_config" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	"key" VARCHAR(128) NOT NULL UNIQUE,
	"value_type" VARCHAR(128) NOT NULL,
	-- String, Array of ?
	"type" VARCHAR(128) NOT NULL
);

-- Audit trail
CREATE TABLE IF NOT EXISTS "rs_audit_trail" (
	"id" SERIAL NOT NULL PRIMARY KEY,
	"event" VARCHAR(128) NOT NULL,
	"subject" VARCHAR(128) NOT NULL,
	"object" VARCHAR(128),
	"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMIT;
