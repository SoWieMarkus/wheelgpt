import type { CommandArguments } from "./arguments";
import { AccessLevel, type User } from "./user";

export type CommandConfig = {
	name: string;
	cooldown?: number;
	aliases?: string[];
	accessLevel?: number;
};

export abstract class Command {
	public readonly name: string;
	protected readonly aliases: string[];

	// The cooldown in seconds for this command
	protected readonly cooldown: number;
	protected readonly accessLevel: number;
	protected readonly channelId: string;

	// The last time this command was used, in milliseconds since the epoch
	protected lastUsed?: number;

	constructor(channelId: string, config: CommandConfig) {
		this.cooldown = config.cooldown ?? 0;
		this.aliases = config.aliases ?? [];
		this.accessLevel = config.accessLevel ?? AccessLevel.USER;
		this.channelId = channelId;
		this.name = config.name;
	}

	public matches(commandArguments: CommandArguments): boolean {
		const key = commandArguments.key;
		return this.name === key || this.aliases.includes(key);
	}

	public isOnCooldown(): boolean {
		const now = Date.now();
		if (this.lastUsed === undefined) {
			return false;
		}
		return now - this.cooldown * 1000 < this.lastUsed;
	}

	public async execute(user: User, args: string[]) {
		// Check if the user has sufficient permissions to execute this command
		if (user.accessLevel < this.accessLevel) {
			return null;
		}

		// Check if the command is on cooldown
		if (this.isOnCooldown()) {
			return null;
		}

		// Update the last used timestamp to now to prevent re-using the command too quickly
		this.lastUsed = Date.now();
		return await this.onExecute(user, args); // Execute the command logic
	}

	protected abstract onExecute(user: User, args: string[]): Promise<string | null>;
}
