import { HttpClient } from "@angular/common/http";
import { Injectable, inject } from "@angular/core";
import type { Schema } from "@wheelgpt/shared";
import type { Observable } from "rxjs";
import type { z } from "zod";
import { environment } from "../../environments/environment";
import { WebTokenSchema } from "../schemas/authentication";
import { ChannelSchema } from "../schemas/channel";
import { PublicChannelsSchema } from "../schemas/landing";

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
			updateSettings: (settings: z.infer<typeof Schema.channel.settings>) =>
				this.post("channel/settings", settings, ChannelSchema),
		};
	}

	public get authentication() {
		return {
			login: (code: string) => this.post("authentication/login", { code }, WebTokenSchema),
		};
	}

	public get landing() {
		return {
			publicChannels: () => this.get("landing/channels", PublicChannelsSchema),
		};
	}
}
