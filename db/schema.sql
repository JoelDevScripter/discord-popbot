-- SQL schema goes here 
-- Inicialmente vacío, solo estructura mínima para usuarios y logs
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    discord_id TEXT UNIQUE NOT NULL,
    coins INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
