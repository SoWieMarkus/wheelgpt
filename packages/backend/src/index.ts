import app from "./app";
import { database } from "./database";
import { Twitch } from "./external";
import { env, logger } from "./utils";
import { wheelgpt } from "./wheelgpt";

const initialize = async () => {
	const channelIds = (await database.channel.findMany()).map((channel) => channel.id);

	logger.info("Updating channel information ...");

	// Update channel information on boot up
	const users = await Twitch.getUsers(channelIds);
	for (const user of users) {
		await database.channel.update({
			where: { id: user.id },
			data: {
				displayName: user.display_name,
				profileImage: user.profile_image_url,
				login: user.login,
				isLive: false, // Initial state, will be updated later
			},
		});
	}

	logger.info("Channel information updated successfully");

	// TODO update web hooks for channels

	console.log(await Twitch.getRegisteredWebhooks());

	const wheelgpt = await database.channel.findUnique({
		where: { login: "wheelgpt" },
	});
	console.log(wheelgpt);

	if (!wheelgpt) {
		logger.error("WheelGPT channel not found in the database. Please ensure it is registered.");
		throw new Error("WheelGPT channel not found");
	}

	await Twitch.registerWebhook("stream.offline", wheelgpt.id, "https://wheelgpt.de/api/twitch/webhook/stream");

	logger.info("Requesting stream status for channels...");

	// Get initial live status of channels. Later it will be updated by the Twitch webhook
	const streams = await Twitch.getStreams(channelIds);
	for (const stream of streams) {
		await database.channel.update({
			where: { id: stream.user_id },
			data: {
				isLive: stream.type === "live",
			},
		});
	}

	logger.info("Stream status updated successfully");
	logger.info("Initializing WheelGPT...");
	await wheelgpt.start();
	logger.info("WheelGPT initialized successfully");
};

logger.info("Starting the backend server...");
app.listen(env.PORT, () => {
	logger.info(`Running on port "${env.PORT}"`);
	initialize().catch((error) => {
		logger.error("Failed to initialize WheelGPT");
		console.error(error);
	});
});
