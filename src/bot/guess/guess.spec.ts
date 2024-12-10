import { evaluateGuesses, type Guess } from './guess'; // Adjust the path to your module file
import { TrackmaniaTime } from './time';
import type { User } from '../user';

describe('evaluateGuesses', () => {
    const createUser = (id: string, name: string): User => ({
        id,
        name,
        displayName: name,
        accessLevel: 1,
        channelId: 'testChannel',
    });

    const createGuess = (user: User, timeInMilliseconds: number): Guess => ({
        user,
        time: new TrackmaniaTime(timeInMilliseconds),
    });

    test('should return the single closest guess as the winner multiple', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 1000)], // 1 second
            ['player2', createGuess(createUser('2', 'Bob'), 2000)], // 2 seconds
            ['player3', createGuess(createUser('3', 'Charlie'), 3000)], // 3 seconds
        ]);
        const result = new TrackmaniaTime(1500); // 1.5 seconds

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(2);
        expect(winners[0]).toEqual(guesses.get('player1'));
        expect(winners[1]).toEqual(guesses.get('player2'));
    });

    test('should return the single closest guess as the winner single', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 1000)], // 1 second
            ['player2', createGuess(createUser('2', 'Bob'), 2000)], // 2 seconds
            ['player3', createGuess(createUser('3', 'Charlie'), 3000)], // 3 seconds
        ]);
        const result = new TrackmaniaTime(1499); // 1.5 seconds

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(1);
        expect(winners[0]).toEqual(guesses.get('player1'));
    });

    test('should return multiple winners in case of a tie', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 2000)],
            ['player2', createGuess(createUser('2', 'Bob'), 2000)],
            ['player3', createGuess(createUser('3', 'Charlie'), 3000)],
        ]);
        const result = new TrackmaniaTime(2000);

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(2);
        expect(winners).toContain(guesses.get('player1'));
        expect(winners).toContain(guesses.get('player2'));
    });

    test('should handle no guesses gracefully', () => {
        const guesses = new Map<string, Guess>();
        const result = new TrackmaniaTime(1500);

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toEqual([]);
    });

    test('should handle all guesses being equidistant from the result', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 1000)],
            ['player2', createGuess(createUser('2', 'Bob'), 2000)],
        ]);
        const result = new TrackmaniaTime(1500); // Both guesses are 500ms away

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(2);
        expect(winners).toContain(guesses.get('player1'));
        expect(winners).toContain(guesses.get('player2'));
    });

    test('should return the closest guess with large time differences multiple', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 3600000)], // 1 hour
            ['player2', createGuess(createUser('2', 'Bob'), 7200000)], // 2 hours
        ]);
        const result = new TrackmaniaTime(5400000); // 1.5 hours

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(2);
        expect(winners[0]).toEqual(guesses.get('player1'));
        expect(winners[1]).toEqual(guesses.get('player2'));
    });

    test('should return the closest guess with large time differences single', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 60 * 60 * 1000)], // 1 hour
            ['player2', createGuess(createUser('2', 'Bob'), 2 * 60 * 60 * 1000)], // 2 hours
        ]);
        const result = new TrackmaniaTime(1.5 * 60 * 60 * 1000 + 1); // 1.5 hours

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(1);
        expect(winners[0]).toEqual(guesses.get('player2'));
    });

    test('should handle scenarios with multiple players and large numbers', () => {
        const guesses = new Map<string, Guess>([
            ['player1', createGuess(createUser('1', 'Alice'), 100000)],
            ['player2', createGuess(createUser('2', 'Bob'), 200000)],
            ['player3', createGuess(createUser('3', 'Charlie'), 150000)],
            ['player4', createGuess(createUser('4', 'Dave'), 150000)],
        ]);
        const result = new TrackmaniaTime(150000);

        const winners = evaluateGuesses(guesses, result);

        expect(winners).toHaveLength(2);
        expect(winners).toContain(guesses.get('player3'));
        expect(winners).toContain(guesses.get('player4'));
    });
});
