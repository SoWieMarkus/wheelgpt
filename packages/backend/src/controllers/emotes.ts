import type { RequestHandler } from "express";
import { Emote } from "../bot/core/emotes";

export const getEmotes: RequestHandler = async (request, response) => {
	response.status(200).json(Object.values(Emote));
};
