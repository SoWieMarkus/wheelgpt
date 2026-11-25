-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_GuesserLeaderboard" (
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "points" INTEGER NOT NULL,
    "perfectGuessCount" INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY ("channelId", "userId")
);
INSERT INTO "new_GuesserLeaderboard" ("channelId", "displayName", "points", "userId") SELECT "channelId", "displayName", "points", "userId" FROM "GuesserLeaderboard";
DROP TABLE "GuesserLeaderboard";
ALTER TABLE "new_GuesserLeaderboard" RENAME TO "GuesserLeaderboard";
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
