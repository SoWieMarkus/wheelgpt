import { database } from "../../database";
import { Command } from "../core";

export class ResetGuessesCommand extends Command {
	protected async onExecute(): Promise<string | null> {
		await database.guess.deleteMany({
			where: {
				channelId: this.channelId,
			},
		});
		return null;
	}
}
