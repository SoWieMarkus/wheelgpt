import { Component, computed, input } from "@angular/core";
import { TranslatePipe } from "@ngx-translate/core";
import { TwitchMessagePipe } from "../../pipes/twitch-message.pipe";
import type { Command } from "../../services/commands.service";

@Component({
	selector: "app-command",
	imports: [TwitchMessagePipe, TranslatePipe],
	templateUrl: "./command.component.html",
	styleUrl: "./command.component.scss",
})
export class CommandComponent {
	public readonly command = input.required<Command>();
	protected readonly commandAliases = computed(() => {
		const command = this.command();
		return [command.name, ...command.aliases];
	});

	protected get showUserBadge() {
		return this.command().accessLevel !== "User";
	}

	protected get userBadge() {
		switch (this.command().accessLevel) {
			case "User":
				return "";
			case "Mod":
				return "badges/mod.png";
			case "Streamer":
				return "badges/streamer.png";
			case "VIP":
				return "badges/vip.png";
			case "Subscriber":
				return "badges/subscriber.png";
			default:
				return "";
		}
	}
}
