import { database } from "../../database";
import type { User } from "../core";
import { EXAMPLE_FORMAT, TrackmaniaTime } from "../core/time";
import { GuessCommand } from "./guess";

// Mock the database
jest.mock("../../database", () => ({
	database: {
		guess: {
			upsert: jest.fn(),
		},
	},
}));

// Mock TrackmaniaTime.parse
jest.mock("../core/time", () => ({
	...jest.requireActual("../core/time"),
	TrackmaniaTime: {
		parse: jest.fn(),
	},
}));

// Type the mocks properly
const mockUpsert = database.guess.upsert as jest.MockedFunction<typeof database.guess.upsert>;
const mockParse = TrackmaniaTime.parse as jest.MockedFunction<typeof TrackmaniaTime.parse>;

describe("GuessCommand", () => {
	let guessCommand: GuessCommand;
	let mockUser: User;

	beforeEach(() => {
		guessCommand = new GuessCommand("test-channel", {
			name: "guess",
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

	test("should return example format when no arguments provided", async () => {
		const result = await guessCommand.execute(mockUser, []);

		expect(result).toBe(`@TestUser ${EXAMPLE_FORMAT}`);
		expect(mockParse).not.toHaveBeenCalled();
		expect(mockUpsert).not.toHaveBeenCalled();
	});

	test("should return example format when time parsing fails", async () => {
		mockParse.mockReturnValue(null);

		const result = await guessCommand.execute(mockUser, ["invalid-time"]);

		expect(result).toBe(`@TestUser ${EXAMPLE_FORMAT}`);
		expect(mockParse).toHaveBeenCalledWith("invalid-time");
		expect(mockUpsert).not.toHaveBeenCalled();
	});

	test("should store guess and return null when valid time provided", async () => {
		const mockTime = {
			getTotalInMilliSeconds: jest.fn().mockReturnValue(30000),
		} as unknown as jest.Mocked<TrackmaniaTime>;

		const mockGuess = {
			channelId: "test-channel",
			userId: "123",
			displayName: "TestUser",
			time: 30000,
			createdAt: new Date(),
		};

		mockParse.mockReturnValue(mockTime);
		mockUpsert.mockResolvedValue(mockGuess);

		const result = await guessCommand.execute(mockUser, ["30.000"]);

		expect(result).toBeNull();
		expect(mockParse).toHaveBeenCalledWith("30.000");
		expect(mockTime.getTotalInMilliSeconds).toHaveBeenCalled();
		expect(mockUpsert).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "123",
					channelId: "test-channel",
				},
			},
			create: {
				userId: "123",
				displayName: "TestUser",
				channelId: "test-channel",
				time: 30000,
			},
			update: {
				time: 30000,
			},
		});
	});

	test("should use only first argument when multiple provided", async () => {
		const mockTime = {
			getTotalInMilliSeconds: jest.fn().mockReturnValue(45000),
		} as unknown as jest.Mocked<TrackmaniaTime>;

		const mockGuess = {
			channelId: "test-channel",
			userId: "123",
			displayName: "TestUser",
			time: 45000,
			createdAt: new Date(),
		};

		mockParse.mockReturnValue(mockTime);
		mockUpsert.mockResolvedValue(mockGuess);

		const result = await guessCommand.execute(mockUser, ["45.000", "ignored", "args"]);

		expect(result).toBeNull();
		expect(mockParse).toHaveBeenCalledWith("45.000");
		expect(mockParse).toHaveBeenCalledTimes(1);
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

		const mockTime = {
			getTotalInMilliSeconds: jest.fn().mockReturnValue(25000),
		} as unknown as jest.Mocked<TrackmaniaTime>;

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
			time: 25000,
			createdAt: new Date(),
		};

		mockParse.mockReturnValue(mockTime);
		mockUpsert.mockResolvedValueOnce(mockGuess1);
		mockUpsert.mockResolvedValueOnce(mockGuess2);

		await guessCommand.execute(user1, ["25.000"]);
		await guessCommand.execute(user2, ["25.000"]);

		expect(mockUpsert).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "1",
					channelId: "test-channel",
				},
			},
			create: {
				userId: "1",
				displayName: "Alice",
				channelId: "test-channel",
				time: 25000,
			},
			update: {
				time: 25000,
			},
		});

		expect(mockUpsert).toHaveBeenCalledWith({
			where: {
				channelId_userId: {
					userId: "2",
					channelId: "test-channel",
				},
			},
			create: {
				userId: "2",
				displayName: "Bob",
				channelId: "test-channel",
				time: 25000,
			},
			update: {
				time: 25000,
			},
		});
	});

	test("should handle invalid time formats", async () => {
		const invalidTimeFormats = ["not-a-time", "25:70.000", "abc", "12.34.56", ""];

		for (const invalidTime of invalidTimeFormats) {
			mockParse.mockReturnValue(null);

			const result = await guessCommand.execute(mockUser, [invalidTime]);

			expect(result).toBe(`@TestUser ${EXAMPLE_FORMAT}`);
			expect(mockParse).toHaveBeenCalledWith(invalidTime);
			expect(mockUpsert).not.toHaveBeenCalled();

			jest.clearAllMocks();
		}
	});
});
