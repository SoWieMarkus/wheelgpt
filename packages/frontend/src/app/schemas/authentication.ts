import { z } from "zod";

export const WebTokenSchema = z.object({
	webToken: z.string(),
});
