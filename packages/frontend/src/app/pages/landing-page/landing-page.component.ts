import { Component, inject, signal } from "@angular/core";
import { MatIconModule } from "@angular/material/icon";
import { TranslatePipe } from "@ngx-translate/core";
import { CommandComponent } from "../../components/command/command.component";
import type { PublicChannels } from "../../schemas/landing";
import { AuthenticationService } from "../../services/authentication.service";
import { BackendService } from "../../services/backend.service";
import { CommandsService } from "../../services/commands.service";

@Component({
	selector: "app-landing-page",
	imports: [CommandComponent, MatIconModule, TranslatePipe],
	templateUrl: "./landing-page.component.html",
	styleUrl: "./landing-page.component.scss",
})
export class LandingPage {
	private readonly backendService = inject(BackendService);
	protected readonly authenticationService = inject(AuthenticationService);

	protected readonly commandsService = inject(CommandsService);

	protected readonly data = signal<PublicChannels | null>(null);

	public ngOnInit() {
		this.backendService.landing
			.publicChannels()
			.then((data) => {
				console.log("Landing data loaded:", data);
				this.data.set(data);
			})
			.catch((error) => {
				console.error("Failed to load landing data:", error);
				this.data.set(null);
			});
	}
}
