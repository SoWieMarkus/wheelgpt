import { database } from "../../database";
import { Command, mentionUser, TrackmaniaMap, type User } from "../core";

export class MapCommand extends Command {
	protected async onExecute(user: User, _: string[]): Promise<string | null> {
		const map = await database.trackmaniaMap.findUnique({
			where: {
				channelId: this.channelId,
			},
		});
		if (map === null) {
			return `${mentionUser(user.displayName)} No map is currently set for this channel.`;
		}

		const trackmaniaMap = new TrackmaniaMap(map);
		return `${mentionUser(user.displayName)} ${trackmaniaMap.toString()}`;
	}
}
