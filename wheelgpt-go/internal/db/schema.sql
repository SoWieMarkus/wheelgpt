CREATE TABLE channel (
    id                    TEXT    NOT NULL PRIMARY KEY,
    login                 TEXT    NOT NULL UNIQUE,
    token                 TEXT    NOT NULL,
    display_name          TEXT    NOT NULL,
    profile_image         TEXT    NOT NULL,
    is_live               INTEGER NOT NULL DEFAULT 0,
    created_at            TEXT    NOT NULL DEFAULT (datetime('now')),
    guess_delay_time      REAL    NOT NULL DEFAULT 2,
    bot_active_when_offline INTEGER NOT NULL DEFAULT 0,
    usage_public          INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE trackmania_map (
    channel_id    TEXT    NOT NULL PRIMARY KEY,
    uid           TEXT    NOT NULL,
    name          TEXT    NOT NULL,
    author        TEXT    NOT NULL,
    author_time   INTEGER NOT NULL,
    gold_time     INTEGER NOT NULL,
    silver_time   INTEGER NOT NULL,
    bronze_time   INTEGER NOT NULL,
    champion_time INTEGER NOT NULL,
    tmx_id        INTEGER,
    world_record  INTEGER
);

CREATE TABLE guess (
    channel_id   TEXT    NOT NULL,
    user_id      TEXT    NOT NULL,
    display_name TEXT    NOT NULL,
    time         INTEGER NOT NULL,
    PRIMARY KEY (channel_id, user_id)
);

CREATE TABLE trackmania_room (
    channel_id        TEXT    NOT NULL PRIMARY KEY,
    login             TEXT    NOT NULL,
    name              TEXT    NOT NULL,
    number_of_players INTEGER NOT NULL,
    max_players       INTEGER NOT NULL
);

CREATE TABLE guesser_leaderboard (
    channel_id          TEXT    NOT NULL,
    user_id             TEXT    NOT NULL,
    display_name        TEXT    NOT NULL,
    points              INTEGER NOT NULL,
    perfect_guess_count INTEGER NOT NULL DEFAULT 0,
    updated_at          TEXT    NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (channel_id, user_id)
);
