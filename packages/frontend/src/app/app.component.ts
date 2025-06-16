import { Component, inject } from "@angular/core";
import { Router } from "@angular/router";
import { AuthenticationService } from "./services/authentication.service";

@Component({
	selector: "app-root",
	imports: [],
	templateUrl: "./app.component.html",
	styleUrl: "./app.component.scss",
})
export class AppComponent {
	protected readonly authenticationService = inject(AuthenticationService);
	protected readonly router = inject(Router);

	protected readonly year = new Date().getFullYear();

	protected logout(): void {
		this.authenticationService.removeToken();
		this.router.navigate(["/landing"]);
	}

	public ngOnInit(): void {
		if (this.authenticationService.isTokenExpired()) {
			this.authenticationService.removeToken();
			return;
		}
		this.authenticationService.initialize();
	}
}
