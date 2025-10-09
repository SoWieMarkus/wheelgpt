import crypto from "node:crypto";
import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import * as z from "zod";
import { env, logger } from "../utils";

const TWITCH_MESSAGE_ID = "twitch-eventsub-message-id";
const TWITCH_MESSAGE_TIMESTAMP = "twitch-eventsub-message-timestamp";
const TWITCH_MESSAGE_SIGNATURE = "twitch-eventsub-message-signature";

export const verifyTwitchSignature: RequestHandler = (request, _, next) => {
	const secret = env.TWITCH_EVENTSUB_SECRET;

	const TwitchEventSubHeaderSchema = z.object({
		[TWITCH_MESSAGE_ID]: z.string(),
		[TWITCH_MESSAGE_TIMESTAMP]: z.string(),
		[TWITCH_MESSAGE_SIGNATURE]: z.string(),
	});

	const { success, data, error } = TwitchEventSubHeaderSchema.safeParse(request.headers);

	if (!success) {
		const httpError = createHttpError(403, "Bad Request. Invalid Twitch EventSub headers.");
		logger.error(`Failed to parse Twitch EventSub headers ${z.prettifyError(error)}`);
		return next(httpError);
	}

	const {
		[TWITCH_MESSAGE_ID]: messageId,
		[TWITCH_MESSAGE_TIMESTAMP]: timestamp,
		[TWITCH_MESSAGE_SIGNATURE]: signature,
	} = data;

	const hmacMessage = messageId + timestamp + request.body.toString("utf8");
	const computedHmac = `sha256=${crypto.createHmac("sha256", secret).update(hmacMessage).digest("hex")}`;

	const verified = crypto.timingSafeEqual(Buffer.from(signature), Buffer.from(computedHmac));

	if (!verified) {
		const httpError = createHttpError(403, "Unauthorized. Invalid Twitch EventSub signature.");
		logger.error("Twitch EventSub signature verification failed", { messageId, timestamp });
		return next(httpError);
	}

	request.legitTwitchSignature = true;

	next();
};
