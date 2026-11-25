-- name: GetPluginTokenByHash :one
SELECT * FROM "PluginToken" WHERE "tokenHash" = sqlc.arg(token_hash);

-- name: GetPluginToken :one
SELECT * FROM "PluginToken" WHERE "channelId" = sqlc.arg(channel_id);

-- name: UpsertPluginToken :exec
INSERT INTO "PluginToken" ("channelId", "tokenHash", "createdAt")
VALUES (sqlc.arg(channel_id), sqlc.arg(token_hash), sqlc.arg(created_at))
ON CONFLICT ("channelId") DO UPDATE SET
    "tokenHash" = excluded."tokenHash",
    "createdAt" = excluded."createdAt",
    "lastUsedAt" = NULL;

-- name: TouchPluginToken :exec
UPDATE "PluginToken" SET "lastUsedAt" = sqlc.arg(last_used_at) WHERE "channelId" = sqlc.arg(channel_id);

-- name: DeletePluginToken :exec
DELETE FROM "PluginToken" WHERE "channelId" = sqlc.arg(channel_id);

-- name: CreateSession :exec
INSERT INTO "Session" ("id", "channelId", "createdAt", "expiresAt")
VALUES (sqlc.arg(id), sqlc.arg(channel_id), sqlc.arg(created_at), sqlc.arg(expires_at));

-- name: GetSession :one
SELECT * FROM "Session" WHERE "id" = sqlc.arg(id) AND "expiresAt" > sqlc.arg(now);

-- name: DeleteSession :exec
DELETE FROM "Session" WHERE "id" = sqlc.arg(id);

-- name: DeleteChannelSessions :exec
DELETE FROM "Session" WHERE "channelId" = sqlc.arg(channel_id);

-- name: DeleteExpiredSessions :exec
DELETE FROM "Session" WHERE "expiresAt" <= sqlc.arg(now);
