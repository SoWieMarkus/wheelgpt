import { Client } from "tmi.js";
import { env, Log } from "../utils";
import { getCommandArguments } from "./command-arguments";
import { getUser } from "./user";
import { Channel } from "./channel";
import { database } from "../database";
import * as uuid from "uuid";

const username = env.BOT_USERNAME;
const password = env.BOT_OAUTH_TOKEN;


class WheelGPTBot extends Client {

    private readonly channelMap: Map<string, Channel>;

    constructor() {
        super({
            identity: { username, password },
            options: { debug: false }
        });
        this.channelMap = new Map();
        this.on("message", (channelId, chatUser, message, self) => {
            if (self) return;

            const commandArguments = getCommandArguments(message);
            if (commandArguments === null) return;

            const user = getUser(chatUser, channelId);
            if (user === null) return;

            const channel = this.channelMap.get(channelId);
            if (channel === undefined) return;

            const response = channel.execute(user, commandArguments);
            if (response === null) return;

            this.say(channelId, response);
        })
    }

    public override async connect(): Promise<[string, number]> {
        const result = await super.connect();
        const channels = await database.channel.findMany();

        for (const channel of channels) {
            this.channelMap.set(channel.channelId, new Channel(channel.channelId, channel.guessDelayTime));
            this.join(channel.channelId);
        }
        return result;
    }

    public async register(channelId: string) {
        const token = uuid.v4();
        const channel = await database.channel.create({
            data: { token, channelId }
        });
        this.channelMap.set(channelId, new Channel(channel.channelId, channel.guessDelayTime));
        this.join(channelId);
        Log.info(`Joining channel "${channelId}".`)
        return channel;
    }

    public async remove(channelId: string) {
        await database.channel.delete({ where: { channelId } });
        this.channelMap.delete(channelId);
        this.part(channelId);
        Log.info(`Removing channel "${channelId}".`);
    }

    public getChannel(channelId: string) {
        return this.channelMap.get(channelId);
    }

}

export const WheelGPT = new WheelGPTBot();




