import { database } from "../../database";
import type { User } from "../core";
import { TrackmaniaMap } from "../core";
import { MapCommand } from "./map";

// Mock the database
jest.mock("../../database", () => ({
	database: {
		trackmaniaMap: {
			findUnique: jest.fn(),
		},
	},
}));

// Mock the TrackmaniaMap class
jest.mock("../core", () => ({
	...jest.requireActual("../core"),
	TrackmaniaMap: jest.fn(),
}));

// Type the mocks properly
const mockFindUnique = database.trackmaniaMap.findUnique as jest.MockedFunction<
	typeof database.trackmaniaMap.findUnique
>;
const MockedTrackmaniaMap = TrackmaniaMap as jest.MockedClass<typeof TrackmaniaMap>;

describe("MapCommand", () => {
	let mapCommand: MapCommand;
	let mockUser: User;

	beforeEach(() => {
		mapCommand = new MapCommand("test-channel", {
			name: "map",
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

	test("should return no map message when no map found", async () => {
		mockFindUnique.mockResolvedValue(null);

		const result = await mapCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser No map is currently set for this channel.");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});

	test("should return map information when map exists", async () => {
		const mockMap = {
			channelId: "test-channel",
			name: "Test Track",
			uid: "test-uid",
			authorTime: 60000,
			goldTime: 65000,
			silverTime: 70000,
			bronzeTime: 80000,
			championTime: 58000,
			worldRecord: null,
			author: "Author",
			tmxId: null,
		};

		const mockTrackmaniaMapInstance = {
			toString: jest.fn().mockReturnValue("Test Track (Champion: 43.000)"),
		} as unknown as jest.Mocked<TrackmaniaMap>;

		mockFindUnique.mockResolvedValue(mockMap);
		MockedTrackmaniaMap.mockImplementation(() => mockTrackmaniaMapInstance);

		const result = await mapCommand.execute(mockUser, []);

		expect(result).toBe("@TestUser Test Track (Champion: 43.000)");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});

	test("should ignore command arguments", async () => {
		const mockMap = {
			channelId: "test-channel",
			name: "Test Track",
			uid: "test-uid",
			authorTime: 45000,
			goldTime: 50000,
			silverTime: 55000,
			bronzeTime: 65000,
			championTime: 43000,
			worldRecord: null,
			author: "Author",
			tmxId: null,
		};

		const mockTrackmaniaMapInstance = {
			toString: jest.fn().mockReturnValue("Test Track (Champion: 43.000)"),
		} as unknown as jest.Mocked<TrackmaniaMap>;

		mockFindUnique.mockResolvedValue(mockMap);
		MockedTrackmaniaMap.mockImplementation(() => mockTrackmaniaMapInstance);

		// Pass some arguments that should be ignored
		const result = await mapCommand.execute(mockUser, ["ignored", "args"]);

		expect(result).toBe("@TestUser Test Track (Champion: 43.000)");
		expect(mockFindUnique).toHaveBeenCalledWith({
			where: {
				channelId: "test-channel",
			},
		});
	});
});
