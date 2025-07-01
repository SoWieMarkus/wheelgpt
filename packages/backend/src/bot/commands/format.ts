import { Command, EXAMPLE_FORMAT, mentionUser, type User } from "../core";

export class FormatCommand extends Command {
	protected async onExecute(user: User, args: string[]): Promise<string | null> {
		if (args.length === 0) {
			return `${mentionUser(user.displayName)} ${EXAMPLE_FORMAT}`;
		}
		return `${mentionUser(args[0])} ${EXAMPLE_FORMAT}`;
	}
}
