-- name: UpsertGuess :exec
-- A changed guess counts as a new guess for the minimum-age rule.
INSERT INTO "Guess" ("channelId", "userId", "displayName", "time", "createdAt")
VALUES (sqlc.arg(channel_id), sqlc.arg(user_id), sqlc.arg(display_name), sqlc.arg(time), sqlc.arg(created_at))
ON CONFLICT ("channelId", "userId") DO UPDATE SET
    "displayName" = excluded."displayName",
    "time" = excluded."time",
    "createdAt" = excluded."createdAt";

-- name: GetGuess :one
SELECT * FROM "Guess" WHERE "channelId" = sqlc.arg(channel_id) AND "userId" = sqlc.arg(user_id);

-- name: ListGuesses :many
SELECT * FROM "Guess" WHERE "channelId" = sqlc.arg(channel_id) ORDER BY "time", "userId";

-- name: DeleteGuesses :execrows
DELETE FROM "Guess" WHERE "channelId" = sqlc.arg(channel_id);
