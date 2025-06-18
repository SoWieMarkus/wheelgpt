import { Pipe, type PipeTransform } from "@angular/core";

const EMOTES: Record<string, string> = {
	AUTHOR_MEDAL: "emotes/author_medal.avif",
	YEK: "emotes/yek.avif",
	DINK_DONK: "emotes/dink_donk.avif",
	OK: "emotes/ok.avif",
	BWOAH: "emotes/bwoah.avif",
	GIGACHAD: "emotes/gigachad.avif",
};

@Pipe({
	name: "twitch",
})
export class TwitchMessagePipe implements PipeTransform {
	public transform(value: unknown, ..._: unknown[]): unknown {
		// Find all emotes inside the message
		// They are wrapped by two $$
		// Example: This message contains an $$EMOTE$$
		if (typeof value !== "string") return value;

		// Replace all $$EMOTE$$ with <img ...>
		return value.replace(/\$\$(\w+)\$\$/g, (match, emote) => {
			const url = EMOTES[emote];
			if (url) {
				return `<img src="${url}" alt="${emote}" class="emote" />`;
			}
			return match; // leave unchanged if not found
		});
	}
}
