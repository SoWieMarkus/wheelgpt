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

const AccessTokenSchema = z.object({
	access_token: z.string(),
});

export const getUserAccessToken = async (code: string): Promise<string | null> => {
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

	const { success, data, error } = AccessTokenSchema.safeParse(response.data);
	if (!success) {
		logger.error("Failed to parse user access token response", {
			error: error.errors[0].message,
		});
		failedTwitchAPIMetric.inc({ endpoint: "user_access_token" });
		return null;
	}
	return data.access_token;
};

export const getAppAccessToken = async (): Promise<string | null> => {
	const url = "https://id.twitch.tv/oauth2/token";
	const response = await axios.post(url, null, {
		params: {
			client_id: clientId,
			client_secret: clientSecret,
			grant_type: "client_credentials",
		},
	});

	const { success, data, error } = AccessTokenSchema.safeParse(response.data);
	if (!success) {
		logger.error("Failed to parse user access token response", {
			error: error.errors[0].message,
		});
		failedTwitchAPIMetric.inc({ endpoint: "app_access_token" });
		return null;
	}

	appAccessToken = data.access_token;
	return appAccessToken;
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

export const getUserOnlineStatus = async (userIds: string[]) => { };
