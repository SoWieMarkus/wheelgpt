import { getLeaderboard, getLeaderboardByName } from "@prisma/client/sql";
import { database } from "../../database";
import { Command, mentionUser, type User } from "../core";

const LEADERBOARD_LIMIT = 5;

export class LeaderboardCommand extends Command {
    protected async onExecute(user: User, args: string[]): Promise<string | null> {
        if (args.length > 0) {
            const userToFind = args[0].replace(/^@/, "").toLowerCase();
            const entry = await database.$queryRawTyped(getLeaderboardByName(this.channelId, userToFind));
            if (entry.length === 0) {
                return `${mentionUser(user.displayName)} No leaderboard entry found for ${userToFind}.`;
            }
            const { displayName, position, points } = entry[0];
            const pointsText = points === 1 ? "pt" : "pts";
            return `${mentionUser(displayName)} is currently in position ${position} with ${points} ${pointsText}.`;
        }

        const leaderboard = await database.$queryRawTyped(getLeaderboard(this.channelId, LEADERBOARD_LIMIT, 0));

        if (leaderboard.length === 0) {
            return "No guesses have been made yet.";
        }
        const leaderboardUrl = `https://wheelgpt.dev/leaderboard/${this.channelId}`;
        const leaderboardMessage = leaderboard
            .map((entry) => `${entry.position}. ${entry.displayName} - ${entry.points}`)
            .join(" | ");
        return `Top ${LEADERBOARD_LIMIT} Guessers 🏆 Perfect guess = 5pts, Closest guess = 1pt: ${leaderboardMessage} | You can view the full leaderboard here: ${leaderboardUrl}`;
    }
}
