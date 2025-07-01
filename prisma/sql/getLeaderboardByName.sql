WITH RankedLeaderboard AS (
    SELECT 
        "userId",
        "channelId",
        "displayName",
        "points",
        "perfectGuessCount",
        CAST(ROW_NUMBER() OVER (ORDER BY "points" DESC) AS TEXT) as position
    FROM "GuesserLeaderboard"
    WHERE "channelId" = :channelId
)
SELECT * FROM RankedLeaderboard 
WHERE "displayName" LIKE :displayName COLLATE NOCASE