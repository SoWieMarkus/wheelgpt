-- CreateTable
CREATE TABLE "Channel" (
    "channelId" TEXT NOT NULL PRIMARY KEY,
    "token" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateIndex
CREATE UNIQUE INDEX "Channel_channelId_key" ON "Channel"("channelId");
