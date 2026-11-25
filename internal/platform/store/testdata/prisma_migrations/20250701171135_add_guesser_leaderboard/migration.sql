-- CreateTable
CREATE TABLE "GuesserLeaderboard" (
    "channelId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "displayName" TEXT NOT NULL,
    "points" INTEGER NOT NULL,

    PRIMARY KEY ("channelId", "userId")
);
