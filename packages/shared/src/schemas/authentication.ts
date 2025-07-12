import * as z from "zod";

export const twitch = z.object({ code: z.string() });
