import { Component, type OnInit, inject, model, signal } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatDividerModule } from "@angular/material/divider";
import { MatIconModule } from "@angular/material/icon";
import { MatSlideToggleModule } from "@angular/material/slide-toggle";
import { MatSliderModule } from "@angular/material/slider";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { Router } from "@angular/router";
import { TranslatePipe, TranslateService } from "@ngx-translate/core";
import { AuthenticationService } from "../../services/authentication.service";
import { BackendService } from "../../services/backend.service";
import { ProfileService } from "../../services/profile.service";
@Component({
	selector: "app-dashboard-page",
	imports: [
		MatSliderModule,
		MatDividerModule,
		TranslatePipe,
		MatButtonModule,
		MatIconModule,
		MatSlideToggleModule,
		FormsModule,
		MatTooltipModule,
	],
	templateUrl: "./dashboard-page.component.html",
	styleUrl: "./dashboard-page.component.scss",
})
export class DashboardPage implements OnInit {
	protected readonly profileService = inject(ProfileService);
	private readonly backendService = inject(BackendService);
	private readonly authenticationService = inject(AuthenticationService);
	private readonly router = inject(Router);
	private readonly snackbar = inject(MatSnackBar);
	private readonly translate = inject(TranslateService);

	protected readonly pluginToken = signal<string | null>(null);

	protected readonly settingBotActiveWhenOffline = model(false);
	protected readonly settingPublicChannel = model(true);
	protected readonly settingGuessDelay = model(2);

	public ngOnInit(): void {
		this.backendService.authentication
			.getPluginToken()
			.then((data) => {
				this.pluginToken.set(data.pluginToken);
			})
			.catch((error) => {
				const message = this.translate.instant("pages.dashboard.init.error");
				this.snackbar.open(message, "", {
					duration: 5000,
				});
				console.error("Failed to load plugin token:", error);
				this.pluginToken.set(null);
				this.authenticationService.removeToken();
				this.router.navigate(["/landing"]);
			});
	}

	protected logout(): void {
		this.authenticationService.removeToken();
		this.router.navigate(["/landing"]);
	}

	protected removeChannel(): void {
		this.backendService.authentication
			.remove()
			.then(() => {
				const message = this.translate.instant("pages.dashboard.remove.success");
				this.snackbar.open(message, "", {
					duration: 5000,
				});

				this.authenticationService.removeToken();
				this.router.navigate(["/landing"]);
			})
			.catch((error) => {
				const message = this.translate.instant("pages.dashboard.remove.error");
				this.snackbar.open(message, "", {
					duration: 5000,
				});
				console.error("Failed to remove channel:", error);
			});
	}

	protected renewPluginToken(): void {
		this.pluginToken.set(null);
		this.backendService.authentication
			.renewPluginToken()
			.then((data) => {
				const message = this.translate.instant("pages.dashboard.token.actions.renew.success");
				this.snackbar.open(message, "", {
					duration: 5000,
				});
				this.pluginToken.set(data.pluginToken);
			})
			.catch((error) => {
				const message = this.translate.instant("pages.dashboard.token.actions.renew.error");
				this.snackbar.open(message, "", {
					duration: 5000,
				});
				console.error("Failed to renew plugin token:", error);
			});
	}

	protected copyPluginToken(): void {
		const token = this.pluginToken();
		if (token === null) return;

		navigator.clipboard
			.writeText(token)
			.then(() => {
				const message = this.translate.instant("pages.dashboard.token.actions.copy.success");
				this.snackbar.open(message, "", {
					duration: 5000,
				});
			})
			.catch((error) => {
				const message = this.translate.instant("pages.dashboard.token.actions.copy.error");
				this.snackbar.open(message, "", {
					duration: 5000,
				});
				console.error("Failed to copy plugin token:", error);
			});
	}
}
