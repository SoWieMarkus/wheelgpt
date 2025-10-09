export type CommandArguments = {
	key: string;
	args: string[];
};

const COMMAND_SYMBOL = "!";

export const getCommandArguments = (message: string): CommandArguments | null => {
	const tokens = message.split(" ");

	if (tokens.length === 0) return null;

	const firstToken = tokens[0];

	if (!firstToken.startsWith(COMMAND_SYMBOL)) {
		return null;
	}

	const key = firstToken.toLowerCase().replace("!", "");
	if (key.length === 0) {
		return null;
	}
	const args = tokens.slice(1).filter((arg) => {
		const visibleChars = arg.replace(/[^\x20-\x7E]/g, "").trim();
		return visibleChars.length > 0;
	});
	return { key, args };
};
