import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { z } from "zod";
import { database } from "../database";
import { TMX } from "../external";
import { logger } from "../utils";
import { wheelgpt } from "../wheelgpt";

const TrackmaniaMapSchema = z
	.object({
		name: z.string().min(1).max(1000),
		uid: z.string().min(1).max(1000),
		author: z.string().min(1).max(1000),
		authorTime: z.number().min(0),
		goldTime: z.number().min(0),
		silverTime: z.number().min(0),
		bronzeTime: z.number().min(0),
		championTime: z.number().min(0),
	})
	.nullable();

export const updateMap: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data: map, error } = TrackmaniaMapSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}
	// Delete the map for the channel to store a new one
	// We use deleteMany here since delete would throw an error if the map doesn't exist
	await database.trackmaniaMap.delete({
		where: { channelId },
	});

	// Delete all guesses for the channel when the map has changed
	await database.guess.deleteMany({
		where: { channelId },
	});

	// New map is null, we already deleted the map and guesses
	// so we can just return a success message
	if (!map) {
		response.status(200).json({ message: "Map deleted successfully." });
		return;
	}

	// Upsert the map into the database (even though there shouldn't be a map for the channel at this point)
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

const PersonalBestSchema = z.object({
	time: z.number().min(0),
});

export const updatePersonalBest: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data, error } = PersonalBestSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	const { time } = data;
	wheelgpt.notifyNewPB(channelId, time);
	response.status(200).json({ message: "Personal best updated successfully." });
};

const RoomSchema = z
	.object({
		login: z.string().min(1).max(1000),
		name: z.string().min(1).max(1000),
		numberOfPlayers: z.number().min(0).max(2000),
		maxPlayers: z.number().min(0).max(2000),
	})
	.nullable();

export const updateRoom: RequestHandler = async (request, response) => {
	const channelId = request.channelId;
	if (!channelId) {
		throw createHttpError(401, "Authentication required.");
	}

	const { success, data: room, error } = RoomSchema.safeParse(request.body);
	if (!success) {
		throw createHttpError(400, error.errors[0].message);
	}

	await database.trackmaniaRoom.deleteMany({
		where: { channelId },
	});

	// If the room is null, we delete the room from the database
	// and return a success message
	if (!room) {
		response.status(200).json({ message: "Room deleted successfully." });
		return;
	}

	// Upsert the room into the database (even though there shouldn't be a room for the channel at this point)
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
