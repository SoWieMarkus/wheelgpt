import type { RequestHandler } from "express";
import { Emote } from "../bot/core/emotes";
import { database } from "../database";

export const getPublicChannels: RequestHandler = async (request, response) => {
	const channels = await database.channel.findMany({
		where: {
			usagePublic: true,
		},
		select: {
			displayName: true,
			profileImage: true,
			isLive: true,
		},
		skip: 0,
		take: 20,
	});

	const emotes = Object.values(Emote);
	response.status(200).json({ channels, emotes });
};
