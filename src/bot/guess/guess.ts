import { TrackmaniaTime } from "./time";
import type { User } from "../user";
import { Medal, type TrackmaniaMap } from "../map/map";

export type Guess = {
    user: User;
    time: TrackmaniaTime;
};

export const evaluateGuesses = (
    guesses: Map<string, Guess>,
    result: TrackmaniaTime,
) => {

    let minDifference: number | null = null;

    const winners: Guess[] = []
    for (const guess of guesses.values()) {
        const difference = result.getDifference(guess.time);
        if (minDifference === null || difference < minDifference) {
            minDifference = difference;
            winners.length = 0;
            winners.push(guess);
            continue;
        }
        if (difference === minDifference) {
            winners.push(guess);
        }
    }

    return winners;
};

export const buildGuessResultMessage = (currentMap: TrackmaniaMap | null, newBestTime: TrackmaniaTime, winners: Guess[]) => {
    const message = buildPbMapResultMessage(currentMap, newBestTime);
    const guessResult = buildBestGuesserMessage(winners, newBestTime);
    return `${message} ${guessResult}`;
}

const buildBestGuesserMessage = (winners: Guess[], newBestTime: TrackmaniaTime) => {
    if (winners.length === 0) {
        return "But no chatter participated ReallyFuckingMad";
    }

    const bestGuess = winners[0];
    const bestGuessAsString = bestGuess.time.toString();
    const difference = newBestTime.getDifference(bestGuess.time);
    const differenceAsString = `${difference < 0 ? "+" : "-"} ${new TrackmaniaTime(Math.abs(difference)).toString()}`;
    const hasPerfectGuess = difference === 0;

    if (winners.length === 0) {
        const username = bestGuess.user.displayName
        return hasPerfectGuess
            ? `@${username} the GIGACHAD  guessed it correctly! BWOAH`
            : `Nobody guessed it correctly but @${username} guessed ${bestGuessAsString} (${differenceAsString}) ok`
    }


    const usernames = winners.map(winner => winner.user.displayName).join(", ");
    return hasPerfectGuess
        ? `@${usernames} the GIGACHAD 's  guessed it correctly! BWOAH`
        : `Nobody guessed it correctly but @${usernames} guessed ${bestGuessAsString} (${differenceAsString}) ok`

}

const buildPbMapResultMessage = (currentMap: TrackmaniaMap | null, newBestTime: TrackmaniaTime) => {
    if (currentMap === null) return "YEK I got a new PB but I didn't know you are on a map?"


    const medal = currentMap.getMedal(newBestTime);

    const differenceToChampion = currentMap.getDifferenceToMedal(Medal.CHAMPION, newBestTime);
    const differenceToAuthor = currentMap.getDifferenceToMedal(Medal.AUTHOR, newBestTime);


    switch (medal) {
        case Medal.CHAMPION:
            return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk That's Champion Medal champion_medal BWOAH`
        case Medal.AUTHOR:
            if (differenceToChampion < 10) {
                const differenceAsTime = new TrackmaniaTime(Math.abs(differenceToChampion)).toString();
                return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk HAHAHAHAH HE MISSED CHAMPION BY ${differenceAsTime}s OMEGALUL ICANT PepePoint`;
            }
            return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk That's AuthorTime ok `
        case Medal.GOLD:
            if (differenceToAuthor < 10) {
                const differenceAsTime = new TrackmaniaTime(Math.abs(differenceToAuthor)).toString();
                return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk HAHAHAHAH HE MISSED AUTHOR BY ${differenceAsTime}s OMEGALUL ICANT PepePoint`;
            }
            return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk That's only GLOD`
        case Medal.SILVER:
            return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk That's only SLIVER`
        case Medal.BRONZE:
            return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk That's only BORNZE`
        case Medal.NONE:
            return `NEW PERSONAL BEST ${newBestTime.toString()} dinkDonk Not even BORNZE PepePoint`
    }
    return "";
}