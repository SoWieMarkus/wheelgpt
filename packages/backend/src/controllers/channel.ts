import { Schema } from "@wheelgpt/shared";
import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { database } from "../database";
import { wheelgpt } from "../wheelgpt";

const query = {
	id: true,
	login: true,
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

	if (!updatedChannel) {
		throw createHttpError(404, "Channel not found.");
	}

	wheelgpt.reload(updatedChannel.login, updatedChannel.id);
	response.status(200).json(updatedChannel);
};

export const getChannelById: RequestHandler = async (request, response) => {
	const channelId = request.params.channelId;
	if (!channelId) {
		throw createHttpError(400, "Channel ID is required.");
	}
	const channel = await database.channel.findUnique({
		where: { id: channelId },
		select: { id: true, displayName: true, profileImage: true, isLive: true },
	});

	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}
	response.status(200).json(channel);
};

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
		orderBy: {
			isLive: "desc",
		},
		skip: 0,
		take: 30,
	});
	response.status(200).json(channels);
};

