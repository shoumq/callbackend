CREATE TABLE IF NOT EXISTS requests
(
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL,
    phone      TEXT,
    text       TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);