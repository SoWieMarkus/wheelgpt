import type { User } from "../core";
import { Emote } from "../core/emotes";
import { EmotesCommand } from "./emotes";

describe("EmotesCommand", () => {
	let emotesCommand: EmotesCommand;
	let mockUser: User;

	beforeEach(() => {
		emotesCommand = new EmotesCommand("test-channel", {
			name: "emotes",
		});
		mockUser = {
			name: "testuser",
			id: "123",
			displayName: "TestUser",
			accessLevel: 1,
			channelId: "test-channel",
		};
	});

	test("should return message with user mention and all emotes", async () => {
		const result = await emotesCommand.execute(mockUser, []);

		expect(result).not.toBeNull();
		expect(result).toContain("@TestUser");
	});

	test("should include all emotes from Emote enum", async () => {
		const result = await emotesCommand.execute(mockUser, []);

		const allEmoteNames = Object.values(Emote).map((emote) => emote.name);

		for (const emoteName of allEmoteNames) {
			expect(result).toContain(emoteName);
		}
	});

	test("should format emotes with spaces between them", async () => {
		const result = await emotesCommand.execute(mockUser, []);

		const expectedEmotes = Object.values(Emote)
			.map((emote) => emote.name)
			.join(" ");

		expect(result).toContain(expectedEmotes);
	});

	test("should mention the user who executed the command", async () => {
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

		const result1 = await emotesCommand.execute(user1, []);
		const result2 = await emotesCommand.execute(user2, []);

		expect(result1).toContain("@Alice");
		expect(result2).toContain("@Bob");
		expect(result1).not.toContain("@Bob");
		expect(result2).not.toContain("@Alice");
	});

	test("should ignore command arguments", async () => {
		const result1 = await emotesCommand.execute(mockUser, []);
		const result2 = await emotesCommand.execute(mockUser, ["ignored", "args"]);
		const result3 = await emotesCommand.execute(mockUser, ["many", "different", "arguments", "here"]);

		expect(result1).toBe(result2);
		expect(result2).toBe(result3);
	});

	test("should handle users with special characters in display name", async () => {
		const specialUser: User = {
			name: "user_with_special",
			id: "999",
			displayName: "User_With_123!",
			accessLevel: 1,
			channelId: "test-channel",
		};

		const result = await emotesCommand.execute(specialUser, []);

		expect(result).toContain("@User_With_123!");
	});

	test("should be consistent output for same user", async () => {
		const result1 = await emotesCommand.execute(mockUser, []);
		const result2 = await emotesCommand.execute(mockUser, []);
		const result3 = await emotesCommand.execute(mockUser, []);

		expect(result1).toBe(result2);
		expect(result2).toBe(result3);
	});
});
