import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { z } from "zod";
import { database } from "../database";
import { TMX } from "../external";
import { logger } from "../utils";
import { wheelgpt } from "../wheelgpt";

const TrackmaniaMapSchema = z.object({
	name: z.string(),
	uid: z.string(),
	author: z.string(),
	authorTime: z.number().min(0),
	goldTime: z.number().min(0),
	silverTime: z.number().min(0),
	bronzeTime: z.number().min(0),
	championTime: z.number().min(0),
});

const TrackmaniaMapPostSchema = z.object({
	map: TrackmaniaMapSchema.nullable().optional(),
});

export const updateMap: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data, error } = TrackmaniaMapPostSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	const { map } = data;

	const exisitingMap = await database.trackmaniaMap.findUnique({
		where: { channelId },
	});

	if (exisitingMap || !map) {
		await database.trackmaniaMap.delete({
			where: { channelId },
		});
	}

	if (!map) {
		response.status(200).json({ message: "Map deleted successfully." });
		return;
	}

	await database.trackmaniaMap.create({
		data: {
			channelId,
			...map,
		},
	});

	TMX.getTrackmaniaExchangeData(map.uid)
		.then(async (tmxData) => {
			if (!tmxData) return;
			await database.trackmaniaMap.update({
				where: { channelId },
				data: { tmxId: tmxData.TrackID, worldRecord: null },
			});
		})
		.catch((error) => {
			console.error("Failed to update map with TMX data:", error);
			logger.error("Failed to update map with TMX data", {
				mapUid: map.uid,
			});
		});

	response.status(200).json({ message: "Successfully updated the map." });
};

const PersonalBestSchema = z.object({
	time: z.number().min(0),
});

export const updatePersonalBest: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data, error } = PersonalBestSchema.safeParse(request.body);
	console.log(success, data, error);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	const { time } = data;
	wheelgpt.notifyNewPB(channelId, time);
	response.status(200).json({ message: "Personal best updated successfully." });
};
