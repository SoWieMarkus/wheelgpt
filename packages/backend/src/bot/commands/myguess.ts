import { database } from "../../database";
import { Command, TrackmaniaTime, type User } from "../core";
import { Emote } from "../core/emotes";

export class MyGuessCommand extends Command {
	protected async onExecute(user: User, args: string[]): Promise<string | null> {
		const guess = await database.guess.findUnique({
			where: {
				channelId_userId: {
					userId: user.id,
					channelId: this.channelId,
				},
			},
		});

		if (guess === null) {
			return `@${user.displayName} I can't find any guess from you ${Emote.YEK.name}`;
		}
		const time = new TrackmaniaTime(guess.time);
		return `@${user.displayName} ${time.toString()}`;
	}
}
