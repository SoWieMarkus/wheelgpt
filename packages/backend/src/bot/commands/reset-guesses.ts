import { database } from "../../database";
import { Command, type User } from "../core";

export class ResetGuessesCommand extends Command {
	protected async onExecute(user: User, _: string[]): Promise<string | null> {
		await database.guess.deleteMany({
			where: {
				channelId: this.channelId,
			},
		});
		return `@${user.displayName} All guesses have been reset.`;
	}
}
