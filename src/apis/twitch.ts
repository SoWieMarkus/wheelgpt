import axios from "axios";
import createHttpError from "http-errors";
import { z } from "zod";
import { env } from "../utils";

const clientId = env.TWITCH_CLIENT_ID;
const clientSecret = env.TWITCH_CLIENT_SECRET;

const AccessTokenSchema = z.object({
    access_token: z.string(),
});
const UsersSchema = z.object({
    data: z.array(z.object({ login: z.string() }))
});

export const getAccessToken = async (code: string) => {
    const url = "https://id.twitch.tv/oauth2/token";
    const tokenResponse = await axios.post(url, null, {
        params: {
            client_id: clientId,
            client_secret: clientSecret,
            code,
            grant_type: "authorization_code",
            redirect_uri: "https://wheelgpt.dev",
        },
    });

    const accessToken = AccessTokenSchema.parse(tokenResponse.data);
    return accessToken.access_token;
};

export const getUser = async (accessToken: string) => {
    const userUrl = "https://api.twitch.tv/helix/users";
    const userResponse = await axios.get(userUrl, {
        headers: {
            Authorization: `Bearer ${accessToken}`,
            "Client-Id": env.TWITCH_CLIENT_ID,
        },
    });

    const users = UsersSchema.parse(userResponse.data).data;
    if (users.length === 0) {
        throw createHttpError(400, "Illegal access token.");
    }

    const { login } = users[0];
    return login;
};