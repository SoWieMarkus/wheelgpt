import { getLeaderboard, getLeaderboardByName } from "@prisma/client/sql";
import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { z } from "zod";
import { database } from "../database";

const LEADERBOARD_PAGE_SIZE = 100;

export const getLeaderboardPage: RequestHandler = async (request, response) => {
	const channelId = request.params.channelId;
	if (!channelId) {
		throw createHttpError(400, "Creator ID is required.");
	}
	const QuerySchema = z.object({
		page: z.coerce.number().min(1, "Page must be a positive integer").default(1),
	});

	const query = QuerySchema.safeParse(request.query);
	if (!query.success) {
		throw createHttpError(400, `Invalid query parameters: ${z.prettifyError(query.error)}`);
	}

	const leaderboard = await database.$queryRawTyped(
		getLeaderboard(channelId, LEADERBOARD_PAGE_SIZE, (query.data.page - 1) * LEADERBOARD_PAGE_SIZE),
	);
	response.status(200).json(leaderboard);
};

export const getLeaderboardPositionByName: RequestHandler = async (request, response) => {
	const { channelId, displayName } = request.params;
	if (!channelId || !displayName) {
		throw createHttpError(400, "Channel ID and display name are required.");
	}
	const leaderboard = await database.$queryRawTyped(getLeaderboardByName(channelId, displayName));
	response.status(200).json(leaderboard.length > 0 ? leaderboard[0] : null);
};
