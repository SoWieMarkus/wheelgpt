import * as z from "zod";

export const ChannelSchema = z.object({
	id: z.string(),
	displayName: z.string(),
	profileImage: z.string().url(),
	guessDelayTime: z.number().int(),
	botActiveWhenOffline: z.boolean(),
	usagePublic: z.boolean(),
});

export type Channel = z.infer<typeof ChannelSchema>;

export const PublicChannelDetailsSchema = z.object({
	id: z.string(),
	login: z.string(),
	displayName: z.string(),
	profileImage: z.string().url(),
	isLive: z.boolean(),
});

export type PublicChannelDetails = z.infer<typeof PublicChannelDetailsSchema>;

export const LandingPageChannelsSchema = z
	.object({
		id: z.string(),
		login: z.string(),
		displayName: z.string(),
		profileImage: z.string().url(),
		isLive: z.boolean(),
	})
	.array();

export type LandingPageChannels = z.infer<typeof LandingPageChannelsSchema>;

export const ChannelSettingsSchema = z.object({
	botActiveWhenOffline: z.boolean(),
	guessDelayTime: z.number().min(0).max(60),
	usagePublic: z.boolean(),
});