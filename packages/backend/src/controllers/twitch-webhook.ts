import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { z } from "zod";
import { logger } from "../utils";

const TWITCH_MESSAGE_TYPE_VERIFICATION = "webhook_callback_verification";
const TWITCH_MESSAGE_TYPE_NOTIFICATION = "notification";
const TWITCH_MESSAGE_TYPE_REVOCATION = "revocation";

const TWITCH_MESSAGE_TYPE_HEADER = "twitch-eventsub-message-type";

const TwitchEventSubHeaderSchema = z.object({
	[TWITCH_MESSAGE_TYPE_HEADER]: z.string(),
});

const TwitchWebhookEventSchema = z.object({
	type: z.string(),
});

export const streamStateWebhook: RequestHandler = async (request, response) => {
	if (!request.legitTwitchSignature) {
		throw createHttpError(403, "Unauthorized. Invalid Twitch signature.");
	}

	const { success, data, error } = TwitchEventSubHeaderSchema.safeParse(request.headers);
	if (!success) {
		logger.error("Failed to parse Twitch EventSub headers", { error: error.errors[0].message });
		throw createHttpError(403, "Bad Request. Invalid Twitch EventSub headers.");
	}
	const messageType = data[TWITCH_MESSAGE_TYPE_HEADER];

	const body = JSON.parse(request.body.toString("utf8"));

	console.log("_________________________________");
	console.log(messageType);
	console.log(request.headers);
	console.log(body);

	switch (messageType) {
		case TWITCH_MESSAGE_TYPE_VERIFICATION: {
			const challenge = body.challenge;
			console.log("Challenge:", challenge);
			if (!challenge) {
				logger.error("Verification challenge missing in request body");
				throw createHttpError(400, "Bad Request. Missing verification challenge.");
			}
			response.status(200).send(challenge);
			return;
		}
		case TWITCH_MESSAGE_TYPE_NOTIFICATION: {
			const event = TwitchWebhookEventSchema.safeParse(body.event);
			if (!event.success) {
				logger.error("Failed to parse Twitch event", { error: event.error.errors[0].message });
				throw createHttpError(400, "Bad Request. Invalid Twitch event data.");
			}

			// TODO update the channel's live status in the database

			response.status(204).send();
			return;
		}
		case TWITCH_MESSAGE_TYPE_REVOCATION: {
			const subscriptionId = body.subscription.id;
			if (!subscriptionId) {
				logger.error("Subscription ID missing in revocation request");
				throw createHttpError(400, "Bad Request. Missing subscription ID.");
			}
			// Handle revocation (e.g., remove the subscription from the database)
			logger.info("Received Twitch revocation", { subscriptionId });
			response.status(204).send();
			return;
		}
		default:
			logger.error("Unknown Twitch EventSub message type", { messageType });
			throw createHttpError(400, "Bad Request. Unknown Twitch EventSub message type.");
	}
};
