import { Injectable, inject } from "@angular/core";
import { Router } from "@angular/router";
import { environment } from "../../environments/environment";

@Injectable({
	providedIn: "root",
})
export class AuthenticationService {
	private readonly router = inject(Router);

	public getToken(): string | null {
		return localStorage.getItem("token");
	}

	public async setToken(token: string) {
		localStorage.setItem("token", token);
	}

	public removeToken(): void {
		localStorage.removeItem("token");
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
