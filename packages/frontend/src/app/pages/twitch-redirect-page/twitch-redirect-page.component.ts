import { Component, inject, type OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ActivatedRoute, Router } from "@angular/router";
import { TranslatePipe, TranslateService } from "@ngx-translate/core";
import { LoadingComponent } from "../../components/loading/loading.component";
import { AuthenticationService } from "../../services/authentication.service";
import { BackendService } from "../../services/backend.service";

@Component({
	selector: "app-twitch-redirect-page",
	imports: [TranslatePipe, LoadingComponent],
	templateUrl: "./twitch-redirect-page.component.html",
	styleUrl: "./twitch-redirect-page.component.scss",
})
export class TwitchRedirectPage implements OnInit {
	private readonly route = inject(ActivatedRoute);
	private readonly router = inject(Router);
	private readonly backendService = inject(BackendService);
	private readonly authenticationService = inject(AuthenticationService);
	private readonly translate = inject(TranslateService);

	private readonly snackbar = inject(MatSnackBar);

	public ngOnInit(): void {
		const code = this.route.snapshot.queryParamMap.get("code");

		if (!code) {
			console.error("Missing code or state in Twitch redirect.");
			this.router.navigate(["/"]);
			return;
		}

		this.authenticationService.twitchLoginProcess.set(true);
		this.backendService.authentication
			.login(code)
			.then((data) => {
				this.authenticationService.setToken(data.webToken);
			})
			.catch((error) => {
				const message = this.translate.instant("pages.redirect.login-failed.message");
				const close = this.translate.instant("pages.redirect.login-failed.close");
				console.error("Login failed:", error);
				this.snackbar.open(message, close, {
					duration: 3000,
				});
			})
			.finally(() => {
				this.router.navigate(["/"]);
				this.authenticationService.twitchLoginProcess.set(false);
			});
	}
}
