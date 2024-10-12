CREATE TABLE recovery_codes
(
    email      TEXT PRIMARY KEY REFERENCES users (email),
    code       VARCHAR(6)  NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL
);
