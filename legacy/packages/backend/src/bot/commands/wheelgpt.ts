import { Command } from "../core";
import { Emote } from "../core/emotes";

export class WheelGPTCommand extends Command {
	protected async onExecute(): Promise<string | null> {
		return `${Emote.HEY_GUYS.name} I am a Trackmania Twitch Bot! You can ask me about the current map, the current room  and guess the next best time. More details: https://wheelgpt.dev | ⭐ Star me on Github: https://github.com/SoWieMarkus/wheelgpt`;
	}
}
