import type { User } from "../core";
import { EXAMPLE_FORMAT } from "../core";
import { FormatCommand } from "./format";

describe("FormatCommand", () => {
	let formatCommand: FormatCommand;
	let mockUser: User;

	beforeEach(() => {
		formatCommand = new FormatCommand("test-channel", {
			name: "format",
		});
		mockUser = {
			name: "testuser",
			id: "123",
			displayName: "TestUser",
			accessLevel: 1,
			channelId: "test-channel",
		};
	});

	test("should return message with user mention when no arguments provided", async () => {
		const result = await formatCommand.execute(mockUser, []);

		expect(result).toBe(`@TestUser ${EXAMPLE_FORMAT}`);
	});

	test("should return message with first argument as mention when arguments provided", async () => {
		const result = await formatCommand.execute(mockUser, ["SomeUsername"]);
		expect(result).toBe(`@SomeUsername ${EXAMPLE_FORMAT}`);
	});

	test("should use only first argument and ignore others", async () => {
		const result = await formatCommand.execute(mockUser, ["FirstUser", "SecondUser", "ThirdUser"]);

		expect(result).toBe(`@FirstUser ${EXAMPLE_FORMAT}`);
		expect(result).not.toContain("SecondUser");
		expect(result).not.toContain("ThirdUser");
	});

	test("should handle whitespace in username argument", async () => {
		const result = await formatCommand.execute(mockUser, ["User With Spaces"]);

		expect(result).toBe(`@User With Spaces ${EXAMPLE_FORMAT}`);
	});

	test("should always include EXAMPLE_FORMAT in response", async () => {
		const result1 = await formatCommand.execute(mockUser, []);
		const result2 = await formatCommand.execute(mockUser, ["SomeUser"]);

		expect(result1).toContain(EXAMPLE_FORMAT);
		expect(result2).toContain(EXAMPLE_FORMAT);
	});

	test("should mention different users based on display name vs argument", async () => {
		const user1: User = {
			name: "alice",
			id: "1",
			displayName: "Alice",
			accessLevel: 1,
			channelId: "test-channel",
		};

		const user2: User = {
			name: "bob",
			id: "2",
			displayName: "Bob",
			accessLevel: 1,
			channelId: "test-channel",
		};

		const resultNoArgs1 = await formatCommand.execute(user1, []);
		const resultNoArgs2 = await formatCommand.execute(user2, []);
		const resultWithArg = await formatCommand.execute(user1, ["Charlie"]);

		expect(resultNoArgs1).toContain("@Alice");
		expect(resultNoArgs2).toContain("@Bob");
		expect(resultWithArg).toContain("@Charlie");
		expect(resultWithArg).not.toContain("@Alice");
	});

	test("should be case sensitive for usernames", async () => {
		const result1 = await formatCommand.execute(mockUser, ["username"]);
		const result2 = await formatCommand.execute(mockUser, ["USERNAME"]);
		const result3 = await formatCommand.execute(mockUser, ["UserName"]);

		expect(result1).toBe(`@username ${EXAMPLE_FORMAT}`);
		expect(result2).toBe(`@USERNAME ${EXAMPLE_FORMAT}`);
		expect(result3).toBe(`@UserName ${EXAMPLE_FORMAT}`);
		expect(result1).not.toBe(result2);
	});
});
