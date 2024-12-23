import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { TrackmaniaExchange } from "../apis";
import { WheelGPT } from "../bot/bot";
import { GuessResultSchema, TrackmaniaTime } from "../bot/guess/time";
import { TrackmaniaMap, TrackmaniaMapPostSchema } from "../bot/map/map";
import { Log, Validation } from "../utils";

export const updateMap: RequestHandler = async (request, response) => {
    const { map } = TrackmaniaMapPostSchema.parse(request.body);
    const channelId = Validation.validateRequiresAuthentication(request);

    const channel = WheelGPT.getChannel(channelId);
    if (channel === undefined) {
        throw createHttpError(400, "Channel not found.");
    }

    if (map === undefined) {
        Log.complete(`Removed map on channel "${channelId}"`)
        channel.setMap(null);
        response.status(200).json({ message: "Successfully updated the map." });
        return;
    }

    const newMap = new TrackmaniaMap(map);
    const mapUid = map.uid;
    TrackmaniaExchange.getTrackmaniaExchangeId(mapUid).then((tmxId) => {
        newMap.setTrackmaniaExchangeId(tmxId);
        Log.complete(`Trackmania Exchange ID for map "${mapUid}" is "${tmxId}"`);
    }).catch(error => console.error(error));
    channel.setMap(newMap);

    Log.complete(`New map on channel "${channelId}": "${newMap.toString()}"`)
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

    Log.complete(`New pb on channel "${channelId}" ("${trackmaniaTime.toString()}")`)
    setTimeout(() => {
        const response = channel.guessResult(trackmaniaTime);
        console.log(response);
        if (response === null) return;
        WheelGPT.say(channelId, response);
    }, channel.getGuessDelay());

    response.status(200).json({ message: "Successfully updated personal best." });
};