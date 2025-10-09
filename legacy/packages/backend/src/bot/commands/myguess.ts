import { database } from "../../database";
import { Command, mentionUser, TrackmaniaTime, type User } from "../core";
import { Emote } from "../core/emotes";

export class MyGuessCommand extends Command {
	protected async onExecute(user: User, _: string[]): Promise<string | null> {
		const guess = await database.guess.findUnique({
			where: {
				channelId_userId: {
					userId: user.id,
					channelId: this.channelId,
				},
			},
		});

		if (guess === null) {
			return `${mentionUser(user.displayName)} I can't find any guess from you ${Emote.YEK.name}`;
		}
		const time = new TrackmaniaTime(guess.time);
		return `${mentionUser(user.displayName)} ${time.toString()}`;
	}
}
