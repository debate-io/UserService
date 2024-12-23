CREATE TYPE approving_status_enum AS ENUM ('PENDING', 'APPROVED', 'DECLINED');

CREATE TABLE IF NOT EXISTS metatopics
(
    id           BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name         TEXT                  NOT NULL UNIQUE,
    status       approving_status_enum NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS topics
(
    id           BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name         TEXT                   NOT NULL UNIQUE,
    status       approving_status_enum  NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS metatopics_topics
(
    metatopics_id           BIGINT REFERENCES metatopics(id) ON DELETE CASCADE,
    topics_id               BIGINT REFERENCES topics(id) ON DELETE CASCADE,
    CONSTRAINT unique_metatopics_topics UNIQUE (metatopics_id, topics_id)
);

CREATE TABLE IF NOT EXISTS users_metatopics
(
    user_id                 BIGINT REFERENCES users(id) ON DELETE CASCADE,
    metatopics_id           BIGINT REFERENCES metatopics(id) ON DELETE CASCADE,
    CONSTRAINT unique_users_metatopics UNIQUE (user_id, metatopics_id)
);
