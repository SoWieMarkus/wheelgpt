import { getLeaderboard, getLeaderboardByName } from "@prisma/client/sql";
import { database } from "../../database";
import type { User } from "../core";
import { LeaderboardCommand } from "./leaderboard";

// Mock the database
jest.mock("../../database", () => ({
	database: {
		$queryRawTyped: jest.fn(),
	},
}));

// Type the mock properly
const mockQueryRawTyped = database.$queryRawTyped as jest.MockedFunction<typeof database.$queryRawTyped>;

describe("LeaderboardCommand", () => {
	let leaderboardCommand: LeaderboardCommand;
	let mockUser: User;

	beforeEach(() => {
		leaderboardCommand = new LeaderboardCommand("test-channel", {
			name: "leaderboard",
		});
		mockUser = {
			name: "testuser",
			id: "123",
			displayName: "TestUser",
			accessLevel: 1,
			channelId: "test-channel",
		};
		jest.clearAllMocks();
	});

	test("should return top leaderboard when no arguments provided", async () => {
		const mockLeaderboard = [
			{ position: "1", displayName: "Alice", points: 15 },
			{ position: "2", displayName: "Bob", points: 10 },
			{ position: "3", displayName: "Charlie", points: 5 },
		];

		mockQueryRawTyped.mockResolvedValue(mockLeaderboard);

		const result = await leaderboardCommand.execute(mockUser, []);

		expect(result).toBe(
			"Top 5 Guessers 🏆 Perfect guess = 5pts, Closest guess = 1pt: 1. Alice - 15 | 2. Bob - 10 | 3. Charlie - 5 | You can view the full leaderboard here: https://wheelgpt.dev/leaderboard/test-channel",
		);
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboard("test-channel", 5, 0));
	});

	test("should return no guesses message when leaderboard is empty", async () => {
		mockQueryRawTyped.mockResolvedValue([]);

		const result = await leaderboardCommand.execute(mockUser, []);

		expect(result).toBe("No guesses have been made yet.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboard("test-channel", 5, 0));
	});

	test("should return specific user position when username provided", async () => {
		const mockUserEntry = [{ displayName: "Alice", position: "3", points: 8 }];

		mockQueryRawTyped.mockResolvedValue(mockUserEntry);

		const result = await leaderboardCommand.execute(mockUser, ["Alice"]);

		expect(result).toBe("@Alice is currently in position 3 with 8 pts.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("test-channel", "alice"));
	});

	test("should handle singular point correctly", async () => {
		const mockUserEntry = [{ displayName: "Bob", position: "5", points: 1 }];

		mockQueryRawTyped.mockResolvedValue(mockUserEntry);

		const result = await leaderboardCommand.execute(mockUser, ["Bob"]);

		expect(result).toBe("@Bob is currently in position 5 with 1 pt.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("test-channel", "bob"));
	});

	test("should remove @ symbol from username", async () => {
		const mockUserEntry = [{ displayName: "Charlie", position: "2", points: 12 }];

		mockQueryRawTyped.mockResolvedValue(mockUserEntry);

		const result = await leaderboardCommand.execute(mockUser, ["@Charlie"]);

		expect(result).toBe("@Charlie is currently in position 2 with 12 pts.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("test-channel", "charlie"));
	});

	test("should convert username to lowercase", async () => {
		const mockUserEntry = [{ displayName: "TestUser", position: "1", points: 20 }];

		mockQueryRawTyped.mockResolvedValue(mockUserEntry);

		const result = await leaderboardCommand.execute(mockUser, ["TESTUSER"]);

		expect(result).toBe("@TestUser is currently in position 1 with 20 pts.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("test-channel", "testuser"));
	});

	test("should return not found message when user not in leaderboard", async () => {
		mockQueryRawTyped.mockResolvedValue([]);

		const result = await leaderboardCommand.execute(mockUser, ["NonExistentUser"]);

		expect(result).toBe("@TestUser No leaderboard entry found for nonexistentuser.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("test-channel", "nonexistentuser"));
	});

	test("should handle single entry leaderboard", async () => {
		const mockLeaderboard = [{ position: "1", displayName: "OnlyUser", points: 3 }];

		mockQueryRawTyped.mockResolvedValue(mockLeaderboard);

		const result = await leaderboardCommand.execute(mockUser, []);

		expect(result).toBe(
			"Top 5 Guessers 🏆 Perfect guess = 5pts, Closest guess = 1pt: 1. OnlyUser - 3 | You can view the full leaderboard here: https://wheelgpt.dev/leaderboard/test-channel",
		);
	});

	test("should work with different channels", async () => {
		const channel1Command = new LeaderboardCommand("channel-1", {
			name: "leaderboard",
		});
		const channel2Command = new LeaderboardCommand("channel-2", {
			name: "leaderboard",
		});

		mockQueryRawTyped.mockResolvedValue([]);

		await channel1Command.execute(mockUser, []);
		await channel2Command.execute(mockUser, ["Alice"]);

		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboard("channel-1", 5, 0));
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("channel-2", "alice"));
	});

	test("should use only first argument when multiple provided", async () => {
		const mockUserEntry = [{ displayName: "FirstUser", position: "4", points: 7 }];

		mockQueryRawTyped.mockResolvedValue(mockUserEntry);

		const result = await leaderboardCommand.execute(mockUser, ["FirstUser", "ignored", "args"]);

		expect(result).toBe("@FirstUser is currently in position 4 with 7 pts.");
		expect(mockQueryRawTyped).toHaveBeenCalledWith(getLeaderboardByName("test-channel", "firstuser"));
	});

	test("should include correct channel URL in leaderboard message", async () => {
		const differentChannelCommand = new LeaderboardCommand("my-special-channel", {
			name: "leaderboard",
		});

		const mockLeaderboard = [{ position: "1", displayName: "User1", points: 5 }];

		mockQueryRawTyped.mockResolvedValue(mockLeaderboard);

		const result = await differentChannelCommand.execute(mockUser, []);

		expect(result).toContain("https://wheelgpt.dev/leaderboard/my-special-channel");
	});
});
