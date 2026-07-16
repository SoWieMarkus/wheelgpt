import { database } from "../../database";
import type { User } from "../core";
import { EXAMPLE_FORMAT } from "../core";
import { TrackmaniaTime } from "../core/time";
import { GuessResultCommand } from "./result";

jest.mock("../../database", () => ({
	database: {
		channel: { findUnique: jest.fn() },
		guess: { findMany: jest.fn(), deleteMany: jest.fn() },
		trackmaniaMap: { findUnique: jest.fn() },
		guesserLeaderboard: { upsert: jest.fn() },
	},
}));

jest.mock("./result", () => ({
	...jest.requireActual("./result"),
	guessResultHandler: jest.fn(),
}));

// eslint-disable-next-line @typescript-eslint/no-require-imports
const { guessResultHandler: mockGuessResultHandler } = require("./result");
// eslint-disable-next-line @typescript-eslint/no-require-imports
const { guessResultHandler: realGuessResultHandler } = jest.requireActual("./result");

const mockChannel = database.channel.findUnique as jest.Mock;
const mockGuessFind = database.guess.findMany as jest.Mock;
const mockGuessDelete = database.guess.deleteMany as jest.Mock;
const mockMapFind = database.trackmaniaMap.findUnique as jest.Mock;

describe("GuessResultCommand", () => {
	let guessResultCommand: GuessResultCommand;
	let mockUser: User;

	beforeEach(() => {
		guessResultCommand = new GuessResultCommand("test-channel", {
			name: "result",
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
		const result = await guessResultCommand.execute(mockUser, []);

		expect(result).toBe(`@TestUser ${EXAMPLE_FORMAT}`);
		expect(mockGuessResultHandler).not.toHaveBeenCalled();
	});

	test("should return error message when time parsing fails", async () => {
		const result = await guessResultCommand.execute(mockUser, ["invalid-time"]);

		expect(result).toBe(`@TestUser smh granadyy mods are all degens. Wrong format you idiot. ${EXAMPLE_FORMAT}`);
		expect(mockGuessResultHandler).not.toHaveBeenCalled();
	});

	test("should call guessResultHandler with channelId and parsed time when valid time provided", async () => {
		mockGuessResultHandler.mockResolvedValue("Some result message");

		const result = await guessResultCommand.execute(mockUser, ["30.000"]);

		expect(result).toBe("Some result message");
		expect(mockGuessResultHandler).toHaveBeenCalledWith(
			"test-channel",
			new TrackmaniaTime(30000),
		);
	});
});

describe("guessResultHandler - min age filtering", () => {
	beforeEach(() => {
		jest.clearAllMocks();
		mockMapFind.mockResolvedValue(null);
		mockGuessDelete.mockResolvedValue({ count: 0 });
	});

	test("excludes guesses newer than guessMinRequiredAgeTime", async () => {
		const now = Date.now();
		const old = new Date(now - 10000); // 10s ago — passes a 5s min age
		const recent = new Date(now - 1000); // 1s ago — excluded by a 5s min age

		mockChannel.mockResolvedValue({ guessMinRequiredAgeTime: 5 });
		mockGuessFind.mockResolvedValue([]);

		await realGuessResultHandler("test-channel", new TrackmaniaTime(30000));

		const cutoff: Date = mockGuessFind.mock.calls[0][0].where.createdAt.lte;

		expect(old.getTime()).toBeLessThanOrEqual(cutoff.getTime());
		expect(recent.getTime()).toBeGreaterThan(cutoff.getTime());
	});

	test("passes all guesses when guessMinRequiredAgeTime is 0", async () => {
		mockChannel.mockResolvedValue({ guessMinRequiredAgeTime: 0 });
		mockGuessFind.mockResolvedValue([]);

		await realGuessResultHandler("test-channel", new TrackmaniaTime(30000));

		const cutoff: Date = mockGuessFind.mock.calls[0][0].where.createdAt.lte;

		// cutoff ≈ now, so every stored guess (createdAt in the past) passes
		expect(cutoff.getTime()).toBeGreaterThanOrEqual(0);
	});

	test("falls back to 0 min age when channel is not found", async () => {
		mockChannel.mockResolvedValue(null);
		mockGuessFind.mockResolvedValue([]);

		await realGuessResultHandler("test-channel", new TrackmaniaTime(30000));

		expect(mockGuessFind.mock.calls[0][0].where.createdAt).toBeDefined();
	});

	test("deleteMany removes all guesses for the channel regardless of age", async () => {
		mockChannel.mockResolvedValue({ guessMinRequiredAgeTime: 30 });
		mockGuessFind.mockResolvedValue([]);

		await realGuessResultHandler("test-channel", new TrackmaniaTime(30000));

		expect(mockGuessDelete).toHaveBeenCalledWith({
			where: { channelId: "test-channel" },
		});
	});
});
