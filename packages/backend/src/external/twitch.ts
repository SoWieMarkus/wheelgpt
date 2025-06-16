import axios from "axios";
import { data } from "cheerio/dist/commonjs/api/attributes";
import { Counter } from "prom-client";
import { z } from "zod";
import { database } from "../database";
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
				id: chunk.join(","),
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
	id: z.string(),
	type: z.string(),
});

const HelixStreamsSchema = z.object({
	data: StreamsSchema.array(),
});

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
				user_id: chunk.join(","),
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
