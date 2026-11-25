/*
  Warnings:

  - The primary key for the `Channel` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `channelId` on the `Channel` table. All the data in the column will be lost.
  - Added the required column `id` to the `Channel` table without a default value. This is not possible if the table is not empty.
  - Added the required column `login` to the `Channel` table without a default value. This is not possible if the table is not empty.

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
    "guessDelayTime" INTEGER NOT NULL DEFAULT 2,
    "botActiveWhenOffline" BOOLEAN NOT NULL DEFAULT false,
    "usagePublic" BOOLEAN NOT NULL DEFAULT true
);
INSERT INTO "new_Channel" ("botActiveWhenOffline", "createdAt", "displayName", "guessDelayTime", "isLive", "profileImage", "token", "usagePublic") SELECT "botActiveWhenOffline", "createdAt", "displayName", "guessDelayTime", "isLive", "profileImage", "token", "usagePublic" FROM "Channel";
DROP TABLE "Channel";
ALTER TABLE "new_Channel" RENAME TO "Channel";
CREATE UNIQUE INDEX "Channel_login_key" ON "Channel"("login");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
