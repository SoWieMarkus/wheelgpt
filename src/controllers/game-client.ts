import { RequestHandler } from "express";
import createHttpError from "http-errors";
import { TrackmaniaExchange } from "../apis";
import { WheelGPT } from "../bot/bot";
import { GuessResultSchema, TrackmaniaTime } from "../bot/guess/time";
import { TrackmaniaMap, TrackmaniaMapPostSchema } from "../bot/map/map";
import { Validation } from "../utils";

export const updateMap: RequestHandler = async (request, response) => {
    const { map } = TrackmaniaMapPostSchema.parse(request.body);
    const channelId = Validation.validateRequiresAuthentication(request);

    const channel = WheelGPT.getChannel(channelId);
    if (channel === undefined) {
        throw createHttpError(400, "Channel not found.");
    }

    if (map === undefined) {
        channel.setMap(null);
        response.status(200).json({ message: "Successfully updated the map." });
        return;
    }

    const newMap = new TrackmaniaMap(map);
    const mapUid = map.uid;
    const tmxId = await TrackmaniaExchange.getTrackmaniaExchangeId(mapUid);

    newMap.setTrackmaniaExchangeId(tmxId);
    channel.setMap(newMap);
    response.status(200).json({ message: "Successfully updated the map." });
}

export const updatePersonalBest: RequestHandler = (request, response, next) => {
    const { time } = GuessResultSchema.parse(request.body);
    const channelId = Validation.validateRequiresAuthentication(request);
    const trackmaniaTime = new TrackmaniaTime(time);
    const channel = WheelGPT.getChannel(channelId);

    if (channel === undefined) {
        throw createHttpError(400, `Channel with id ${channelId} not found.`);
    }

    setTimeout(() => {
        const response = channel.guessResult(trackmaniaTime);
        if (response === null) return;
        WheelGPT.say(channelId, response);
    }, channel.getGuessDelay() * 1000);

    response.status(200).json({ message: "Successfully updated personal best." });
};