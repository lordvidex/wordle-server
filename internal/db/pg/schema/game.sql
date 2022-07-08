CREATE TABLE game (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    word_id UUID REFERENCES word(id)
);

CREATE TABLE player(
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE game_session (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES game(id),
    player_id UUID REFERENCES player(id)
);

CREATE TABLE game_session_guess(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_session_id UUID REFERENCES game_session(id),
    word_id UUID REFERENCES word(id)
);