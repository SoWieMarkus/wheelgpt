import { database } from "../../database";
import { Command, mentionUser, type User } from "../core";

export class ResetGuessesCommand extends Command {
	protected async onExecute(user: User, _: string[]): Promise<string | null> {
		await database.guess.deleteMany({
			where: {
				channelId: this.channelId,
			},
		});
		return `${mentionUser(user.displayName)} All guesses have been reset.`;
	}
}
