-- +goose Up
-- Baseline: the schema the Prisma backend left behind (prisma/migrations up to
-- 20260716095641). Every statement is guarded with IF NOT EXISTS so it is a
-- no-op on the production database and creates the same tables on a fresh one.
-- Never change these tables in a later migration while the old backend may
-- still run against the same file.

CREATE TABLE IF NOT EXISTS "Channel" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "profileImage" TEXT NOT NULL,
    "isLive" BOOLEAN NOT NULL DEFAULT false,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "guessDelayTime" REAL NOT NULL DEFAULT 2,
    "guessMinRequiredAgeTime" REAL NOT NULL DEFAULT 0,
    "botActiveWhenOffline" BOOLEAN NOT NULL DEFAULT false,
    "usagePublic" BOOLEAN NOT NULL DEFAULT true
);
CREATE UNIQUE INDEX IF NOT EXISTS "Channel_login_key" ON "Channel"("login");

CREATE TABLE IF NOT EXISTS "TrackmaniaMap" (
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

CREATE TABLE IF NOT EXISTS "Guess" (
    "channelId" TEXT NOT NULL,
    "time" INTEGER NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("channelId", "userId")
);

CREATE TABLE IF NOT EXISTS "TrackmaniaRoom" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "numberOfPlayers" INTEGER NOT NULL,
    "maxPlayers" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS "GuesserLeaderboard" (
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "points" INTEGER NOT NULL,
    "perfectGuessCount" INTEGER NOT NULL DEFAULT 0,
    "updatedAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("channelId", "userId")
);
