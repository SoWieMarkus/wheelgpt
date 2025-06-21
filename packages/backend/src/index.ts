import app from "./app";
import { database } from "./database";
import { Twitch } from "./external";
import { env, logger } from "./utils";
import { wheelgpt } from "./wheelgpt";

// Update channel information on boot up
const updateChannelDetails = async (channelIds: string[]) => {
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
};

// Get initial live status of channels. Later it will be updated by the Twitch webhook
const updateStreamStatus = async (channelIds: string[]) => {
	const streams = await Twitch.getStreams(channelIds);
	for (const stream of streams) {
		await database.channel.update({
			where: { id: stream.user_id },
			data: {
				isLive: stream.type === "live",
			},
		});
	}
};

const initialize = async () => {
	const channelIds = (await database.channel.findMany()).map((channel) => channel.id);

	logger.info("Updating channel information ...");
	await updateChannelDetails(channelIds);
	logger.info("Channel information updated successfully");

	logger.info("Registering Twitch webhook subscriptions...");
	await Twitch.syncWebhooks(channelIds, env.TWITCH_STREAM_WEB_HOOK_URL);
	logger.info("Twitch webhook subscriptions registered successfully");
	console.log(await Twitch.getRegisteredWebhooks());

	logger.info("Requesting stream status for channels...");
	await updateStreamStatus(channelIds);
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
