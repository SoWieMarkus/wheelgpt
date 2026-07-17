-- name: UpsertGuess :exec
INSERT INTO guess (channel_id, user_id, display_name, time)
VALUES (?, ?, ?, ?)
ON CONFLICT (channel_id, user_id) DO UPDATE SET
    time         = excluded.time,
    display_name = excluded.display_name;

-- name: GetGuess :one
SELECT * FROM guess WHERE channel_id = ? AND user_id = ?;

-- name: GetGuessesByChannel :many
SELECT * FROM guess WHERE channel_id = ?;

-- name: DeleteGuessesByChannel :exec
DELETE FROM guess WHERE channel_id = ?;

-- name: UpsertLeaderboard :exec
INSERT INTO guesser_leaderboard (channel_id, user_id, display_name, points, perfect_guess_count, updated_at)
VALUES (?, ?, ?, ?, ?, datetime('now'))
ON CONFLICT (channel_id, user_id) DO UPDATE SET
    points              = guesser_leaderboard.points + excluded.points,
    perfect_guess_count = guesser_leaderboard.perfect_guess_count + excluded.perfect_guess_count,
    display_name        = excluded.display_name,
    updated_at          = datetime('now');

-- name: GetLeaderboard :many
SELECT
    display_name,
    points,
    perfect_guess_count,
    ROW_NUMBER() OVER (ORDER BY points DESC) AS position
FROM guesser_leaderboard
WHERE channel_id = ?
ORDER BY points DESC
LIMIT ? OFFSET ?;

-- name: GetLeaderboardByName :one
SELECT
    display_name,
    points,
    ROW_NUMBER() OVER (ORDER BY points DESC) AS position
FROM guesser_leaderboard
WHERE channel_id = ?
  AND lower(display_name) = lower(?);

-- name: GetChannel :one
SELECT * FROM channel WHERE id = ?;

-- name: GetChannelByLogin :one
SELECT * FROM channel WHERE login = ?;

-- name: ListChannels :many
SELECT id, login, display_name, guess_delay_time, bot_active_when_offline, usage_public FROM channel;

-- name: UpsertChannel :exec
INSERT INTO channel (id, login, token, display_name, profile_image)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    display_name  = excluded.display_name,
    profile_image = excluded.profile_image,
    login         = excluded.login;

-- name: UpdateChannelSettings :exec
UPDATE channel SET
    guess_delay_time        = ?,
    bot_active_when_offline = ?,
    usage_public            = ?
WHERE id = ?;

-- name: SetChannelLive :exec
UPDATE channel SET is_live = ? WHERE id = ?;

-- name: UpdateChannelDetails :exec
UPDATE channel SET display_name = ?, login = ? WHERE id = ?;

-- name: DeleteChannel :exec
DELETE FROM channel WHERE id = ?;

-- name: RotateChannelToken :exec
UPDATE channel SET token = ? WHERE id = ?;

-- name: GetTrackmaniaMap :one
SELECT * FROM trackmania_map WHERE channel_id = ?;

-- name: UpsertTrackmaniaMap :exec
INSERT INTO trackmania_map (channel_id, uid, name, author, author_time, gold_time, silver_time, bronze_time, champion_time, tmx_id, world_record)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (channel_id) DO UPDATE SET
    uid           = excluded.uid,
    name          = excluded.name,
    author        = excluded.author,
    author_time   = excluded.author_time,
    gold_time     = excluded.gold_time,
    silver_time   = excluded.silver_time,
    bronze_time   = excluded.bronze_time,
    champion_time = excluded.champion_time,
    tmx_id        = excluded.tmx_id,
    world_record  = excluded.world_record;

-- name: DeleteTrackmaniaMap :exec
DELETE FROM trackmania_map WHERE channel_id = ?;

-- name: GetTrackmaniaRoom :one
SELECT * FROM trackmania_room WHERE channel_id = ?;

-- name: UpsertTrackmaniaRoom :exec
INSERT INTO trackmania_room (channel_id, login, name, number_of_players, max_players)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (channel_id) DO UPDATE SET
    login             = excluded.login,
    name              = excluded.name,
    number_of_players = excluded.number_of_players,
    max_players       = excluded.max_players;

-- name: DeleteTrackmaniaRoom :exec
DELETE FROM trackmania_room WHERE channel_id = ?;

-- name: ClearChannelGameState :exec
DELETE FROM guess WHERE channel_id = ?;
