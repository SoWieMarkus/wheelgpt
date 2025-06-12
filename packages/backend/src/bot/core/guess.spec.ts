import { evaluateGuesses } from "./guess";
import { TrackmaniaTime } from "./time";

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
