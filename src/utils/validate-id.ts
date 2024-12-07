import type { Request } from "express";
import createHttpError from "http-errors";
import { z } from "zod";

export const validateChannelId = (request: Request) => {
    const ChannelIdSchema = z.object({ channelId: z.string() });
    const { channelId } = ChannelIdSchema.parse(request.params);
    return channelId;
};

export const validateRequiresAuthentication = (request: Request) => {
    const { channelId } = request;
    if (channelId === undefined) {
        throw createHttpError(403, "Please provide a valid token.");
    }
    return channelId;
};

