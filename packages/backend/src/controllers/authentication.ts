import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import * as uuid from "uuid";
import { database } from "../database";
import { Twitch } from "../external";
import { env } from "../utils";
import { wheelgpt } from "../wheelgpt";

const clientId = env.TWITCH_CLIENT_ID;
const redirectUri = env.TWITCH_REDIRECT_URL;

export const twitchLoginRequest: RequestHandler = async (request, response) => {
	const scope = "user:read:email";
	const authUrl = `https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=${clientId}&redirect_uri=${redirectUri}&scope=${scope}`;
	response.redirect(authUrl);
};

export const login: RequestHandler = async (request, response) => {
	const { success, data, error } = request.body;
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	const { code } = data;
	const userAccessToken = await Twitch.getUserAccessToken(code);

	if (!userAccessToken) {
		throw createHttpError(400, "Failed to retrieve user access token.");
	}

	const user = await Twitch.getUser(userAccessToken);
	if (!user) {
		throw createHttpError(400, "Failed to retrieve user information.");
	}

	const displayName = user.display_name;
	const profileImage = user.profile_image_url;

	const existingChannel = await database.channel.findUnique({
		where: { id: user.id },
	});

	const channel = await database.channel.upsert({
		where: { id: user.id },
		create: {
			id: user.id,
			login: user.login,
			displayName,
			profileImage,
			token: uuid.v4(),
		},
		update: {
			displayName,
			profileImage,
		},
		select: {
			id: true,
			login: true,
			displayName: true,
			guessDelayTime: true,
			botActiveWhenOffline: true,
			usagePublic: true,
		},
	});

	if (!existingChannel) {
		wheelgpt.register(channel);
	}

	const webToken = jwt.sign({ channelId: user.id }, env.JWT_SECRET_WEB, { expiresIn: "1d" });
	response.status(200).json({ webToken, displayName, profileImage });
};

export const remove: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(400, "Channel ID is required.");
	}

	await wheelgpt.remove(channelId);
	await database.channel.delete({ where: { id: channelId } });
	response.status(204).json({});
};

export const getPluginToken: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(400, "Channel ID is required.");
	}

	const channel = await database.channel.findUnique({
		where: { id: channelId },
		select: { token: true },
	});

	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}

	const pluginToken = jwt.sign({ channelId, token: channel.token }, env.JWT_SECRET_CHANNEL);
	response.status(200).json({ pluginToken });
};

export const updatePluginToken: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(400, "Channel ID is required.");
	}

	const channel = await database.channel.update({
		where: { id: channelId },
		data: { token: uuid.v4() },
		select: { token: true, login: true },
	});

	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}

	const pluginToken = jwt.sign({ id: channel.login, token: channel.token }, env.JWT_SECRET_CHANNEL);
	response.status(200).json({ pluginToken });
};
