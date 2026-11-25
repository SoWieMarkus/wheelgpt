/*
  Warnings:

  - Added the required column `displayName` to the `Channel` table without a default value. This is not possible if the table is not empty.
  - Added the required column `profileImage` to the `Channel` table without a default value. This is not possible if the table is not empty.

*/
-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Channel" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "token" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "profileImage" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "guessDelayTime" INTEGER NOT NULL DEFAULT 2
);
INSERT INTO "new_Channel" ("channelId", "createdAt", "guessDelayTime", "token") SELECT "channelId", "createdAt", "guessDelayTime", "token" FROM "Channel";
DROP TABLE "Channel";
ALTER TABLE "new_Channel" RENAME TO "Channel";
CREATE UNIQUE INDEX "Channel_channelId_key" ON "Channel"("channelId");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
