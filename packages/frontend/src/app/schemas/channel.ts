import { z } from "zod";

export const ChannelSchema = z.object({
	id: z.string(),
	displayName: z.string(),
	profileImage: z.string().url(),
	guessDelayTime: z.number().int(),
	botActiveWhenOffline: z.boolean(),
	usagePublic: z.boolean(),
});

export type Channel = z.infer<typeof ChannelSchema>;
