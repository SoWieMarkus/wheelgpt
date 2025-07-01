 SELECT 
    "userId",
    "channelId",
    "displayName", 
    "points",
    "perfectGuessCount",
     CAST(ROW_NUMBER() OVER (ORDER BY "points" DESC) AS TEXT) as position
FROM "GuesserLeaderboard"
WHERE "channelId" = :channelId
ORDER BY "points" DESC
LIMIT :leaderboardPageSize
OFFSET :leaderboardOffset