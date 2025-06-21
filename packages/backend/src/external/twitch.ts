import axios from "axios";
import { Counter } from "prom-client";
import { z } from "zod";
import { prometheus } from "../prometheus";
import { env, logger } from "../utils";

export const failedTwitchAPIMetric = new Counter({
	name: "wheelgpt_twitch_failed_requests_total",
	help: "Total number of failed requests to Twitch API",
	labelNames: ["endpoint"],
	registers: [prometheus],
});

const clientId = env.TWITCH_CLIENT_ID;
const clientSecret = env.TWITCH_CLIENT_SECRET;

let appAccessToken: string | null = null;
let appAccessTokenExpiration: number | null = null;

const AppAccessTokenSchema = z.object({
	access_token: z.string(),
	expires_in: z.number(),
});

// Request a new app access token from Twitch
export const requestAppAccessToken = async (): Promise<string | null> => {
	const url = "https://id.twitch.tv/oauth2/token";
	const response = await axios.post(url, null, {
		params: {
			client_id: clientId,
			client_secret: clientSecret,
			grant_type: "client_credentials",
		},
	});

	const { success, data, error } = AppAccessTokenSchema.safeParse(response.data);
	if (!success) {
		logger.error("Failed to parse user access token response", {
			error: error.errors[0].message,
		});
		failedTwitchAPIMetric.inc({ endpoint: "app_access_token" });
		return null;
	}

	appAccessToken = data.access_token;
	appAccessTokenExpiration = data.expires_in;
	return appAccessToken;
};

// Get the app access token, caching it if it's still valid
const getAppAccessToken = async (): Promise<string> => {
	if (appAccessToken && appAccessTokenExpiration && Date.now() < appAccessTokenExpiration * 1000) {
		return appAccessToken;
	}
	const token = await requestAppAccessToken();
	if (token === null) throw new Error("Failed to retrieve app access token");
	return token;
};

const UserAccessTokenSchema = z.object({
	access_token: z.string(),
});

// Request a user access token using the authorization code returned from the Twitch OAuth flow
export const requestUserAccessToken = async (code: string): Promise<string | null> => {
	const url = "https://id.twitch.tv/oauth2/token";
	const response = await axios.post(url, null, {
		params: {
			client_id: clientId,
			client_secret: clientSecret,
			code,
			grant_type: "authorization_code",
			redirect_uri: "https://wheelgpt.dev",
		},
	});

	const { success, data, error } = UserAccessTokenSchema.safeParse(response.data);
	if (!success) {
		logger.error("Failed to parse user access token response", {
			error: error.errors[0].message,
		});
		failedTwitchAPIMetric.inc({ endpoint: "user_access_token" });
		return null;
	}
	return data.access_token;
};

const UserSchema = z.object({
	id: z.string(),
	login: z.string(),
	display_name: z.string(),
	profile_image_url: z.string(),
});

const HelixUsersSchema = z.object({
	data: UserSchema.array(),
});

// Used to get the user information by a user access token to identify who is trying to log in
export const getUser = async (accessToken: string) => {
	const userUrl = "https://api.twitch.tv/helix/users";
	const response = await axios.get(userUrl, {
		headers: {
			Authorization: `Bearer ${accessToken}`,
			"Client-Id": env.TWITCH_CLIENT_ID,
		},
	});

	const { success, data, error } = HelixUsersSchema.safeParse(response.data);
	if (!success) {
		logger.error("Failed to parse user response", {
			error: error.errors[0].message,
		});
		failedTwitchAPIMetric.inc({ endpoint: "get_user" });
		return null;
	}
	if (data.data.length === 0) {
		logger.error("No user data found in response");
		failedTwitchAPIMetric.inc({ endpoint: "get_user" });
		return null;
	}

	return data.data[0];
};

// Used to update the channel information in the database on boot up
export const getUsers = async (channelIds: string[]) => {
	const accessToken = await getAppAccessToken();
	const userUrl = "https://api.twitch.tv/helix/users";

	// Split channelIds into chunks of 100, as Twitch API has a limit of 100 IDs per request
	const chunkSize = 100;
	const chunks = [];
	for (let i = 0; i < channelIds.length; i += chunkSize) {
		chunks.push(channelIds.slice(i, i + chunkSize));
	}

	const users: z.infer<typeof UserSchema>[] = [];
	for (const chunk of chunks) {
		const response = await axios.get(userUrl, {
			headers: {
				Authorization: `Bearer ${accessToken}`,
				"Client-Id": env.TWITCH_CLIENT_ID,
			},
			params: {
				id: chunk,
			},
		});

		const { success, data, error } = HelixUsersSchema.safeParse(response.data);
		if (!success) {
			logger.error("Failed to parse user response", {
				error: error.errors[0].message,
			});
			failedTwitchAPIMetric.inc({ endpoint: "get_users" });
			continue;
		}

		users.push(...data.data);
	}
	return users;
};

const StreamsSchema = z.object({
	user_id: z.string(),
	type: z.string(),
});

const HelixStreamsSchema = z.object({
	data: StreamsSchema.array(),
});

// Used to initially check if a channel is live or not
// See: https://dev.twitch.tv/docs/api/reference/#get-streams
export const getStreams = async (channelIds: string[]) => {
	const accessToken = await getAppAccessToken();
	const streamsUrl = "https://api.twitch.tv/helix/streams";

	// Split channelIds into chunks of 100, as Twitch API has a limit of 100 IDs per request
	const chunkSize = 100;
	const chunks = [];
	for (let i = 0; i < channelIds.length; i += chunkSize) {
		chunks.push(channelIds.slice(i, i + chunkSize));
	}
	const streams: z.infer<typeof StreamsSchema>[] = [];
	for (const chunk of chunks) {
		const response = await axios.get(streamsUrl, {
			headers: {
				Authorization: `Bearer ${accessToken}`,
				"Client-Id": env.TWITCH_CLIENT_ID,
			},
			params: {
				user_id: chunk,
			},
		});

		const { success, data, error } = HelixStreamsSchema.safeParse(response.data);
		if (!success) {
			logger.error("Failed to parse streams response", {
				error: error.errors[0].message,
			});
			failedTwitchAPIMetric.inc({ endpoint: "get_streams" });
			continue;
		}

		streams.push(...data.data);
	}
	return streams;
};

// Zod schema for EventSub subscriptions
const EventSubSubscriptionSchema = z.object({
	id: z.string(),
	type: z.string(),
	status: z.string(),
	condition: z.record(z.string(), z.string()),
	transport: z.object({
		method: z.string(),
		callback: z.string(),
	}),
	created_at: z.string(),
});
const HelixEventSubSchema = z.object({
	data: EventSubSubscriptionSchema.array(),
	total: z.number(),
	total_cost: z.number(),
	max_total_cost: z.number(),
});

// Fetch all registered EventSub webhooks
export const getRegisteredWebhooks = async () => {
	const accessToken = await getAppAccessToken();
	const url = "https://api.twitch.tv/helix/eventsub/subscriptions";
	const response = await axios.get(url, {
		headers: {
			Authorization: `Bearer ${accessToken}`,
			"Client-Id": env.TWITCH_CLIENT_ID,
		},
	});
	const { success, data, error } = HelixEventSubSchema.safeParse(response.data);
	if (!success) {
		logger.error("Failed to parse EventSub subscriptions", { error: error.errors[0].message });
		return [];
	}
	return data.data;
};

// Register a new webhook for a channel (online/offline)
export const registerWebhook = async (
	type: "stream.online" | "stream.offline",
	broadcasterUserId: string,
	callback: string,
) => {
	const accessToken = await getAppAccessToken();
	const url = "https://api.twitch.tv/helix/eventsub/subscriptions";
	await axios.post(
		url,
		{
			type,
			version: "1",
			condition: { broadcaster_user_id: broadcasterUserId },
			transport: {
				method: "webhook",
				callback,
				secret: env.TWITCH_EVENTSUB_SECRET,
			},
		},
		{
			headers: {
				Authorization: `Bearer ${accessToken}`,
				"Client-Id": env.TWITCH_CLIENT_ID,
				"Content-Type": "application/json",
			},
		},
	);
};

// Remove a webhook by ID
export const removeWebhook = async (id: string) => {
	const accessToken = await getAppAccessToken();
	const url = `https://api.twitch.tv/helix/eventsub/subscriptions?id=${id}`;
	await axios.delete(url, {
		headers: {
			Authorization: `Bearer ${accessToken}`,
			"Client-Id": env.TWITCH_CLIENT_ID,
		},
	});
};

// Main logic to sync webhooks on boot
export const syncWebhooks = async (channelIds: string[], callbackUrl: string) => {
	const registered = await getRegisteredWebhooks();
	console.log(registered);

	// Build a set of required webhooks
	const required = new Set<string>();
	for (const id of channelIds) {
		required.add(`stream.online:${id}`);
		required.add(`stream.offline:${id}`);
	}

	// Build a map of currently registered webhooks
	const registeredMap = new Map<string, { id: string; type: string; broadcaster_user_id: string }>();
	for (const sub of registered) {
		if (sub.status !== "enabled") {
			console.log(sub.status);
			logger.info(`Skipping disabled webhook: ${sub.id}`);
			continue;
		}
		if ((sub.type === "stream.online" || sub.type === "stream.offline") && sub.condition.broadcaster_user_id) {
			registeredMap.set(`${sub.type}:${sub.condition.broadcaster_user_id}`, {
				id: sub.id,
				type: sub.type,
				broadcaster_user_id: sub.condition.broadcaster_user_id,
			});
		}
	}

	// Register missing webhooks
	for (const id of channelIds) {
		for (const type of ["stream.online", "stream.offline"] as const) {
			const key = `${type}:${id}`;
			if (!registeredMap.has(key)) {
				logger.info(`Registering webhook: ${type} for channel ${id}`);
				await registerWebhook(type, id, callbackUrl);
			}
		}
	}

	// Remove unused webhooks
	for (const [key, sub] of registeredMap.entries()) {
		if (!required.has(key)) {
			logger.info(`Removing unused webhook: ${sub.type} for channel ${sub.broadcaster_user_id}`);
			await removeWebhook(sub.id);
		}
	}
};
