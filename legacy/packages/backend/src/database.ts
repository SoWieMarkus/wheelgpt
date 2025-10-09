import { PrismaClient } from "@prisma/client";

export const database = new PrismaClient();

process.on("SIGINT", async () => {
	await database.$disconnect();
	process.exit(0);
});
