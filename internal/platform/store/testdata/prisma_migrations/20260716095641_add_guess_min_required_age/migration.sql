-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Channel" (
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
INSERT INTO "new_Channel" ("botActiveWhenOffline", "createdAt", "displayName", "guessDelayTime", "id", "isLive", "login", "profileImage", "token", "usagePublic") SELECT "botActiveWhenOffline", "createdAt", "displayName", "guessDelayTime", "id", "isLive", "login", "profileImage", "token", "usagePublic" FROM "Channel";
DROP TABLE "Channel";
ALTER TABLE "new_Channel" RENAME TO "Channel";
CREATE UNIQUE INDEX "Channel_login_key" ON "Channel"("login");
CREATE TABLE "new_Guess" (
    "channelId" TEXT NOT NULL,
    "time" INTEGER NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ("channelId", "userId")
);
INSERT INTO "new_Guess" ("channelId", "displayName", "time", "userId") SELECT "channelId", "displayName", "time", "userId" FROM "Guess";
DROP TABLE "Guess";
ALTER TABLE "new_Guess" RENAME TO "Guess";
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
