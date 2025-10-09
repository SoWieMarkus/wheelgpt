import { database } from "../../database";
import type { User } from "../core";
import { Emote } from "../core/emotes";
import { MyGuessCommand } from "./myguess";

// Mock the database
jest.mock("../../database", () => ({
	database: {
		guess: {
			findUnique: jest.fn(),
		},
	},
}));

// Type the mock properly
const mockFindUnique = database.guess.findUnique as jest.MockedFunction<typeof database.guess.findUnique>;

describe("MyGuessCommand", () => {
	let myGuessCommand: MyGuessCommand;
	let mockUser: User;

	beforeEach(() => {
		myGuessCommand = new MyGuessCommand("test-channel", {
			name: "myguess",
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

	test("should return no guess found message when user has no guess", async () => {
		mockFindUnique.mockResolvedValue(null);

		const result = await myGuessCommand.execute(mockUser, []);

		expect(result).toBe(`@TestUser I can't find any guess from you ${Emote.YEK.name}`);
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "123",
					channelId: "test-channel",
				},
			},
		});
	});

	test("should return user guess when guess exists", async () => {
		const mockGuess = {
			channelId: "test-channel",
			userId: "123",
			displayName: "TestUser",
			time: 30000, // 30 seconds in milliseconds
			createdAt: new Date(),
		};

		mockFindUnique.mockResolvedValue(mockGuess);

		const result = await myGuessCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser 30.000");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "123",
					channelId: "test-channel",
				},
			},
		});
	});

	test("should format time correctly for different time values", async () => {
		const testCases = [
			{ time: 12345, expected: "12.345" }, // 12.345 seconds
			{ time: 60000, expected: "1:00.000" }, // 1 minute
			{ time: 3661000, expected: "1:01:01.000" }, // 1 hour 1 minute 1 second
			{ time: 500, expected: "0.500" }, // 0.5 seconds
			{ time: 123456, expected: "2:03.456" }, // 2 minutes 3.456 seconds
		];

		for (const testCase of testCases) {
			const mockGuess = {
				channelId: "test-channel",
				userId: "123",
				displayName: "TestUser",
				time: testCase.time,
				createdAt: new Date(),
			};

			mockFindUnique.mockResolvedValue(mockGuess);

			const result = await myGuessCommand.execute(mockUser, []);

			expect(result).toBe(`@TestUser ${testCase.expected}`);
		}
	});

	test("should work with different users", async () => {
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

		const mockGuess1 = {
			channelId: "test-channel",
			userId: "1",
			displayName: "Alice",
			time: 25000,
			createdAt: new Date(),
		};

		const mockGuess2 = {
			channelId: "test-channel",
			userId: "2",
			displayName: "Bob",
			time: 35000,
			createdAt: new Date(),
		};

		mockFindUnique.mockResolvedValueOnce(mockGuess1);
		const result1 = await myGuessCommand.execute(user1, []);

		mockFindUnique.mockResolvedValueOnce(mockGuess2);
		const result2 = await myGuessCommand.execute(user2, []);

		expect(result1).toBe("@Alice 25.000");
		expect(result2).toBe("@Bob 35.000");

		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "1",
					channelId: "test-channel",
				},
			},
		});

		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "2",
					channelId: "test-channel",
				},
			},
		});
	});

	test("should ignore command arguments", async () => {
		const mockGuess = {
			channelId: "test-channel",
			userId: "123",
			displayName: "TestUser",
			time: 30000,
			createdAt: new Date(),
		};

		mockFindUnique.mockResolvedValue(mockGuess);

		const result1 = await myGuessCommand.execute(mockUser, []);
		const result2 = await myGuessCommand.execute(mockUser, ["ignored", "args"]);

		expect(result1).toBe(result2);
		expect(mockFindUnique).toHaveBeenCalledTimes(2);

		// Verify both calls had identical parameters
		const calls = mockFindUnique.mock.calls;
		expect(calls[0]).toEqual(calls[1]);
	});

	test("should handle very large time values", async () => {
		const mockGuess = {
			channelId: "test-channel",
			userId: "123",
			displayName: "TestUser",
			time: 7200000, // 2 hours
		};

		mockFindUnique.mockResolvedValue(mockGuess);

		const result = await myGuessCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser 2:00:00.000");
	});

	test("should use correct composite key for database query", async () => {
		const mockGuess = {
			channelId: "test-channel",
			userId: "123",
			displayName: "TestUser",
			time: 30000,
		};

		mockFindUnique.mockResolvedValue(mockGuess);

		await myGuessCommand.execute(mockUser, []);

		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "123",
					channelId: "test-channel",
				},
			},
		});
	});
});
