import { z } from "zod";

export const LeaderboardEntrySchema = z.object({
	userId: z.string(),
	displayName: z.string(),
	channelId: z.string(),
	points: z.number().int(),
	position: z.coerce.number().int(),
});

export type LeaderboardEntry = z.infer<typeof LeaderboardEntrySchema>;
