-- CreateTable
CREATE TABLE "TrackmaniaMap" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "uid" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "authorTime" INTEGER NOT NULL,
    "goldTime" INTEGER NOT NULL,
    "silverTime" INTEGER NOT NULL,
    "bronzeTime" INTEGER NOT NULL,
    "championTime" INTEGER NOT NULL,
    "tmxId" INTEGER,
    "worldRecord" INTEGER
);

-- CreateTable
CREATE TABLE "Guess" (
    "channelId" TEXT NOT NULL,
    "time" INTEGER NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,

    PRIMARY KEY ("channelId", "userId")
);

-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Channel" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "token" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "profileImage" TEXT NOT NULL,
    "isLive" BOOLEAN NOT NULL DEFAULT false,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "guessDelayTime" INTEGER NOT NULL DEFAULT 2,
    "botActiveWhenOffline" BOOLEAN NOT NULL DEFAULT false,
    "usagePublic" BOOLEAN NOT NULL DEFAULT true
);
INSERT INTO "new_Channel" ("channelId", "createdAt", "displayName", "guessDelayTime", "profileImage", "token") SELECT "channelId", "createdAt", "displayName", "guessDelayTime", "profileImage", "token" FROM "Channel";
DROP TABLE "Channel";
ALTER TABLE "new_Channel" RENAME TO "Channel";
CREATE UNIQUE INDEX "Channel_channelId_key" ON "Channel"("channelId");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
