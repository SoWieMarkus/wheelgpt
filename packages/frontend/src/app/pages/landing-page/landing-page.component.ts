import { Component, inject, signal } from "@angular/core";
import { MatBadgeModule } from "@angular/material/badge";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";
import { RouterLink } from "@angular/router";
import { TranslatePipe } from "@ngx-translate/core";
import { CommandComponent } from "../../components/command/command.component";
import type { LandingPageChannels } from "../../schemas/channel";
import { AuthenticationService } from "../../services/authentication.service";
import { BackendService } from "../../services/backend.service";
import { CommandsService } from "../../services/commands.service";
import { ProfileService } from "../../services/profile.service";
@Component({
	selector: "app-landing-page",
	imports: [CommandComponent, MatIconModule, TranslatePipe, RouterLink, MatChipsModule, MatBadgeModule],
	templateUrl: "./landing-page.component.html",
	styleUrl: "./landing-page.component.scss",
})
export class LandingPage {
	private readonly backendService = inject(BackendService);
	protected readonly authenticationService = inject(AuthenticationService);
	protected readonly profileService = inject(ProfileService);
	protected readonly commandsService = inject(CommandsService);

	protected readonly data = signal<LandingPageChannels>([]);

	public ngOnInit() {
		this.backendService.channel
			.publicChannels()
			.then((data) => {
				// sort channels by live status and display name
				data.sort((a, b) => {
					if (a.isLive && !b.isLive) return -1; // a is live, b is not
					if (!a.isLive && b.isLive) return 1; // b is live, a is not
					// both are live or both are not live, sort by display name
					return a.displayName.localeCompare(b.displayName);
				});
				this.data.set(data);
			})
			.catch((error) => {
				console.error("Failed to load landing data:", error);
				this.data.set([]);
			});
	}

	public openTwitch(login: string) {
		const url = `https://www.twitch.tv/${login}`;
		window.open(url, "_blank");
	}
}
