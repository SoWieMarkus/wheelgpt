import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import * as z from "zod";
import { database } from "../database";
import { logger } from "../utils";

const TWITCH_MESSAGE_TYPE_VERIFICATION = "webhook_callback_verification";
const TWITCH_MESSAGE_TYPE_NOTIFICATION = "notification";
const TWITCH_MESSAGE_TYPE_REVOCATION = "revocation";

const TWITCH_MESSAGE_TYPE_HEADER = "twitch-eventsub-message-type";

const TwitchEventSubHeaderSchema = z.object({
	[TWITCH_MESSAGE_TYPE_HEADER]: z.string(),
});

const TwitchWebhookEventSchema = z.object({
	subscription: z.object({
		id: z.string(),
		status: z.string(),
		type: z.string(),
		version: z.string(),
		condition: z.object({
			broadcaster_user_id: z.string(),
		}),
		transport: z.object({
			method: z.string(),
			callback: z.string(),
		}),
		created_at: z.string(),
		cost: z.number(),
	}),
	event: z.object({
		broadcaster_user_id: z.string(),
		broadcaster_user_login: z.string(),
		broadcaster_user_name: z.string(),
	}),
});

export const streamStateWebhook: RequestHandler = async (request, response) => {
	if (!request.legitTwitchSignature) {
		throw createHttpError(403, "Unauthorized. Invalid Twitch signature.");
	}

	const { success, data, error } = TwitchEventSubHeaderSchema.safeParse(request.headers);
	if (!success) {
		logger.error(`Invalid Twitch EventSub headers ${z.prettifyError(error)} `);
		throw createHttpError(403, "Bad Request. Invalid Twitch EventSub headers.");
	}

	const messageType = data[TWITCH_MESSAGE_TYPE_HEADER];
	const body = JSON.parse(request.body.toString("utf8"));

	switch (messageType) {
		case TWITCH_MESSAGE_TYPE_VERIFICATION: {
			const challenge = body.challenge;
			if (!challenge) {
				logger.error("Verification challenge missing in request body");
				throw createHttpError(400, "Bad Request. Missing verification challenge.");
			}
			response.status(200).send(challenge);
			return;
		}
		case TWITCH_MESSAGE_TYPE_NOTIFICATION: {
			const { success, data } = TwitchWebhookEventSchema.safeParse(body);
			if (!success) {
				logger.error("Failed to parse Twitch event");
				console.error(body);
				throw createHttpError(400, "Bad Request. Invalid Twitch event data.");
			}

			const channelId = data.event.broadcaster_user_id;
			const type = data.subscription.type;
			const channel = await database.channel.update({
				where: { id: channelId },
				data: {
					isLive: type === "stream.online",
				},
			});

			// If the channel is offline and the bot is not active when offline,
			// delete the map, room and guesses
			if (channel?.botActiveWhenOffline === false && type === "stream.offline") {
				// We use deleteMany here since delete would throw an error if the map/room doesn't exist
				await database.trackmaniaMap.deleteMany({
					where: { channelId },
				});
				await database.guess.deleteMany({
					where: { channelId },
				});
				await database.trackmaniaRoom.deleteMany({
					where: { channelId },
				});
			}

			logger.info(
				`Channel ${data.event.broadcaster_user_name} updated to ${type === "stream.online" ? "live" : "offline"} `,
			);
			response.status(204).send();
			return;
		}
		case TWITCH_MESSAGE_TYPE_REVOCATION: {
			const subscriptionId = body.subscription.id;
			if (!subscriptionId) {
				logger.error("Subscription ID missing in revocation request");
				throw createHttpError(400, "Bad Request. Missing subscription ID.");
			}
			logger.warn("Received Twitch revocation");
			response.status(204).send();
			return;
		}
		default:
			logger.error("Unknown Twitch EventSub message type", { messageType });
			throw createHttpError(400, "Bad Request. Unknown Twitch EventSub message type.");
	}
};
