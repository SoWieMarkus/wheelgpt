import { z } from "zod";

export const PublicChannelsSchema = z
	.object({
		displayName: z.string(),
		profileImage: z.string().url(),
		isLive: z.boolean(),
	})
	.array();

export type PublicChannels = z.infer<typeof PublicChannelsSchema>;
