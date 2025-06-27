import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import { z } from "zod";
import { database } from "../database";
import { env } from "../utils";

const PayloadSchema = z.object({
	channelId: z.string().optional(), // Channel login, which was used in the past to generate tokens. To still support existing tokens it is kept here.
	id: z.string().optional(), // Channel ID, which is the new standard for tokens. One of these two needs to be present.
	token: z.string(),
});

export const requiresPluginAuthentication: RequestHandler = (request, response, next) => {
	const token = request.headers.authorization;
	if (token === undefined || token === "") {
		const httpError = createHttpError(401, "Unauthorized. Please provide a valid token.");
		return next(httpError);
	}

	jwt.verify(token, env.JWT_SECRET_CHANNEL, async (error, decoded) => {
		if (error !== null || decoded === undefined) {
			const httpError = createHttpError(401, "Unauthorized. Please provide a valid token.");
			return next(httpError);
		}

		const { success, data } = PayloadSchema.safeParse(decoded);

		if (!success) {
			const httpError = createHttpError(401, "Unauthorized. Invalid token payload.");
			return next(httpError);
		}

		const { channelId, id, token } = data;

		// Check if exactly one of channelId or id is provided
		// channelId === login, used in the past to identifiy channels
		// To support existing tokens, we keep this check
		if (channelId === undefined && id === undefined) {
			const httpError = createHttpError(401, "Unauthorized. Invalid token payload.");
			return next(httpError);
		}

		if (channelId !== undefined && id !== undefined) {
			const httpError = createHttpError(401, "Unauthorized. Invalid token payload.");
			return next(httpError);
		}

		const channel =
			channelId === undefined
				? await database.channel.findUnique({
						where: { id: id },
					})
				: await database.channel.findUnique({
						where: { login: channelId },
					});

		if (channel === null) {
			const httpError = createHttpError(401, "Unauthorized. User not found.");
			return next(httpError);
		}

		if (channel.token !== token) {
			const httpError = createHttpError(401, "Unauthorized. Invalid token.");
			return next(httpError);
		}

		if (!channel.botActiveWhenOffline && !channel.isLive) {
			response.status(200).json({
				message: "Channel is offline.",
			});
			return;
		}

		request.channelId = channel.id;
		request.login = channel.login;
		next();
	});
};
