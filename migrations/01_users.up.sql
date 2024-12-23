CREATE TYPE role_enum AS ENUM ('ADMIN', 'CONTENT_MANAGER', 'USER');

CREATE TABLE images
(
    id           BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    hash         TEXT        NOT NULL,
    content_type TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file         BYTEA       NOT NULL
);

CREATE TABLE users
(
    id         BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    role       role_enum   NOT NULL DEFAULT 'USER',
    username   TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    image_id   BIGINT REFERENCES images (id)
);
