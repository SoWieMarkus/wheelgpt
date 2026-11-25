-- name: CreateRound :one
INSERT INTO "GuessRound" ("channelId", "resultTime", "mapUid", "mapName", "createdAt")
VALUES (sqlc.arg(channel_id), sqlc.arg(result_time), sqlc.narg(map_uid), sqlc.narg(map_name), sqlc.arg(created_at))
RETURNING *;

-- name: CreateResult :exec
INSERT INTO "GuessResult" ("roundId", "channelId", "userId", "displayName", "guessTime", "resultTime", "points", "perfect", "createdAt")
VALUES (sqlc.arg(round_id), sqlc.arg(channel_id), sqlc.arg(user_id), sqlc.arg(display_name), sqlc.arg(guess_time), sqlc.arg(result_time), sqlc.arg(points), sqlc.arg(perfect), sqlc.arg(created_at));

-- name: DeleteChannelRounds :exec
DELETE FROM "GuessRound" WHERE "channelId" = sqlc.arg(channel_id);

-- name: DeleteChannelResults :exec
DELETE FROM "GuessResult" WHERE "channelId" = sqlc.arg(channel_id);

-- name: DeleteChannelAllTime :exec
DELETE FROM "GuesserLeaderboard" WHERE "channelId" = sqlc.arg(channel_id);
