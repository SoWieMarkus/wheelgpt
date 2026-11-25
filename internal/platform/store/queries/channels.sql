-- name: GetChannel :one
SELECT c.*,
       CAST(COALESCE(s."botEnabled", 1) AS BOOLEAN) AS "botEnabled",
       CAST(COALESCE(s."scopes", '') AS TEXT) AS "scopes",
       s."authorizedAt" AS "authorizedAt",
       CAST(COALESCE(s."chatSubscription", 'unknown') AS TEXT) AS "chatSubscription"
FROM "Channel" c
LEFT JOIN "ChannelState" s ON s."channelId" = c."id"
WHERE c."id" = sqlc.arg(id);

-- name: GetChannelByLogin :one
SELECT c.*,
       CAST(COALESCE(s."botEnabled", 1) AS BOOLEAN) AS "botEnabled",
       CAST(COALESCE(s."scopes", '') AS TEXT) AS "scopes",
       s."authorizedAt" AS "authorizedAt",
       CAST(COALESCE(s."chatSubscription", 'unknown') AS TEXT) AS "chatSubscription"
FROM "Channel" c
LEFT JOIN "ChannelState" s ON s."channelId" = c."id"
WHERE c."login" = lower(sqlc.arg(login));

-- name: ListChannels :many
SELECT c."id", c."login", c."isLive",
       CAST(COALESCE(s."botEnabled", 1) AS BOOLEAN) AS "botEnabled"
FROM "Channel" c
LEFT JOIN "ChannelState" s ON s."channelId" = c."id"
ORDER BY c."id";

-- name: ListPublicChannels :many
SELECT "id", "login", "displayName", "profileImage", "isLive"
FROM "Channel"
WHERE "usagePublic" = 1
  AND (CAST(sqlc.arg(live_only) AS INTEGER) = 0 OR "isLive" = 1)
  AND (CAST(sqlc.arg(query) AS TEXT) = '' OR "login" LIKE '%' || sqlc.arg(query) || '%' OR "displayName" LIKE '%' || sqlc.arg(query) || '%')
ORDER BY "isLive" DESC, lower("displayName"), "id"
LIMIT sqlc.arg(lim) OFFSET sqlc.arg(off);

-- name: CountPublicChannels :one
SELECT COUNT(*) FROM "Channel"
WHERE "usagePublic" = 1
  AND (CAST(sqlc.arg(live_only) AS INTEGER) = 0 OR "isLive" = 1)
  AND (CAST(sqlc.arg(query) AS TEXT) = '' OR "login" LIKE '%' || sqlc.arg(query) || '%' OR "displayName" LIKE '%' || sqlc.arg(query) || '%');

-- name: UpsertChannel :exec
INSERT INTO "Channel" ("id", "login", "token", "displayName", "profileImage", "createdAt")
VALUES (sqlc.arg(id), sqlc.arg(login), sqlc.arg(token), sqlc.arg(display_name), sqlc.arg(profile_image), sqlc.arg(created_at))
ON CONFLICT ("id") DO UPDATE SET
    "login" = excluded."login",
    "displayName" = excluded."displayName",
    "profileImage" = excluded."profileImage";

-- name: UpdateChannelProfile :exec
UPDATE "Channel"
SET "login" = sqlc.arg(login), "displayName" = sqlc.arg(display_name), "profileImage" = sqlc.arg(profile_image)
WHERE "id" = sqlc.arg(id);

-- name: SetChannelLive :exec
UPDATE "Channel" SET "isLive" = sqlc.arg(is_live) WHERE "id" = sqlc.arg(id);

-- name: SetAllChannelsOffline :exec
UPDATE "Channel" SET "isLive" = 0;

-- name: UpdateChannelSettings :exec
UPDATE "Channel"
SET "guessDelayTime" = sqlc.arg(guess_delay_time),
    "guessMinRequiredAgeTime" = sqlc.arg(guess_min_required_age_time),
    "usagePublic" = sqlc.arg(usage_public)
WHERE "id" = sqlc.arg(id);

-- name: SetChannelLegacyToken :exec
UPDATE "Channel" SET "token" = sqlc.arg(token) WHERE "id" = sqlc.arg(id);

-- name: DeleteChannel :exec
DELETE FROM "Channel" WHERE "id" = sqlc.arg(id);

-- name: CountChannels :one
SELECT COUNT(*) FROM "Channel";

-- name: UpsertChannelAuthorization :exec
INSERT INTO "ChannelState" ("channelId", "botEnabled", "scopes", "authorizedAt")
VALUES (sqlc.arg(channel_id), 1, sqlc.arg(scopes), sqlc.arg(authorized_at))
ON CONFLICT ("channelId") DO UPDATE SET
    "scopes" = excluded."scopes",
    "authorizedAt" = excluded."authorizedAt";

-- name: SetBotEnabled :exec
INSERT INTO "ChannelState" ("channelId", "botEnabled")
VALUES (sqlc.arg(channel_id), sqlc.arg(bot_enabled))
ON CONFLICT ("channelId") DO UPDATE SET "botEnabled" = excluded."botEnabled";

-- name: SetChatSubscription :exec
INSERT INTO "ChannelState" ("channelId", "chatSubscription")
VALUES (sqlc.arg(channel_id), sqlc.arg(status))
ON CONFLICT ("channelId") DO UPDATE SET "chatSubscription" = excluded."chatSubscription";

-- name: DeleteChannelState :exec
DELETE FROM "ChannelState" WHERE "channelId" = sqlc.arg(channel_id);
