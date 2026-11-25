-- +goose Up
-- Additive tables for the Go backend. The Prisma tables are not touched, so
-- the old backend can still run against this file (rollback path).
-- All new timestamps are INTEGER unix milliseconds, the format Prisma uses.

-- Per-channel state that did not exist before. A channel without a row is a
-- legacy channel: bot enabled, never authorized through the new OAuth flow.
CREATE TABLE IF NOT EXISTS "ChannelState" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "botEnabled" BOOLEAN NOT NULL DEFAULT true,
    "scopes" TEXT NOT NULL DEFAULT '',
    "authorizedAt" INTEGER,
    -- unknown | enabled | needs_authorization | disabled
    "chatSubscription" TEXT NOT NULL DEFAULT 'unknown'
);

-- Static plugin tokens. Only the SHA-256 of the token is stored.
CREATE TABLE IF NOT EXISTS "PluginToken" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "tokenHash" TEXT NOT NULL,
    "createdAt" INTEGER NOT NULL,
    "lastUsedAt" INTEGER
);
CREATE UNIQUE INDEX IF NOT EXISTS "PluginToken_tokenHash_key" ON "PluginToken"("tokenHash");

-- Website sessions. "id" is the SHA-256 of the cookie value.
CREATE TABLE IF NOT EXISTS "Session" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "channelId" TEXT NOT NULL,
    "createdAt" INTEGER NOT NULL,
    "expiresAt" INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS "Session_channelId_idx" ON "Session"("channelId");

-- One row per evaluated personal best.
CREATE TABLE IF NOT EXISTS "GuessRound" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "channelId" TEXT NOT NULL,
    "resultTime" INTEGER NOT NULL,
    "mapUid" TEXT,
    "mapName" TEXT,
    "createdAt" INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS "GuessRound_channel_created_idx" ON "GuessRound"("channelId", "createdAt");

-- One row per guess that took part in a round. Winners have points > 0.
-- Monthly leaderboards are sums over this table.
CREATE TABLE IF NOT EXISTS "GuessResult" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "roundId" INTEGER NOT NULL REFERENCES "GuessRound"("id") ON DELETE CASCADE,
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "guessTime" INTEGER NOT NULL,
    "resultTime" INTEGER NOT NULL,
    "points" INTEGER NOT NULL,
    "perfect" BOOLEAN NOT NULL,
    "createdAt" INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS "GuessResult_channel_created_idx" ON "GuessResult"("channelId", "createdAt");
CREATE INDEX IF NOT EXISTS "GuessResult_channel_user_idx" ON "GuessResult"("channelId", "userId", "id");
CREATE INDEX IF NOT EXISTS "GuessResult_round_idx" ON "GuessResult"("roundId");
