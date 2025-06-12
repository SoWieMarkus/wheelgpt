import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import { z } from "zod";
import { database } from "../database";
import { env } from "../utils";

const PayloadSchema = z.object({
	channelId: z.string(),
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

		const { channelId, token } = data;

		const channel = await database.channel.findUnique({
			where: { channelId },
		});

		if (channel === null) {
			const httpError = createHttpError(401, "Unauthorized. User not found.");
			return next(httpError);
		}

		if (channel.token !== token) {
			const httpError = createHttpError(401, "Unauthorized. Invalid token.");
			return next(httpError);
		}

		if (channel.botActiveWhenOffline && !channel.isLive) {
			response.status(200).json({
				message: "Channel is offline, but bot is active.",
			});
			return;
		}

		request.channelId = channelId;
		next();
	});
};
