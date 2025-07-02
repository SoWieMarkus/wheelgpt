import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { z } from "zod";
import { database } from "../database";
import { TMX } from "../external";
import { logger } from "../utils";
import { wheelgpt } from "../wheelgpt";

const TrackmaniaMapSchema = z.object({
	name: z.string().min(1).max(1000),
	uid: z.string().min(1).max(1000),
	author: z.string().min(1).max(1000),
	authorTime: z.number().min(0),
	goldTime: z.number().min(0),
	silverTime: z.number().min(0),
	bronzeTime: z.number().min(0),
	championTime: z.number().min(0),
});

export const updateMap: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data: map, error } = TrackmaniaMapSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	// Delete all guesses for the channel when the map has changed
	await database.guess.deleteMany({
		where: { channelId },
	});

	await database.trackmaniaMap.upsert({
		where: { channelId },
		update: map,
		create: {
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

export const deleteMap: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	// We use deleteMany here since delete would throw an error if the map doesn't exist
	await database.trackmaniaMap.deleteMany({
		where: { channelId },
	});

	// Delete all guesses for the channel when the map has changed
	await database.guess.deleteMany({
		where: { channelId },
	});
	response.status(200).json({ message: "Map deleted successfully." });
};

const PersonalBestSchema = z.object({
	time: z.number().min(0),
});

export const updatePersonalBest: RequestHandler = async (request, response) => {
	const { channelId, login } = request;
	if (!channelId || !login) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data, error } = PersonalBestSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	const { time } = data;
	wheelgpt.notifyNewPB(login, time);
	response.status(200).json({ message: "Personal best updated successfully." });
};

const RoomSchema = z.object({
	login: z.string().min(1).max(1000),
	name: z.string().min(1).max(1000),
	numberOfPlayers: z.number().min(0).max(2000),
	maxPlayers: z.number().min(0).max(2000),
});

export const updateRoom: RequestHandler = async (request, response) => {
	console.log(request.body);

	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data: room, error } = RoomSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	await database.trackmaniaRoom.upsert({
		where: { channelId },
		update: room,
		create: {
			channelId,
			...room,
		},
	});

	response.status(200).json({ message: "Room updated successfully." });
};

export const deleteRoom: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	await database.trackmaniaRoom.deleteMany({
		where: { channelId },
	});

	response.status(200).json({ message: "Room deleted successfully." });
};
