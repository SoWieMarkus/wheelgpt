import { Counter } from "prom-client";
import { Client } from "tmi.js";
import { database } from "../database";
import { prometheus } from "../prometheus";
import { logger } from "../utils";
import { guessResultHandler } from "./commands";
import { type ChannelConfig, getCommandArguments, getUser, TrackmaniaTime, TwitchChannel } from "./core";

export const failedConnectionAttemptsCounterMetric = new Counter({
	name: "wheelgpt_failed_connection_attempts_total",
	help: "Total failed connection attempts per channel",
	labelNames: ["login"],
	registers: [prometheus],
});

export const commandCounterMetric = new Counter({
	name: "wheelgpt_commands_total",
	help: "Total number of commands executed",
	labelNames: ["login", "commandName"],
	registers: [prometheus],
});

export class WheelGPT extends Client {
	private readonly channelMap: Map<string, TwitchChannel>;

	constructor(username: string, password: string) {
		super({
			identity: { username, password },
			options: { debug: false },
		});
		this.channelMap = new Map();
		this.on("message", async (tmiChannel, userstate, message, self) => {
			if (self) return; // Ignore messages from the bot itself

			const commandArguments = getCommandArguments(message);
			if (commandArguments === null) return;

			const login = tmiChannel.startsWith("#") ? tmiChannel.slice(1) : tmiChannel;
			const channel = this.getChannel(login);
			if (channel === null) return;

			const command = channel.getCommand(commandArguments);
			if (command === null) return;

			const user = getUser(userstate, tmiChannel);
			if (user === null) return;

			commandCounterMetric.inc({ login, commandName: command.name });
			command
				.execute(user, commandArguments.args)
				.then((result) => {
					if (result === null) return;
					this.say(login, result);
				})
				.catch((error) => {
					logger.error(`Error executing command ${command.name} in channel ${login}:`, error);
				});
		});
	}

	public async start() {
		await this.connect();
		const channels = await database.channel.findMany({
			select: {
				id: true,
				login: true,
				displayName: true,
				guessDelayTime: true,
				botActiveWhenOffline: true,
				usagePublic: true,
			},
		});
		for (const channel of channels) {
			await this.register(channel);
		}
	}

	public async register(channel: ChannelConfig) {
		try {
			const twitchChannel = new TwitchChannel(channel.id, channel);
			this.channelMap.set(channel.login, twitchChannel);
			await this.join(channel.login);
			logger.info(`Registered channel ${channel.id} (${channel.login})`);
		} catch (error) {
			failedConnectionAttemptsCounterMetric.inc({ login: channel.login });
			logger.error(`Failed to register channel ${channel.id}:`, error);
			console.error(error);
		}
	}

	public async remove(channelId: string) {
		this.channelMap.delete(channelId);
		await this.part(channelId);
	}

	public getChannel(channelId: string): TwitchChannel | null {
		return this.channelMap.get(channelId) ?? null;
	}

	public async notifyNewPB(channelId: string, time: number) {
		const channel = this.getChannel(channelId);
		if (!channel) {
			logger.warn(`Channel ${channelId} not found for PB notification.`);
			return;
		}

		setTimeout(async () => {
			const result = new TrackmaniaTime(time);
			const message = await guessResultHandler(channelId, result);
			if (message === null) return;
			this.say(channelId, message);
		}, channel.config.guessDelayTime * 1000);
	}

	public async reload(id: string) {
		this.channelMap.delete(id);
		const channel = await database.channel.findUnique({
			where: { id },
			select: {
				id: true,
				login: true,
				displayName: true,
				guessDelayTime: true,
				botActiveWhenOffline: true,
				usagePublic: true,
			},
		});
		if (!channel) {
			logger.warn(`Channel ${id} not found for reload.`);
			return;
		}
		this.channelMap.set(channel.login, new TwitchChannel(channel.id, channel));
	}
}
