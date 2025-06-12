import {
    EmotesCommand,
    FormatCommand,
    GuessCommand,
    GuessResultCommand,
    MapCommand,
    MyGuessCommand,
} from "../commands";
import type { CommandArguments } from "./arguments";
import type { Command } from "./command";
import { AccessLevel } from "./user";

export type ChannelConfig = {
    channelId: string;
    displayName: string;
    guessDelayTime: number;
    botActiveWhenOffline: boolean;
    usagePublic: boolean;
}

export class TwitchChannel {
    private readonly commands: Command[];

    constructor(
        public readonly channelId: string,
        public readonly config: ChannelConfig,
    ) {
        this.commands = [
            new GuessCommand(channelId, {
                name: "guess",
                accessLevel: AccessLevel.USER,
                cooldown: 0,
                aliases: ["g"],
            }),
            new MyGuessCommand(channelId, {
                name: "myguess",
                accessLevel: AccessLevel.USER,
                cooldown: 0,
                aliases: ["myg"],
            }),
            new GuessResultCommand(channelId, {
                name: "result",
                accessLevel: AccessLevel.MOD,
                cooldown: 0,
                aliases: ["gr"],
            }),
            new MapCommand(channelId, {
                name: "map",
                accessLevel: AccessLevel.USER,
                cooldown: 10,
                aliases: ["m"],
            }),
            new EmotesCommand(channelId, {
                name: "wgpt-emotes",
                accessLevel: AccessLevel.MOD,
                cooldown: 10,
            }),
            new FormatCommand(channelId, {
                name: "format",
                accessLevel: AccessLevel.USER,
                cooldown: 5,
                aliases: ["f"],
            }),
        ];
    }

    public getCommand(commandArguments: CommandArguments) {
        for (const command of this.commands) {
            if (command.matches(commandArguments)) {
                return command;
            }
        }
        return null;
    }
}
