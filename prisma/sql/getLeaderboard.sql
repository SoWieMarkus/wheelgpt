 SELECT 
    "userId",
    "channelId",
    "displayName", 
    "points",
    "perfectGuessCount",
    "updatedAt",
     CAST(ROW_NUMBER() OVER (ORDER BY "points" DESC, "updatedAt" ASC) AS TEXT) as position
FROM "GuesserLeaderboard"
WHERE "channelId" = :channelId
ORDER BY "points" DESC, "updatedAt" ASC
LIMIT :leaderboardPageSize
OFFSET :leaderboardOffset