import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import * as uuid from "uuid";
import z from "zod";
import { database } from "../database";
import { Twitch } from "../external";
import { env, logger } from "../utils";
import { wheelgpt } from "../wheelgpt";

const clientId = env.TWITCH_CLIENT_ID;
const redirectUri = env.TWITCH_REDIRECT_URL;

export const twitchLoginRequest: RequestHandler = async (_, response) => {
	const scope = "user:read:email";
	const authUrl = `https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=${clientId}&redirect_uri=${redirectUri}&scope=${scope}`;
	response.redirect(authUrl);
};

const TwitchCodeSchema = z.object({
	code: z.string(),
});

const generatePluginToken = (id: string, token: string) => {
	return jwt.sign({ id, token }, env.JWT_SECRET_CHANNEL);
};

const generateWebToken = (id: string) => {
	return jwt.sign({ id }, env.JWT_SECRET_WEB, { expiresIn: "1d" });
};

export const login: RequestHandler = async (request, response) => {
	const { success, data, error } = TwitchCodeSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	const { code } = data;
	const userAccessToken = await Twitch.requestUserAccessToken(code);

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
		Twitch.addWebhooksByChannel(channel.id)
			.then(() => {
				logger.info(`Added Twitch webhooks for channel ${channel.id}`);
			})
			.catch(() => {
				logger.error("Failed to add Twitch webhooks");
			});
		Twitch.checkStreamState(channel.id)
			.then(async (isLive) => {
				logger.info(`Checked stream state for channel ${channel.id}: ${isLive ? "Live" : "Offline"}`);
				await database.channel.update({
					where: { id: channel.id },
					data: { isLive },
				});
			})
			.catch(() => {
				logger.error("Failed to add Twitch webhooks");
			});
	}

	const webToken = generateWebToken(channel.id);
	response.status(200).json({ webToken });
};

export const remove: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(400, "Channel ID is required.");
	}

	const channel = await database.channel.findUnique({
		where: { id: channelId },
		select: { id: true, login: true },
	});
	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}

	await wheelgpt.remove(channel.login);
	await database.channel.delete({ where: { id: channelId } });
	Twitch.removeWebhooksByChannel(channelId)
		.then(() => {
			logger.info(`Removed Twitch webhooks for channel ${channelId}`);
		})
		.catch((_) => {
			logger.error("Failed to remove Twitch webhooks");
		});
	response.status(204).json({});
};

export const getPluginToken: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(400, "Channel ID is required.");
	}

	const channel = await database.channel.findUnique({
		where: { id: channelId },
		select: { token: true, id: true },
	});

	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}

	const pluginToken = generatePluginToken(channel.id, channel.token);
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
		select: { token: true, id: true },
	});

	if (!channel) {
		throw createHttpError(404, "Channel not found.");
	}

	const pluginToken = generatePluginToken(channel.id, channel.token);
	response.status(200).json({ pluginToken });
};
