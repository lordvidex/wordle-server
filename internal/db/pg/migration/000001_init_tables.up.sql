-- for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS game (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invite_id VARCHAR(16) NOT NULL,
    word VARCHAR(255),
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS game_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES game(id) ON DELETE CASCADE,
    word_length SMALLINT DEFAULT 5,
    trials SMALLINT DEFAULT 6,
    player_count SMALLINT DEFAULT 2,
    has_analytics BOOLEAN DEFAULT TRUE,
    should_record_time BOOLEAN DEFAULT TRUE,
    can_view_opponents_sessions BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS wordlewf_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    points INTEGER DEFAULT 0,
    email VARCHAR(255) NOT NULL, -- FIXME: change to username?
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS wordlewf_user_games (
    user_id UUID REFERENCES wordlewf_user(id) ON DELETE CASCADE,
    game_id UUID REFERENCES game(id) ON DELETE CASCADE,
    user_game_name VARCHAR(255) NOT NULL,
    points INTEGER DEFAULT 0,
    position INTEGER DEFAULT -1,
    CONSTRAINT wordlewf_user_games_pk PRIMARY KEY (user_id, game_id)
);