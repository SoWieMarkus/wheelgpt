import { Command, type User } from "../core";
import { Emote } from "../core/emotes";

export class EmotesCommand extends Command {
	protected async onExecute(user: User, _: string[]): Promise<string | null> {
		const emotes = Object.values(Emote)
			.map((emote) => emote.name)
			.join(" ");
		return `@${user.displayName} ${emotes} ${Emote.YEK.name}`;
	}
}
