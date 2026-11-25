-- name: GetMap :one
SELECT * FROM "TrackmaniaMap" WHERE "channelId" = sqlc.arg(channel_id);

-- name: UpsertMap :exec
INSERT INTO "TrackmaniaMap" ("channelId", "uid", "name", "author", "authorTime", "goldTime", "silverTime", "bronzeTime", "championTime")
VALUES (sqlc.arg(channel_id), sqlc.arg(uid), sqlc.arg(name), sqlc.arg(author), sqlc.arg(author_time), sqlc.arg(gold_time), sqlc.arg(silver_time), sqlc.arg(bronze_time), sqlc.arg(champion_time))
ON CONFLICT ("channelId") DO UPDATE SET
    "uid" = excluded."uid",
    "name" = excluded."name",
    "author" = excluded."author",
    "authorTime" = excluded."authorTime",
    "goldTime" = excluded."goldTime",
    "silverTime" = excluded."silverTime",
    "bronzeTime" = excluded."bronzeTime",
    "championTime" = excluded."championTime",
    "tmxId" = NULL,
    "worldRecord" = NULL;

-- name: SetMapTmxID :exec
UPDATE "TrackmaniaMap" SET "tmxId" = sqlc.arg(tmx_id)
WHERE "channelId" = sqlc.arg(channel_id) AND "uid" = sqlc.arg(uid);

-- name: DeleteMap :exec
DELETE FROM "TrackmaniaMap" WHERE "channelId" = sqlc.arg(channel_id);

-- name: GetRoom :one
SELECT * FROM "TrackmaniaRoom" WHERE "channelId" = sqlc.arg(channel_id);

-- name: UpsertRoom :exec
INSERT INTO "TrackmaniaRoom" ("channelId", "login", "name", "numberOfPlayers", "maxPlayers")
VALUES (sqlc.arg(channel_id), sqlc.arg(login), sqlc.arg(name), sqlc.arg(number_of_players), sqlc.arg(max_players))
ON CONFLICT ("channelId") DO UPDATE SET
    "login" = excluded."login",
    "name" = excluded."name",
    "numberOfPlayers" = excluded."numberOfPlayers",
    "maxPlayers" = excluded."maxPlayers";

-- name: DeleteRoom :exec
DELETE FROM "TrackmaniaRoom" WHERE "channelId" = sqlc.arg(channel_id);
