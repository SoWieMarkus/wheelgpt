import { database } from "../../database";
import type { User } from "../core";
import {
	buildGuessResultMessage,
	Command,
	EXAMPLE_FORMAT,
	evaluateGuesses,
	mentionUser,
	TrackmaniaMap,
	TrackmaniaTime,
	updateGuessResultLeaderboard,
} from "../core";

export const guessResultHandler = async (channelId: string, time: TrackmaniaTime) => {
	const guesses = await database.guess.findMany({
		where: {
			channelId,
		},
	});
	const map = await database.trackmaniaMap.findUnique({
		where: {
			channelId,
		},
	});
	const trackmaniaMap = map ? new TrackmaniaMap(map) : null;

	const winners = evaluateGuesses(guesses, time);
	database.guess.deleteMany({
		where: {
			channelId,
		},
	});

	updateGuessResultLeaderboard(winners, time);

	return buildGuessResultMessage(trackmaniaMap, time, winners);
};

export class GuessResultCommand extends Command {
	protected async onExecute(user: User, args: string[]): Promise<string | null> {
		if (args.length === 0) {
			return `${mentionUser(user.displayName)} ${EXAMPLE_FORMAT}`;
		}
		const time = TrackmaniaTime.parse(args[0]);
		if (time === null) {
			return `${mentionUser(user.displayName)} smh granadyy mods are all degens. Wrong format you idiot. ${EXAMPLE_FORMAT}`;
		}

		return guessResultHandler(this.channelId, time);
	}
}
