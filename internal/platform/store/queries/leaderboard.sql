-- All-time leaderboard: the table the Prisma backend maintained, so points
-- from before the switchover are kept.

-- name: AddAllTimePoints :exec
INSERT INTO "GuesserLeaderboard" ("channelId", "userId", "displayName", "points", "perfectGuessCount", "updatedAt")
VALUES (sqlc.arg(channel_id), sqlc.arg(user_id), sqlc.arg(display_name), sqlc.arg(points), sqlc.arg(perfect), sqlc.arg(updated_at))
ON CONFLICT ("channelId", "userId") DO UPDATE SET
    "displayName" = excluded."displayName",
    "points" = "GuesserLeaderboard"."points" + excluded."points",
    "perfectGuessCount" = "GuesserLeaderboard"."perfectGuessCount" + excluded."perfectGuessCount",
    "updatedAt" = excluded."updatedAt";

-- name: ListAllTimeLeaderboard :many
WITH ranked AS (
    SELECT "userId", "displayName", "points", "perfectGuessCount",
           ROW_NUMBER() OVER (ORDER BY "points" DESC, "perfectGuessCount" DESC, "updatedAt" ASC, "userId") AS position
    FROM "GuesserLeaderboard"
    WHERE "channelId" = sqlc.arg(channel_id)
)
SELECT CAST(position AS INTEGER) AS position, "userId", "displayName", "points", "perfectGuessCount"
FROM ranked
WHERE CAST(sqlc.arg(query) AS TEXT) = '' OR "displayName" LIKE '%' || sqlc.arg(query) || '%'
ORDER BY position
LIMIT sqlc.arg(lim) OFFSET sqlc.arg(off);

-- name: CountAllTimeLeaderboard :one
SELECT COUNT(*) FROM "GuesserLeaderboard"
WHERE "channelId" = sqlc.arg(channel_id)
  AND (CAST(sqlc.arg(query) AS TEXT) = '' OR "displayName" LIKE '%' || sqlc.arg(query) || '%');

-- name: FindAllTimeEntry :one
-- Looks a guesser up by user id or, if that is empty, by display name.
WITH ranked AS (
    SELECT "userId", "displayName", "points", "perfectGuessCount",
           ROW_NUMBER() OVER (ORDER BY "points" DESC, "perfectGuessCount" DESC, "updatedAt" ASC, "userId") AS position
    FROM "GuesserLeaderboard"
    WHERE "channelId" = sqlc.arg(channel_id)
)
SELECT CAST(position AS INTEGER) AS position, "userId", "displayName", "points", "perfectGuessCount"
FROM ranked
WHERE (CAST(sqlc.arg(user_id) AS TEXT) != '' AND "userId" = sqlc.arg(user_id))
   OR (CAST(sqlc.arg(user_id) AS TEXT) = '' AND lower("displayName") = lower(sqlc.arg(display_name)))
LIMIT 1;

-- Monthly leaderboards: sums over GuessResult within [from, to) in unix ms.
-- The display name is taken from the user's most recent result (SQLite bare
-- column semantics with a single max() aggregate).

-- name: ListPeriodLeaderboard :many
WITH totals AS (
    SELECT "userId", "displayName", MAX("id") AS last_id,
           CAST(SUM("points") AS INTEGER) AS points,
           CAST(SUM("perfect") AS INTEGER) AS perfect_count,
           CAST(MAX(CASE WHEN "points" > 0 THEN "createdAt" END) AS INTEGER) AS last_scored_at
    FROM "GuessResult"
    WHERE "channelId" = sqlc.arg(channel_id) AND "createdAt" >= sqlc.arg(from_ms) AND "createdAt" < sqlc.arg(to_ms)
    GROUP BY "userId"
    HAVING SUM("points") > 0
), ranked AS (
    SELECT "userId", "displayName", points, perfect_count,
           ROW_NUMBER() OVER (ORDER BY points DESC, perfect_count DESC, last_scored_at ASC, "userId") AS position
    FROM totals
)
SELECT CAST(position AS INTEGER) AS position, "userId", "displayName", points, perfect_count
FROM ranked
WHERE CAST(sqlc.arg(query) AS TEXT) = '' OR "displayName" LIKE '%' || sqlc.arg(query) || '%'
ORDER BY position
LIMIT sqlc.arg(lim) OFFSET sqlc.arg(off);

-- name: CountPeriodLeaderboard :one
SELECT COUNT(*) FROM (
    SELECT "userId", "displayName", MAX("id")
    FROM "GuessResult"
    WHERE "channelId" = sqlc.arg(channel_id) AND "createdAt" >= sqlc.arg(from_ms) AND "createdAt" < sqlc.arg(to_ms)
    GROUP BY "userId"
    HAVING SUM("points") > 0
) t
WHERE CAST(sqlc.arg(query) AS TEXT) = '' OR t."displayName" LIKE '%' || sqlc.arg(query) || '%';

-- name: FindPeriodEntry :one
WITH totals AS (
    SELECT "userId", "displayName", MAX("id") AS last_id,
           CAST(SUM("points") AS INTEGER) AS points,
           CAST(SUM("perfect") AS INTEGER) AS perfect_count,
           CAST(MAX(CASE WHEN "points" > 0 THEN "createdAt" END) AS INTEGER) AS last_scored_at
    FROM "GuessResult"
    WHERE "channelId" = sqlc.arg(channel_id) AND "createdAt" >= sqlc.arg(from_ms) AND "createdAt" < sqlc.arg(to_ms)
    GROUP BY "userId"
    HAVING SUM("points") > 0
), ranked AS (
    SELECT "userId", "displayName", points, perfect_count,
           ROW_NUMBER() OVER (ORDER BY points DESC, perfect_count DESC, last_scored_at ASC, "userId") AS position
    FROM totals
)
SELECT CAST(position AS INTEGER) AS position, "userId", "displayName", points, perfect_count
FROM ranked
WHERE (CAST(sqlc.arg(user_id) AS TEXT) != '' AND "userId" = sqlc.arg(user_id))
   OR (CAST(sqlc.arg(user_id) AS TEXT) = '' AND lower("displayName") = lower(sqlc.arg(display_name)))
LIMIT 1;

-- name: ListLeaderboardMonths :many
SELECT DISTINCT CAST(strftime('%Y-%m', "createdAt" / 1000, 'unixepoch') AS TEXT) AS month
FROM "GuessResult"
WHERE "channelId" = sqlc.arg(channel_id) AND "points" > 0
ORDER BY month DESC;
