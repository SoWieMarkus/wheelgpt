import type { User } from "../core";
import { Emote } from "../core/emotes";
import { WheelGPTCommand } from "./wheelgpt";

describe("WheelGPTCommand", () => {
	let wheelgptCommand: WheelGPTCommand;
	let mockUser: User;

	beforeEach(() => {
		wheelgptCommand = new WheelGPTCommand("test-channel", {
			name: "wheelgpt",
		});
		mockUser = {
			name: "testuser",
			id: "123",
			displayName: "TestUser",
			accessLevel: 1,
			channelId: "test-channel",
		};
	});

	test("should return information message about the bot", async () => {
		const result = await wheelgptCommand.execute(mockUser, []);

		const expectedMessage = `${Emote.HEY_GUYS.name} I am a Trackmania Twitch Bot! You can ask me about the current map, the current room  and guess the next best time. More details: https://wheelgpt.dev | ⭐ Star me on Github: https://github.com/SoWieMarkus/wheelgpt`;

		expect(result).toBe(expectedMessage);
	});

	test("should return same message regardless of arguments", async () => {
		const result1 = await wheelgptCommand.execute(mockUser, []);
		const result2 = await wheelgptCommand.execute(mockUser, ["arg1", "arg2"]);
		const result3 = await wheelgptCommand.execute(mockUser, ["random", "arguments", "here"]);

		expect(result1).toBe(result2);
		expect(result2).toBe(result3);
	});

	test("should return same message regardless of user", async () => {
		const user1: User = {
			name: "user1",
			id: "1",
			displayName: "User1",
			accessLevel: 1,
			channelId: "test-channel",
		};

		const user2: User = {
			name: "user2",
			id: "2",
			displayName: "User2",
			accessLevel: 5,
			channelId: "test-channel",
		};

		const result1 = await wheelgptCommand.execute(user1, []);
		const result2 = await wheelgptCommand.execute(user2, []);

		expect(result1).toBe(result2);
	});

	test("should include correct URLs in message", async () => {
		const result = await wheelgptCommand.execute(mockUser, []);

		expect(result).toContain("https://wheelgpt.dev");
		expect(result).toContain("https://github.com/SoWieMarkus/wheelgpt");
	});
});
