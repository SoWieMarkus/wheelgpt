import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import { env } from "../utils";
import { z } from "zod";
import { database } from "../database";

export const requiresChannelAuthentication: RequestHandler = (
    request,
    response,
    next,
) => {
    const token = request.headers.authorization;
    if (token === undefined) {
        return next(createHttpError(401, "Please provide a valid token."));
    }

    jwt.verify(token, env.JWT_SECRET_CHANNEL, async (error, decoded) => {

        const PayloadSchema = z.object({
            channelId: z.string(),
            token: z.string()
        })

        if (error !== null || decoded === undefined) {
            return next(createHttpError(401, "Please provide a valid token."));
        }

        const { success, data } = PayloadSchema.safeParse(decoded);
        if (!success) return next(createHttpError(401, "Please provide a valid token."));
        const { channelId, token } = data;

        const channel = await database.channel.findUnique({ where: { channelId } });
        if (channel === null || channel.token !== token) return next(createHttpError(401, "Please provide a valid token."));

        request.channelId = channelId;
        next();
    });
};
