import { RequestHandler } from "express";
import { z } from "zod";
import { Twitch } from "../apis";
import { database } from "../database";
import { env, Validation } from "../utils";
import { WheelGPT } from "../bot/bot";
import jwt from "jsonwebtoken";
import * as uuid from "uuid";
import createHttpError from "http-errors";

const clientId = env.TWITCH_CLIENT_ID;
const redirectUri = env.TWITCH_REDIRECT_URL;

export const twitch: RequestHandler = (request, response) => {
    const scope = 'user:read:email';
    const authUrl = `https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=${clientId}&redirect_uri=${redirectUri}&scope=${scope}`;
    response.redirect(authUrl);
};

export const authenticate: RequestHandler = async (request, response) => {
    const AuthenticationBodySchema = z.object({ code: z.string() });
    const { code } = AuthenticationBodySchema.parse(request.body);

    const accessToken = await Twitch.getAccessToken(code);
    const channelId = await Twitch.getUser(accessToken);

    const existingChannel = await database.channel.findUnique({
        where: { channelId }
    });

    if (existingChannel === null) {
        await WheelGPT.register(channelId);
    }

    const webToken = jwt.sign({ channelId: channelId }, env.JWT_SECRET_WEB, { expiresIn: "1d" });
    response.status(200).json({ webToken });
}

export const remove: RequestHandler = async (request, response) => {
    const channelId = Validation.validateRequiresAuthentication(request);
    await WheelGPT.remove(channelId);
    response.status(200).json({});
};

export const getToken: RequestHandler = async (request, response) => {
    const channelId = Validation.validateRequiresAuthentication(request);
    const channel = await database.channel.findUnique({
        where: { channelId }
    });
    if (channel === null) {
        throw createHttpError(400, "Bad Request.");
    }
    const channelToken = jwt.sign({ channelId: channelId, token: channel.token }, env.JWT_SECRET_CHANNEL);
    response.status(200).json({ channelToken });
};

export const updateToken: RequestHandler = async (request, response) => {
    const channelId = Validation.validateRequiresAuthentication(request);
    const channel = await database.channel.update({ where: { channelId }, data: { token: uuid.v4() } });
    if (channel === null) {
        throw createHttpError(400, "Bad Request.");
    }
    const channelToken = jwt.sign({ channelId: channelId, token: channel.token }, env.JWT_SECRET_CHANNEL);
    response.status(200).json({ channelToken });
};