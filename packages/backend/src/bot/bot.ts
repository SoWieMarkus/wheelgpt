import { Counter } from "prom-client";
import { Client } from "tmi.js";
import { database } from "../database";
import { prometheus } from "../prometheus";
import { logger } from "../utils";
import { type ChannelConfig, TwitchChannel } from "./channel";
import { guessResultHandler } from "./commands";
import { getCommandArguments, getUser, TrackmaniaTime } from "./core";

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
	private _initialized = false;
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
		this._initialized = true;
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

	public async remove(login: string) {
		this.channelMap.delete(login);
		await this.part(login);
	}

	public getChannel(channelId: string): TwitchChannel | null {
		return this.channelMap.get(channelId) ?? null;
	}

	public async notifyNewPB(login: string, time: number) {
		const channel = this.getChannel(login);
		if (!channel) {
			logger.warn(`Channel ${login} not found for PB notification.`);
			return;
		}

		setTimeout(async () => {
			try {
				const result = new TrackmaniaTime(time);
				const message = await guessResultHandler(channel.id, result);
				if (message === null) return;
				this.say(login, message);
			} catch (error) {
				logger.error(`Error notifying new PB for channel ${login}:`, error);
			}
		}, channel.config.guessDelayTime * 1000);
	}

	public async reload(login: string, channelId: string) {
		this.channelMap.delete(login);
		const channel = await database.channel.findUnique({
			where: { id: channelId },
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
			logger.warn(`Channel ${login} not found for reload.`);
			return;
		}
		this.channelMap.set(channel.login, new TwitchChannel(channel.id, channel));
	}

	public get initialized(): boolean {
		return this._initialized;
	}
}
