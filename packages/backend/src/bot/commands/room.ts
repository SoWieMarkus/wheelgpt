import { database } from "../../database";
import { Command, type User } from "../core";

export class RoomCommand extends Command {
    protected async onExecute(user: User, _: string[]): Promise<string | null> {
        const room = await database.trackmaniaRoom.findUnique({
            where: {
                channelId: this.channelId,
            },
        });
        if (room === null) {
            return `@${user.displayName} I am currently not in a room.`;
        }

        const message = `${room.name} [${room.numberOfPlayers}/${room.maxPlayers}]`;
        return `@${user.displayName} ${message}`;
    }
}
