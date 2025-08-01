import { Injectable, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { environment } from "../../environments/environment";
import { ProfileService } from "./profile.service";

@Injectable({
	providedIn: "root",
})
export class AuthenticationService {
	private readonly router = inject(Router);
	private readonly profileService = inject(ProfileService);

	// Indicates if the service is currently initializing the profile data
	public readonly isFetchingData = signal(false);

	// Indicates if the user is currently in the process of logging in via Twitch
	public readonly twitchLoginProcess = signal<boolean>(false);

	public getToken(): string | null {
		return localStorage.getItem("token");
	}

	public async setToken(token: string) {
		localStorage.setItem("token", token);
		this.initialize();
	}

	public async initialize() {
		this.isFetchingData.set(true);
		try {
			await this.profileService.initialize();
		} catch (error) {
			console.error("Error during initialization:", error);
			this.removeToken();
			this.router.navigate(["/login"]);
		} finally {
			this.isFetchingData.set(false);
		}
	}

	public removeToken(): void {
		localStorage.removeItem("token");
		this.profileService.reset();
	}

	public isTokenExpired(): boolean {
		const token = this.getToken();
		if (!token) {
			return true;
		}

		const expiry = this.getTokenExpiry(token);
		return expiry ? expiry < Date.now() / 1000 : false;
	}

	private getTokenExpiry(token: string): number | null {
		try {
			const decodedToken = JSON.parse(atob(token.split(".")[1]));
			return decodedToken.exp ?? null;
		} catch (e) {
			console.error(e);
			return null;
		}
	}

	public get twitchSSOUrl(): string {
		return `${environment.backend}/authentication/twitch`;
	}
}
