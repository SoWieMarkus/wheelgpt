import * as z from "zod";

export const WebTokenSchema = z.object({
	webToken: z.string(),
});

export const PluginTokenSchema = z.object({
	pluginToken: z.string(),
});
