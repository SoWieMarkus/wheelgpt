-- CreateTable
CREATE TABLE "TrackmaniaRoom" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "numberOfPlayers" INTEGER NOT NULL,
    "maxPlayers" INTEGER NOT NULL
);
