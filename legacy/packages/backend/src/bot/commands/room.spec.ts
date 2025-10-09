import { database } from "../../database";
import type { User } from "../core";
import { RoomCommand } from "./room";

// Mock the database
jest.mock("../../database", () => ({
	database: {
		trackmaniaRoom: {
			findUnique: jest.fn(),
		},
	},
}));

// Type the mock properly
const mockFindUnique = database.trackmaniaRoom.findUnique as jest.MockedFunction<
	typeof database.trackmaniaRoom.findUnique
>;

describe("RoomCommand", () => {
	let roomCommand: RoomCommand;
	let mockUser: User;

	beforeEach(() => {
		roomCommand = new RoomCommand("test-channel", {
			name: "room",
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

	test("should return not in room message when no room found", async () => {
		mockFindUnique.mockResolvedValue(null);

		const result = await roomCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser I am currently not in a room.");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});

	test("should return room information when room exists", async () => {
		const mockRoom = {
			channelId: "test-channel",
			login: "test-room",
			name: "Cool Track Server",
			numberOfPlayers: 15,
			maxPlayers: 20,
		};

		mockFindUnique.mockResolvedValue(mockRoom);

		const result = await roomCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser Cool Track Server [15/20]");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});

	test("should ignore command arguments", async () => {
		const mockRoom = {
			channelId: "test-channel",
			login: "test-room",
			name: "Test Server",
			numberOfPlayers: 3,
			maxPlayers: 8,
		};

		mockFindUnique.mockResolvedValue(mockRoom);

		// Pass some arguments that should be ignored
		const result = await roomCommand.execute(mockUser, ["ignored", "args"]);

		expect(result).toBe("@TestUser Test Server [3/8]");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});
});
