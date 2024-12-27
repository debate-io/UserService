ALTER TABLE games
    ADD COLUMN room_uid VARCHAR(255) NOT NULL,
    ADD CONSTRAINT games_room_uid_unique UNIQUE (room_uid);
