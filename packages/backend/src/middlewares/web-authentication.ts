import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import { z } from "zod";
import { database } from "../database";
import { env } from "../utils";

const PayloadSchema = z.object({
	id: z.string(),
	exp: z.number(),
});

export const requiresWebAuthentication: RequestHandler = (request, _, next) => {
	const token = request.headers.authorization;
	if (token === undefined || token === "") {
		const httpError = createHttpError(401, "Unauthorized. Please provide a valid token.");
		return next(httpError);
	}

	jwt.verify(token, env.JWT_SECRET_WEB, async (error, decoded) => {
		if (error !== null || decoded === undefined) {
			const httpError = createHttpError(401, "Unauthorized. Please provide a valid token.");
			return next(httpError);
		}

		const { success, data } = PayloadSchema.safeParse(decoded);

		if (!success) {
			const httpError = createHttpError(401, "Unauthorized. Invalid token payload.");
			return next(httpError);
		}

		const { id, exp } = data;
		const now = Date.now() / 1000;
		if (exp < now) {
			const httpError = createHttpError(401, "Unauthorized. Token has expired.");
			return next(httpError);
		}

		const channel = await database.channel.findUnique({
			where: { id },
		});

		if (channel === null) {
			const httpError = createHttpError(401, "Unauthorized. User not found.");
			return next(httpError);
		}

		request.channelId = channel.id;
		next();
	});
};
