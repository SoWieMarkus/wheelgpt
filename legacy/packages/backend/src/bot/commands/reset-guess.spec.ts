import { database } from "../../database";
import type { User } from "../core";
import { ResetGuessesCommand } from "./reset-guesses";

// Mock the database
jest.mock("../../database", () => ({
	database: {
		guess: {
			deleteMany: jest.fn(),
		},
	},
}));

// Type the mock properly
const mockDeleteMany = database.guess.deleteMany as jest.MockedFunction<typeof database.guess.deleteMany>;

describe("ResetGuessesCommand", () => {
	let resetGuessesCommand: ResetGuessesCommand;
	let mockUser: User;

	beforeEach(() => {
		resetGuessesCommand = new ResetGuessesCommand("test-channel", {
			name: "reset-guesses",
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

	test("should delete all guesses and return success message", async () => {
		mockDeleteMany.mockResolvedValue({ count: 3 });

		const result = await resetGuessesCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser All guesses have been reset.");
		expect(mockDeleteMany).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
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

		mockDeleteMany.mockResolvedValue({ count: 2 });

		const result1 = await resetGuessesCommand.execute(user1, []);
		const result2 = await resetGuessesCommand.execute(user2, []);

		expect(result1).toBe("@Alice All guesses have been reset.");
		expect(result2).toBe("@Bob All guesses have been reset.");
	});

	test("should ignore command arguments", async () => {
		mockDeleteMany.mockResolvedValue({ count: 1 });

		// Pass some arguments that should be ignored
		const result = await resetGuessesCommand.execute(mockUser, ["ignored", "args"]);

		expect(result).toBe("@TestUser All guesses have been reset.");
		expect(mockDeleteMany).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});

	test("should return success message even when no guesses exist", async () => {
		mockDeleteMany.mockResolvedValue({ count: 0 });

		const result = await resetGuessesCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser All guesses have been reset.");
		expect(mockDeleteMany).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});
});
