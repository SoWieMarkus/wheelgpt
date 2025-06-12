import { Command, EXAMPLE_FORMAT, type User } from "../core";

export class FormatCommand extends Command {
	protected async onExecute(user: User, args: string[]): Promise<string | null> {
		if (args.length === 0) {
			return `@${user.displayName} ${EXAMPLE_FORMAT}`;
		}
		return `@${args[0]} ${EXAMPLE_FORMAT}`;
	}
}
