SELECT 'CREATE DATABASE challenge_chat'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'challenge_chat')\gexec

\c challenge_chat;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    user_name VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_log(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    message VARCHAR,
    created_at TIMESTAMP NOT NULL
);