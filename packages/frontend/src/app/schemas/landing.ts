import { z } from "zod";

export const PublicChannelsSchema = z
	.object({
		login: z.string(),
		displayName: z.string(),
		profileImage: z.string().url(),
		isLive: z.boolean(),
	})
	.array();

export type PublicChannels = z.infer<typeof PublicChannelsSchema>;
