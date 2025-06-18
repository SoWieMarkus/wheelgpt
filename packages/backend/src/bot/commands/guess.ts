import { database } from "../../database";
import { Command, type User } from "../core";
import { EXAMPLE_FORMAT, TrackmaniaTime } from "../core/time";

export class GuessCommand extends Command {
	protected async onExecute(user: User, args: string[]): Promise<string | null> {
		if (args.length === 0) {
			return `@${user.displayName} ${EXAMPLE_FORMAT}`;
		}

		const time = TrackmaniaTime.parse(args[0]);
		if (time === null) {
			return `@${user.displayName} ${EXAMPLE_FORMAT}`;
		}

		await database.guess.upsert({
			where: {
				channelId_userId: {
					userId: user.id,
					channelId: this.channelId,
				},
			},
			create: {
				userId: user.id,
				displayName: user.displayName,
				channelId: this.channelId,
				time: time.getTotalInMilliSeconds(),
			},
			update: {
				time: time.getTotalInMilliSeconds(),
			},
		});
		return null;
	}
}
