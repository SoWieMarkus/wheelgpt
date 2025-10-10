import { HttpClient } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import type { Observable } from "rxjs";
import type * as z from "zod";
import { environment } from "../../environments/environment";
import { PluginTokenSchema, WebTokenSchema } from "../schemas/authentication";
import { ChannelSchema, ChannelSettingsSchema, LandingPageChannelsSchema, PublicChannelDetailsSchema } from "../schemas/channel";
import { LeaderboardEntrySchema } from "../schemas/leaderboard";

const BACKEND_URL = environment.backend;

@Injectable({
	providedIn: "root",
})
export class BackendService {
	private readonly http = inject(HttpClient);

	private handleResponse<T>(observable: Observable<T>, responseModel?: z.Schema<T>): Promise<T> {
		return new Promise((resolve, reject) => {
			observable.subscribe({
				next: (response) => {
					if (!responseModel) {
						resolve(response);
						return;
					}
					const { success, data, error } = responseModel.safeParse(response);
					if (!success) {
						reject(new Error(error.message));
						return;
					}
					resolve(data);
				},
				error: (error) => {
					console.error("Error in backend service:", error);
					reject(new Error(error.error?.error ?? "An error occurred while processing the request."));
				},
			});
		});
	}

	private get<T>(url: string, responseModel?: z.Schema<T>): Promise<T> {
		return this.handleResponse(this.http.get<T>(`${BACKEND_URL}/${url}`), responseModel);
	}

	private post<T>(url: string, body: unknown, responseModel?: z.Schema<T>): Promise<T> {
		return this.handleResponse(this.http.post<T>(`${BACKEND_URL}/${url}`, body), responseModel);
	}

	private delete<T>(url: string, responseModel?: z.Schema<T>): Promise<T> {
		return this.handleResponse(this.http.delete<T>(`${BACKEND_URL}/${url}`), responseModel);
	}

	public get channel() {
		return {
			me: () => this.get("channel/me", ChannelSchema),
			updateSettings: (settings: z.infer<typeof ChannelSettingsSchema>) =>
				this.post("channel/settings", settings, ChannelSchema),
			publicChannels: () => this.get("channel", LandingPageChannelsSchema),
			getChannelById: (channelId: string) => this.get(`channel/${channelId}`, PublicChannelDetailsSchema),
		};
	}

	public get authentication() {
		return {
			login: (code: string) => this.post("authentication/login", { code }, WebTokenSchema),
			renewPluginToken: () => this.post("authentication/renew", {}, PluginTokenSchema),
			remove: () => this.delete("authentication/remove"),
			getPluginToken: () => this.get("authentication/token", PluginTokenSchema),
		};
	}

	public get leaderboard() {
		return {
			get: (channelId: string, page = 1) =>
				this.get(`leaderboard/${channelId}?page=${page}`, LeaderboardEntrySchema.array()),
			getByName: (channelId: string, username: string) =>
				this.get(`leaderboard/${channelId}/user/${username}`, LeaderboardEntrySchema),
		};
	}
}
