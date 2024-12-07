import { database } from "../database";
import { Command, isCommandOnCooldown, matchesCommand } from "./command";
import { CommandArguments } from "./command-arguments";
import { buildGuessResultMessage, evaluateGuesses, Guess } from "./guess/guess";
import { parseTrackmaniaTime, TrackmaniaTime } from "./guess/time";
import { TrackmaniaMap } from "./map/map";
import { AccessLevel, mentionUser, User } from "./user";

const EXAMPLE_FORMAT = "hh:mm:ss.xxx (hours and minutes are optional)";
const MAX_DELAY = 60;
const MIN_DELAY = 0;

export const WheelGPTUser: User = {
    id: "-1",
    name: "wheelgpt",
    displayName: "WheelGPT",
    accessLevel: AccessLevel.MOD,
    channelId: "#wheelgpt",
};

export class Channel {

    private readonly guesses: Map<string, Guess>;
    private currentMap: TrackmaniaMap | null;

    private readonly commands: Command[] = [
        {
            name: "map",
            cooldown: 10,
            task: (user) => {
                if (this.currentMap === null) {
                    return mentionUser(user, "No Map");
                }
                return mentionUser(user, this.currentMap.toString())
            },
            accessLevel: AccessLevel.USER
        },
        {
            name: "format",
            cooldown: 10,
            task: (user, args) => {
                if (args.length === 0) {
                    return mentionUser(user, EXAMPLE_FORMAT)
                }
                return `@${args[0]} ${EXAMPLE_FORMAT}`
            },
            accessLevel: AccessLevel.USER
        },
        {
            name: "guess",
            aliases: ["g"],
            task: (user, args) => {
                if (args.length === 0) return null;
                const time = args[0];
                const guess = parseTrackmaniaTime(time);
                if (guess === null) return null;

                this.guesses.set(user.id, {
                    user,
                    time: guess,
                });
                return null;
            },
            accessLevel: AccessLevel.USER,
        },
        {
            name: "wgpt-time-delay",
            accessLevel: AccessLevel.MOD,
            task: (user, args) => {
                if (args.length === 0) return null;

                const delay = Number.parseInt(args[0]);
                if (Number.isNaN(delay) || delay < MIN_DELAY || delay > MAX_DELAY) {
                    return mentionUser(user, `Please provide a delay time (in seconds, max ${MAX_DELAY}s) as a number.`);
                }

                this.guessResultDelay = delay;
                database.channel.update({
                    where: { channelId: this.id },
                    data: {
                        guessDelayTime: delay
                    }
                });
                return mentionUser(user, "Successfully updated delay time.");
            },
        },
        {
            name: "result",
            accessLevel: AccessLevel.MOD,
            task: (user, args) => {
                if (args.length === 0) {
                    return mentionUser(user, "Please provide a result time.");
                }

                const time = args[0];
                const result = parseTrackmaniaTime(time);
                if (result === null) {
                    return mentionUser(user, `smh granadyy mods are all degens. Wrong format you idiot. ${EXAMPLE_FORMAT}`);
                }
                return this.guessResult(result);
            },
        },
        {
            name: "myguess",
            aliases: ["mg"],
            cooldown: 2,
            accessLevel: AccessLevel.USER,
            task: (user, args) => {
                const guess = this.guesses.get(user.id);
                if (guess === undefined) {
                    return mentionUser(user, "I can't find any guess from you YEK");
                }
                return mentionUser(user, guess.time.toString());
            },
        },
        {
            name: "resetTimes",
            accessLevel: AccessLevel.MOD,
            task: () => {
                this.guesses.clear();
                return null;
            }
        },
    ]

    constructor(private readonly id: string, private guessResultDelay: number) {
        this.guesses = new Map();
        this.currentMap = null;
    }


    private getCommand(commandArguments: CommandArguments) {
        for (const command of this.commands) {
            if (matchesCommand(command, commandArguments)) {
                return command;
            }
        }
        return null;
    }

    public execute(user: User, commandArguments: CommandArguments): string | null {
        const command = this.getCommand(commandArguments);
        if (
            command === null ||
            user.accessLevel < command.accessLevel ||
            isCommandOnCooldown(command)
        )
            return null;

        command.lastUsed = Date.now();
        return command.task(user, commandArguments.args);
    }

    public setMap(map: TrackmaniaMap | null) {
        this.currentMap = map;
    }

    public getGuessDelay() {
        return this.guessResultDelay * 1000;
    }

    public guessResult(time: TrackmaniaTime) {
        const winners = evaluateGuesses(this.guesses, time);
        this.guesses.clear();
        return buildGuessResultMessage(this.currentMap, time, winners);
    }

}