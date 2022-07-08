CREATE TABLE word (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    time_played TIMESTAMP NOT NULL,
    -- keys are the letters in the word, values are either -1, 0, 1
    letters JSON NOT NULL
);