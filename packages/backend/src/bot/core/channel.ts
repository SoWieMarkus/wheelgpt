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
	id: string;
	login: string;
	displayName: string;
	guessDelayTime: number;
	botActiveWhenOffline: boolean;
	usagePublic: boolean;
};

export class TwitchChannel {
	private readonly commands: Command[];

	constructor(
		public readonly id: string,
		public readonly config: ChannelConfig,
	) {
		this.commands = [
			new GuessCommand(id, {
				name: "guess",
				accessLevel: AccessLevel.USER,
				cooldown: 0,
				aliases: ["g"],
			}),
			new MyGuessCommand(id, {
				name: "myguess",
				accessLevel: AccessLevel.USER,
				cooldown: 0,
				aliases: ["mg"],
			}),
			new GuessResultCommand(id, {
				name: "result",
				accessLevel: AccessLevel.MOD,
				cooldown: 0,
			}),
			new MapCommand(id, {
				name: "map",
				accessLevel: AccessLevel.USER,
				cooldown: 10,
			}),
			new EmotesCommand(id, {
				name: "wgpt-emotes",
				accessLevel: AccessLevel.MOD,
				cooldown: 10,
			}),
			new FormatCommand(id, {
				name: "format",
				accessLevel: AccessLevel.USER,
				cooldown: 5,
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
