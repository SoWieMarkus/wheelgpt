import { Command } from "./command";
import { AccessLevel, type User } from "./user";

// Concrete implementation for testing
class TestCommand extends Command {
	protected async onExecute(user: User, args: string[]): Promise<string | null> {
		return `Handled by ${user.name} with args: ${args.join(",")}`;
	}
}

const makeUser = (accessLevel: number = AccessLevel.USER): User => ({
	id: "1",
	name: "tester",
	displayName: "Tester",
	accessLevel,
	channelId: "chan",
});

describe("Command", () => {
	const baseConfig = {
		name: "foo",
		cooldown: 1,
		aliases: ["alias1", "alias2"],
		accessLevel: AccessLevel.USER,
	};

	it("matches command name and aliases", () => {
		const cmd = new TestCommand("test", baseConfig);
		expect(cmd.matches({ key: "foo", args: [] })).toBe(true);
		expect(cmd.matches({ key: "alias1", args: [] })).toBe(true);
		expect(cmd.matches({ key: "alias2", args: [] })).toBe(true);
		expect(cmd.matches({ key: "bar", args: [] })).toBe(false);
	});

	it("is not on cooldown if never used", () => {
		const cmd = new TestCommand("test", baseConfig);
		expect(cmd.isOnCooldown()).toBe(false);
	});

	it("is on cooldown if used recently", () => {
		const cmd = new TestCommand("test", baseConfig);
		// biome-ignore lint/suspicious/noExplicitAny: Mock for testing purposes
		(cmd as any).lastUsed = Date.now();
		expect(cmd.isOnCooldown()).toBe(true);
	});

	it("is not on cooldown if enough time has passed", () => {
		const cmd = new TestCommand("test", baseConfig);
		// biome-ignore lint/suspicious/noExplicitAny: Mock for testing purposes
		(cmd as any).lastUsed = Date.now() - 2000;
		expect(cmd.isOnCooldown()).toBe(false);
	});

	it("execute returns null if user accessLevel is too low", async () => {
		const cmd = new TestCommand("test", { ...baseConfig, accessLevel: AccessLevel.MOD });
		const user = makeUser(AccessLevel.USER);
		expect(await cmd.execute(user, ["a"])).toBeNull();
	});

	it("execute returns null if on cooldown", async () => {
		const cmd = new TestCommand("test", baseConfig);

		// biome-ignore lint/suspicious/noExplicitAny: Mock for testing purposes
		(cmd as any).lastUsed = Date.now();
		const user = makeUser(AccessLevel.USER);
		expect(await cmd.execute(user, ["a"])).toBeNull();
	});

	it("execute calls handle and updates lastUsed if allowed", async () => {
		const cmd = new TestCommand("test", baseConfig);
		const user = makeUser(AccessLevel.USER);
		const result = await cmd.execute(user, ["x", "y"]);
		expect(result).toBe("Handled by tester with args: x,y");
	});
});
