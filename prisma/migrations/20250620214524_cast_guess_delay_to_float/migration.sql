/*
  Warnings:

  - You are about to alter the column `guessDelayTime` on the `Channel` table. The data in that column could be lost. The data in that column will be cast from `Int` to `Float`.

*/
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
    "botActiveWhenOffline" BOOLEAN NOT NULL DEFAULT false,
    "usagePublic" BOOLEAN NOT NULL DEFAULT true
);
INSERT INTO "new_Channel" ("botActiveWhenOffline", "createdAt", "displayName", "guessDelayTime", "id", "isLive", "login", "profileImage", "token", "usagePublic") SELECT "botActiveWhenOffline", "createdAt", "displayName", "guessDelayTime", "id", "isLive", "login", "profileImage", "token", "usagePublic" FROM "Channel";
DROP TABLE "Channel";
ALTER TABLE "new_Channel" RENAME TO "Channel";
CREATE UNIQUE INDEX "Channel_login_key" ON "Channel"("login");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
