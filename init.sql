CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    age INTEGER NULL,
    gender TEXT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS llm_interactions (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    prompt JSONB NOT NULL,
    response JSONB,
    latency_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE goals (
    id SERIAL PRIMARY KEY,
    user_id TEXT,
    goal_data JSONB
);
