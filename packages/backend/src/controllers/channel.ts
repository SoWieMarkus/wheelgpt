import { Schema } from "@wheelgpt/shared";
import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { database } from "../database";

const query = {
	id: true,
	displayName: true,
	profileImage: true,
	botActiveWhenOffline: true,
	guessDelayTime: true,
	usagePublic: true,
};

export const me: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Unauthorized: Authentication required.");
	}

	const channel = await database.channel.findUnique({
		where: { id: channelId },
		select: query,
	});

	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}

	response.status(200).json(channel);
};

export const updateSettings: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Unauthorized: Authentication required.");
	}

	const { success, data, error } = Schema.channel.settings.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, `Invalid request body: ${error.errors[0].message}`);
	}

	const updatedChannel = await database.channel.update({
		where: { id: channelId },
		data: {
			botActiveWhenOffline: data.botActiveWhenOffline,
			guessDelayTime: data.guessDelayTime,
			usagePublic: data.usagePublic,
		},
		select: query,
	});

	response.status(200).json(updatedChannel);
};
