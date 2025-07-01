import { database } from "../../database";
import { Command } from "../core";

const LEADERBOARD_LIMIT = 10;

export class LeaderboardCommand extends Command {
    protected async onExecute(): Promise<string | null> {
        const leaderboard = await database.guesserLeaderboard.findMany({
            orderBy: {
                points: "desc",
            },
            take: LEADERBOARD_LIMIT,
        });

        if (leaderboard.length === 0) {
            return "No guesses have been made yet.";
        }

        const leaderboardMessage = leaderboard
            .map((entry, index) => `${index + 1}. ${entry.displayName} - ${entry.points}`)
            .join(" | ");
        return `Top ${LEADERBOARD_LIMIT} Guessers 🏆 Perfect guess = 5pts, Closest guess = 1pt: ${leaderboardMessage}`;
    }
}   