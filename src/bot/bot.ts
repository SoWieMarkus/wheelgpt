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

            const commandArguments = getCommandArguments(message);
            if (commandArguments === null) return;

            const user = getUser(chatUser, channelId);
            if (user === null) return;

            const channel = this.getChannel(channelId);
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
            Log.info(`Joining channel "${channel.channelId}".`)
            this.channelMap.set(channel.channelId, new Channel(channel.channelId, channel.guessDelayTime));
            await this.join(channel.channelId);
            await this.say(channel.channelId, "Hello World HeyGuys")
        }
        return result;
    }

    public async register(user: { login: string, display_name: string, profile_image_url: string }) {
        const token = uuid.v4();
        const { login: channelId, display_name: displayName, profile_image_url: profileImage } = user;
        const channel = await database.channel.create({
            data: { token, channelId, displayName, profileImage }
        });
        this.channelMap.set(channelId, new Channel(channel.channelId, channel.guessDelayTime));
        await this.join(channelId);
        await this.say(channelId, "Hello World HeyGuys")
        Log.info(`Joining channel "${channelId}".`)
        return channel;
    }

    public async remove(channelId: string) {
        await database.channel.delete({ where: { channelId } });
        this.channelMap.delete(channelId);
        await this.part(channelId);
        Log.info(`Removing channel "${channelId}".`);
    }

    public getChannel(channelId: string) {
        return this.channelMap.get(channelId.replace("#", ""));
    }

}

export const WheelGPT = new WheelGPTBot();




