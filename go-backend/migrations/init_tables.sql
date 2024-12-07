CREATE TABLE IF NOT EXISTS refresh_sessions (
    id INT GENERATED ALWAYS AS IDENTITY,
    user_id UUID NOT NULL,
    refresh_token UUID NOT NULL DEFAULT gen_random_uuid(),
    issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expiration TIMESTAMPTZ NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    oauth_id VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    avatar VARCHAR(255),
    PRIMARY KEY(id)
)