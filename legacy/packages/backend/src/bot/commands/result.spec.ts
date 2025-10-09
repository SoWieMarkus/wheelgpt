import type { User } from "../core";
import { EXAMPLE_FORMAT } from "../core";
import { GuessResultCommand, guessResultHandler } from "./result";

// Mock the guessResultHandler function
jest.mock("./result", () => ({
	...jest.requireActual("./result"),
	guessResultHandler: jest.fn(),
}));

// Type the mocks properly
const mockGuessResultHandler = guessResultHandler as jest.MockedFunction<typeof guessResultHandler>;

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

	test("should return result message when valid time provided", async () => {
		mockGuessResultHandler.mockResolvedValue("Some result message");

		const result = await guessResultCommand.execute(mockUser, ["30.000"]);

		expect(result).toBe(
			"YEK I got a new PB but I didn't know you are on a map? But no chatter participated ReallyFuckingMad",
		);
	});
});
