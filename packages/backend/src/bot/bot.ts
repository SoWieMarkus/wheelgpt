import { Counter } from "prom-client";
import { Client } from "tmi.js";
import { database } from "../database";
import { prometheus } from "../prometheus";
import { logger } from "../utils";
import { type ChannelConfig, getCommandArguments, getUser, TrackmaniaTime, TwitchChannel } from "./core";
import { guessResultHandler } from "./commands";


export const failedConnectionAttemptsCounterMetric = new Counter({
    name: "failed_connection_attempts_total",
    help: "Total failed connection attempts per channel",
    labelNames: ["channelId"],
    registers: [prometheus],
});

export const commandCounterMetric = new Counter({
    name: "commands_total",
    help: "Total number of commands executed",
    labelNames: ["channelId", "commandName"],
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
        this.on("message", async (channelId, userstate, message, self) => {
            const commandArguments = getCommandArguments(message);
            if (commandArguments === null) return;

            const channel = this.getChannel(channelId);
            if (channel === null) return;

            const command = channel.getCommand(commandArguments);
            if (command === null) return;

            const user = getUser(userstate, channelId);
            if (user === null) return;

            commandCounterMetric.inc({ channelId, commandName: command.name });
            command
                .execute(user, commandArguments.args)
                .then((result) => {
                    if (result === null) return;
                    this.say(channelId, result);
                })
                .catch((error) => {
                    logger.error(`Error executing command ${command.name} in channel ${channelId}:`, error);
                });
        });
    }

    public async start() {
        await this.connect();
        const channels = await database.channel.findMany({
            select: {
                channelId: true,
                displayName: true,
                guessDelayTime: true,
                botActiveWhenOffline: true,
                usagePublic: true,
            },
        });
        for (const channel of channels) {
            this.register(channel);
        }
    }

    public async register(channel: ChannelConfig) {
        try {
            const twitchChannel = new TwitchChannel(channel.channelId, channel);
            this.channelMap.set(channel.channelId, twitchChannel);
            await this.join(channel.channelId);
        } catch (error) {
            failedConnectionAttemptsCounterMetric.inc({ channelId: channel.channelId });
            logger.error(`Failed to register channel ${channel.channelId}:`, error);
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

    public async reload(channelId: string) {
        this.channelMap.delete(channelId);
        const channel = await database.channel.findUnique({
            where: { channelId },
            select: {
                channelId: true,
                displayName: true,
                guessDelayTime: true,
                botActiveWhenOffline: true,
                usagePublic: true,
            },
        });
        if (!channel) {
            logger.warn(`Channel ${channelId} not found for reload.`);
            return;
        }
        this.channelMap.set(channel.channelId, new TwitchChannel(channel.channelId, channel));
    }
}
