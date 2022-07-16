-- for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- WORD
CREATE TABLE word (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    time_played TIMESTAMPTZ NOT NULL,
    -- keys are the letters in the word, values are either -1, 0, 1
    letters JSON NOT NULL
);

CREATE TABLE IF NOT EXISTS game (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invite_id VARCHAR(16),
    word_id UUID REFERENCES word(id),
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
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS game_player (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES wordlewf_user(id) NOT NULL,
    game_id UUID REFERENCES game(id) NOT NULL,
    name VARCHAR(255) NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE game_player_word (
    player_id UUID REFERENCES game_player(id),
    word_id UUID REFERENCES word(id) ON DELETE CASCADE,
    CONSTRAINT game_player_word_pkey PRIMARY KEY (player_id, word_id)
);