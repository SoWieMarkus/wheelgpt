import type { Guess } from "@prisma/client";
import { Emote } from "./emotes";
import {
	buildBestGuesserMessage,
	buildGuessResultMessage,
	buildPbMapResultMessage,
	evaluateGuesses,
	POINTS_CLOSEST_GUESS,
	POINTS_PERFECT_GUESS,
} from "./guess";
import { Medal, type TrackmaniaMap } from "./map";
import { TrackmaniaTime } from "./time";

describe("buildBestGuesserMessage", () => {
	test("should return no participants message when no winners", () => {
		const winners: Guess[] = [];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildBestGuesserMessage(winners, newBestTime);

		expect(result).toBe(`But no chatter participated ${Emote.REALLY_FUCKING_MAD.name}`);
	});

	test("should return perfect guess message for single winner", () => {
		const winners: Guess[] = [
			{
				channelId: "test",
				userId: "1",
				displayName: "Alice",
				time: 30000,
			},
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildBestGuesserMessage(winners, newBestTime);

		expect(result).toBe(
			`@Alice the ${Emote.GIGACHAD.name} guessed it correctly! ${Emote.BWOAH.name} (+${POINTS_PERFECT_GUESS} points)`,
		);
	});

	test("should return perfect guess message for multiple winners", () => {
		const winners: Guess[] = [
			{ channelId: "test", userId: "1", displayName: "Alice", time: 30000 },
			{ channelId: "test", userId: "2", displayName: "Bob", time: 30000 },
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildBestGuesserMessage(winners, newBestTime);

		expect(result).toBe(
			`@Alice, @Bob the ${Emote.GIGACHAD.name} 's guessed it correctly! ${Emote.BWOAH.name} (+${POINTS_PERFECT_GUESS} points)`,
		);
	});

	test("should return closest guess message for single winner", () => {
		const winners: Guess[] = [
			{
				channelId: "test",
				userId: "1",
				displayName: "Alice",
				time: 29500, // 0.5 seconds off
			},
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildBestGuesserMessage(winners, newBestTime);

		expect(result).toBe(
			`Nobody guessed it correctly but @Alice guessed 29.500 (-0.500) ${Emote.OK.name} (+${POINTS_CLOSEST_GUESS} point)`,
		);
	});

	test("should return closest guess message for multiple winners", () => {
		const winners: Guess[] = [
			{ channelId: "test", userId: "1", displayName: "Alice", time: 29500 },
			{ channelId: "test", userId: "2", displayName: "Bob", time: 29500 },
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildBestGuesserMessage(winners, newBestTime);

		expect(result).toBe(
			`Nobody guessed it correctly but @Alice, @Bob guessed 29.500 (-0.500) ${Emote.OK.name} (+${POINTS_CLOSEST_GUESS} point)`,
		);
	});

	test("should show positive difference when guess is higher than result", () => {
		const winners: Guess[] = [
			{
				channelId: "test",
				userId: "1",
				displayName: "Alice",
				time: 31000, // 1 second over
			},
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildBestGuesserMessage(winners, newBestTime);

		expect(result).toBe(
			`Nobody guessed it correctly but @Alice guessed 31.000 (+1.000) ${Emote.OK.name} (+${POINTS_CLOSEST_GUESS} point)`,
		);
	});
});

describe("buildPbMapResultMessage", () => {
	test("should return unknown map message when no map", () => {
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(null, newBestTime);

		expect(result).toBe(`${Emote.YEK.name} I got a new PB but I didn't know you are on a map?`);
	});

	test("should return champion medal message", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.CHAMPION),
			getDifferenceToMedal: jest.fn(),
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's Champion Medal ${Emote.CHAMPION_MEDAL.name} ${Emote.BWOAH.name}`,
		);
	});

	test("should return missed champion message when close to champion", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.AUTHOR),
			getDifferenceToMedal: jest.fn().mockReturnValue(5), // 5ms from champion
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} HAHAHAHAH HE MISSED CHAMPION BY 0.005s ${Emote.OMEGALUL.name} ${Emote.ICANT.name} ${Emote.PEPE_POINT.name}`,
		);
	});

	test("should return author medal message", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.AUTHOR),
			getDifferenceToMedal: jest.fn().mockReturnValue(1000), // 1 second from champion
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's ${Emote.AUTHOR_MEDAL.name} ${Emote.OK.name}`,
		);
	});

	test("should return missed author message when close to author", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.GOLD),
			getDifferenceToMedal: jest
				.fn()
				.mockReturnValueOnce(5000) // champion difference
				.mockReturnValueOnce(8), // 8ms from author
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} HAHAHAHAH HE MISSED AUTHOR BY 0.008s ${Emote.OMEGALUL.name} ${Emote.ICANT.name} ${Emote.PEPE_POINT.name}`,
		);
	});

	test("should return gold medal message", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.GOLD),
			getDifferenceToMedal: jest
				.fn()
				.mockReturnValueOnce(5000) // champion difference
				.mockReturnValueOnce(1000), // 1 second from author
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's only ${Emote.GOLD_MEDAL.name}`);
	});

	test("should return silver medal message", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.SILVER),
			getDifferenceToMedal: jest.fn(),
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's only ${Emote.SILVER_MEDAL.name}`);
	});

	test("should return bronze medal message", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.BRONZE),
			getDifferenceToMedal: jest.fn(),
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's only ${Emote.BRONZE_MEDAL.name}`);
	});

	test("should return no medal message", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.NONE),
			getDifferenceToMedal: jest.fn(),
		} as unknown as TrackmaniaMap;

		const newBestTime = new TrackmaniaTime(30000);

		const result = buildPbMapResultMessage(mockMap, newBestTime);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} Not even ${Emote.BRONZE_MEDAL.name} ${Emote.PEPE_POINT.name}`,
		);
	});
});

describe("buildGuessResultMessage", () => {
	test("should combine map result and guess result messages", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.CHAMPION),
			getDifferenceToMedal: jest.fn(),
		} as unknown as TrackmaniaMap;

		const winners: Guess[] = [
			{
				channelId: "test",
				userId: "1",
				displayName: "Alice",
				time: 30000,
			},
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildGuessResultMessage(mockMap, newBestTime, winners);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's Champion Medal ${Emote.CHAMPION_MEDAL.name} ${Emote.BWOAH.name} @Alice the ${Emote.GIGACHAD.name} guessed it correctly! ${Emote.BWOAH.name} (+5 points)`,
		);
	});

	test("should handle null map with winners", () => {
		const winners: Guess[] = [
			{
				channelId: "test",
				userId: "1",
				displayName: "Alice",
				time: 30000,
			},
		];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildGuessResultMessage(null, newBestTime, winners);

		expect(result).toBe(
			`${Emote.YEK.name} I got a new PB but I didn't know you are on a map? @Alice the ${Emote.GIGACHAD.name} guessed it correctly! ${Emote.BWOAH.name} (+5 points)`,
		);
	});

	test("should handle map with no winners", () => {
		const mockMap = {
			getMedal: jest.fn().mockReturnValue(Medal.AUTHOR),
			getDifferenceToMedal: jest.fn().mockReturnValue(1000),
		} as unknown as TrackmaniaMap;

		const winners: Guess[] = [];
		const newBestTime = new TrackmaniaTime(30000);

		const result = buildGuessResultMessage(mockMap, newBestTime, winners);

		expect(result).toBe(
			`NEW PERSONAL BEST 30.000 ${Emote.DINK_DONK.name} That's ${Emote.AUTHOR_MEDAL.name} ${Emote.OK.name} But no chatter participated ${Emote.REALLY_FUCKING_MAD.name}`,
		);
	});
});

describe("evaluateGuesses", () => {
	test("should handle all guesses being equidistant from the result", () => {
		const alice = {
			channelId: "testChannel",
			displayName: "Alice",
			userId: "1",
			time: 1000,
		};
		const bob = {
			channelId: "testChannel",
			displayName: "Bob",
			userId: "2",
			time: 2000,
		};
		const charlie = {
			channelId: "testChannel",
			displayName: "Charlie",
			userId: "3",
			time: 3000,
		};

		const guesses = [alice, bob, charlie];
		const result = new TrackmaniaTime(1500);
		const winners = evaluateGuesses(guesses, result);

		expect(winners.length).toBe(2);
		expect(winners[0]).toEqual(alice);
		expect(winners[1]).toEqual(bob);
	});

	test("should return the single closest guess as the winner", () => {
		const alice = {
			channelId: "testChannel",
			displayName: "Alice",
			userId: "1",
			time: 1000,
		};
		const bob = {
			channelId: "testChannel",
			displayName: "Bob",
			userId: "2",
			time: 2000,
		};
		const charlie = {
			channelId: "testChannel",
			displayName: "Charlie",
			userId: "3",
			time: 3000,
		};

		const guesses = [alice, bob, charlie];
		const result = new TrackmaniaTime(1499);

		const winners = evaluateGuesses(guesses, result);

		expect(winners.length).toBe(1);
		expect(winners[0]).toEqual(alice);
	});

	test("should return multiple winners in case of a tie", () => {
		const alice = {
			channelId: "testChannel",
			displayName: "Alice",
			userId: "1",
			time: 2000,
		};
		const bob = {
			channelId: "testChannel",
			displayName: "Bob",
			userId: "2",
			time: 2000,
		};
		const charlie = {
			channelId: "testChannel",
			displayName: "Charlie",
			userId: "3",
			time: 3000,
		};

		const guesses = [alice, bob, charlie];
		const result = new TrackmaniaTime(2000);

		const winners = evaluateGuesses(guesses, result);

		expect(winners.length).toBe(2);
		expect(winners).toContain(alice);
		expect(winners).toContain(bob);
	});

	test("should handle no guesses gracefully", () => {
		const result = new TrackmaniaTime(1500);
		const winners = evaluateGuesses([], result);
		expect(winners).toEqual([]);
	});

	test("should return the closest guess with large time differences multiple", () => {
		const alice = {
			channelId: "testChannel",
			displayName: "Alice",
			userId: "1",
			time: 3600000, // 1 hour
		};
		const bob = {
			channelId: "testChannel",
			displayName: "Bob",
			userId: "2",
			time: 7200000, // 2 hours
		};

		const result = new TrackmaniaTime(1.5 * 60 * 60 * 1000); // 1.5 hours
		const guesses = [alice, bob];
		const winners = evaluateGuesses(guesses, result);

		expect(winners.length).toBe(2);
		expect(winners[0]).toEqual(alice);
		expect(winners[1]).toEqual(bob);
	});

	test("should return the closest guess with large time differences single", () => {
		const alice = {
			channelId: "testChannel",
			displayName: "Alice",
			userId: "1",
			time: 3600000, // 1 hour
		};
		const bob = {
			channelId: "testChannel",
			displayName: "Bob",
			userId: "2",
			time: 7200000, // 2 hours
		};

		const result = new TrackmaniaTime(1.5 * 60 * 60 * 1000 + 1); // 1.5 hours and 1 ms
		const guesses = [alice, bob];
		const winners = evaluateGuesses(guesses, result);

		expect(winners.length).toBe(1);
		expect(winners[0]).toEqual(bob);
	});

	test("should handle scenarios with multiple players and large numbers", () => {
		const alice = {
			channelId: "testChannel",
			displayName: "Alice",
			userId: "1",
			time: 100000, // 100 seconds
		};
		const bob = {
			channelId: "testChannel",
			displayName: "Bob",
			userId: "2",
			time: 200000, // 200 seconds
		};
		const charlie = {
			channelId: "testChannel",
			displayName: "Charlie",
			userId: "3",
			time: 150000, // 150 seconds
		};
		const dave = {
			channelId: "testChannel",
			displayName: "Dave",
			userId: "4",
			time: 150000, // 150 seconds
		};

		const guesses = [alice, bob, charlie, dave];
		const result = new TrackmaniaTime(150000);
		const winners = evaluateGuesses(guesses, result);

		expect(winners.length).toBe(2);
		expect(winners).toContain(charlie);
		expect(winners).toContain(dave);
	});
});
