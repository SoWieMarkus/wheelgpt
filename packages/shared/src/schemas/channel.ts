import * as z from "zod";

export const settings = z.object({
	botActiveWhenOffline: z.boolean(),
	guessDelayTime: z.number().min(0).max(60),
	usagePublic: z.boolean(),
});
