CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discord_id VARCHAR(100) NOT NULL UNIQUE,
    discord_username VARCHAR(100) NOT NULL,
    discord_avatar_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(100) NOT NULL,
    discord_server_id varchar(100) NOT NULL UNIQUE,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    FOREIGN KEY (creator_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS list_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    list_id UUID NOT NULL,
    discord_nickname VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (list_id) REFERENCES lists(id),
    UNIQUE  (user_id, list_id)
);

CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(15) NOT NULL,
    seed INTEGER NOT NULL,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    FOREIGN KEY (list_id) REFERENCES lists(id),
    FOREIGN KEY (creator_id) REFERENCES list_users(id)
);

CREATE TABLE IF NOT EXISTS ratings (
    movie_id UUID NOT NULL,
    list_user_id UUID NOT NULL,
    rating SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    PRIMARY KEY (movie_id, list_user_id)
);

CREATE TABLE IF NOT EXISTS votes (
    movie_id UUID NOT NULL,
    list_user_id UUID NOT NULL,
    is_approval BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    PRIMARY KEY (movie_id, list_user_id)
);