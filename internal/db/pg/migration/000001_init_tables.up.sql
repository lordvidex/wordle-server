-- for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS game (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invite_id VARCHAR(16) NOT NULL,
    word VARCHAR(255) NOT NULL,
    player_count SMALLINT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS game_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES game(id) ON DELETE CASCADE,
    word_length SMALLINT DEFAULT 5,
    trials SMALLINT DEFAULT 6,
    max_player_count SMALLINT DEFAULT 2,
    has_analytics BOOLEAN DEFAULT TRUE,
    should_record_time BOOLEAN DEFAULT TRUE,
    can_view_opponents_sessions BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS player (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    points BIGINT DEFAULT 0,
    email VARCHAR(255) NOT NULL, -- FIXME: change to username?
    password VARCHAR(255) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS player_games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id UUID REFERENCES player(id) ON DELETE CASCADE,
    game_id UUID REFERENCES game(id) ON DELETE CASCADE,
    user_game_name VARCHAR(255) NOT NULL,
    points INTEGER DEFAULT 0,
    position INTEGER DEFAULT -1,
    CONSTRAINT player_games_uk UNIQUE(player_id, game_id)
);

CREATE TABLE IF NOT EXISTS player_game_words (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    player_games_id UUID REFERENCES player_games(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    played_at TIMESTAMPTZ NOT NULL
);