import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import jwt from "jsonwebtoken";
import { env } from "../utils";
import { z } from "zod";

export const requiresWebAuthentication: RequestHandler = (
    request,
    response,
    next,
) => {
    const token = request.headers.authorization;
    if (token === undefined) {
        return next(createHttpError(401, "Please provide a valid token."));
    }

    jwt.verify(token, env.JWT_SECRET_WEB, (error, decoded) => {
        const PayloadSchema = z.object({
            channelId: z.string(),
            exp: z.number()
        })
        if (error !== null || decoded === undefined) {
            return next(createHttpError(401, "Please provide a valid token."));
        }

        const { success, data } = PayloadSchema.safeParse(decoded);
        if (!success) return next(createHttpError(401, "Please provide a valid token."));
        const { channelId, exp } = data;
        if (exp <= Date.now() / 1000) {
            return next(createHttpError(401, "Token has expired."));
        }
        request.channelId = channelId;
        next();
    });
};
