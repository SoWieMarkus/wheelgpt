import type { CommandArguments } from "./command-arguments";
import type { User } from "./user";

export const matchesCommand = (
    command: Command,
    commandArguments: CommandArguments,
) => {
    const key = commandArguments.key;
    return command.name === key || command.aliases?.includes(key);
};

export const isCommandOnCooldown = (command: Command) => {
    const now = Date.now();
    if (command.lastUsed === undefined || command.cooldown === undefined)
        return false;
    return now - command.cooldown * 1000 < command.lastUsed;
};

export type Command = {
    name: string;
    cooldown?: number;
    task: (user: User, args: string[]) => string | null;
    aliases?: string[];
    accessLevel: number;
    lastUsed?: number;
};