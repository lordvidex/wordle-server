-- for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE game (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    word_id UUID REFERENCES word(id)
);

CREATE TABLE game_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES game(id),
    word_length SMALLINT DEFAULT 5,
    trials SMALLINT DEFAULT 6,
    player_count SMALLINT DEFAULT 2,
    has_analytics BOOLEAN DEFAULT TRUE,
    should_record_time BOOLEAN DEFAULT TRUE,
    can_view_opponents_sessions BOOLEAN DEFAULT TRUE
);

CREATE TABLE game_player (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE game_session (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES game(id) NOT NULL,
    player_id UUID REFERENCES game_player(id) NOT NULL
);

CREATE TABLE game_session_guess(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_session_id UUID REFERENCES game_session(id),
    word_id UUID REFERENCES word(id)
);