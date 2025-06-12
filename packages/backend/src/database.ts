import { PrismaClient } from "../generated/prisma";

export const database = new PrismaClient();

process.on("SIGINT", async () => {
	await database.$disconnect();
	process.exit(0);
});
