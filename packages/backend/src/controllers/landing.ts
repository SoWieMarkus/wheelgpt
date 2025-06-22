import type { RequestHandler } from "express";
import { database } from "../database";

export const getPublicChannels: RequestHandler = async (_, response) => {
	const channels = await database.channel.findMany({
		where: {
			usagePublic: true,
		},
		select: {
			displayName: true,
			profileImage: true,
			isLive: true,
			login: true,
		},
		skip: 0,
		take: 30,
	});
	response.status(200).json(channels);
};
