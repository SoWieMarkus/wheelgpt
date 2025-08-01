import { Injectable, inject, signal } from "@angular/core";
import type { Channel } from "../schemas/channel";
import { BackendService } from "./backend.service";

@Injectable({
	providedIn: "root",
})
export class ProfileService {
	private readonly backend = inject(BackendService);
	public readonly channel = signal<Channel | null>(null);

	public async initialize() {
		const channel = await this.backend.channel.me();
		this.channel.set(channel);
	}

	public reset() {
		this.channel.set(null);
	}
}
