import type { Guess } from "../../../generated/prisma";
import { Emote } from "./emotes";
import { Medal, type TrackmaniaMap } from "./map";
import { TrackmaniaTime } from "./time";

export const evaluateGuesses = (guesses: Guess[], result: TrackmaniaTime) => {
	let minDifference: number | null = null;
	let winners: Guess[] = [];
	for (const guess of guesses) {
		const time = new TrackmaniaTime(guess.time);
		const difference = Math.abs(result.getDifference(time));
		if (minDifference === null || difference < minDifference) {
			minDifference = difference;
			winners = [guess];
			continue;
		}
		if (difference === minDifference) {
			winners.push(guess);
		}
	}
	return winners;
};

export const buildGuessResultMessage = (
	currentMap: TrackmaniaMap | null,
	newBestTime: TrackmaniaTime,
	winners: Guess[],
) => {
	const message = buildPbMapResultMessage(currentMap, newBestTime);
	const guessResult = buildBestGuesserMessage(winners, newBestTime);
	return `${message} ${guessResult}`;
};

const buildBestGuesserMessage = (winners: Guess[], newBestTime: TrackmaniaTime) => {
	if (winners.length === 0) {
		return `But no chatter participated ${Emote.REALLY_FUCKING_MAP.name}`;
	}

	const bestGuess = winners[0];
	const bestGuessAsString = bestGuess.time.toString();
	const bestGuessTime = new TrackmaniaTime(bestGuess.time);

	const difference = newBestTime.getDifference(bestGuessTime);
	const differenceAsString = `${difference < 0 ? "+" : "-"} ${new TrackmaniaTime(Math.abs(difference)).toString()}`;
	const hasPerfectGuess = difference === 0;

	if (winners.length === 0) {
		const username = bestGuess.displayName;
		return hasPerfectGuess
			? `@${username} the ${Emote.GIGACHAD.name} guessed it correctly! ${Emote.BWOAH.name}`
			: `Nobody guessed it correctly but @${username} guessed ${bestGuessAsString} (${differenceAsString}) ${Emote.OK.name}`;
	}

	const usernames = winners.map((winner) => winner.displayName).join(", ");
	return hasPerfectGuess
		? `@${usernames} the ${Emote.GIGACHAD.name} 's guessed it correctly! ${Emote.BWOAH.name}`
		: `Nobody guessed it correctly but @${usernames} guessed ${bestGuessAsString} (${differenceAsString}) ${Emote.OK.name}`;
};

const buildPbMapResultMessage = (currentMap: TrackmaniaMap | null, newBestTime: TrackmaniaTime) => {
	if (currentMap === null) return `${Emote.YEK.name} I got a new PB but I didn't know you are on a map?`;

	const medal = currentMap.getMedal(newBestTime);

	const differenceToChampion = currentMap.getDifferenceToMedal(Medal.CHAMPION, newBestTime);
	const differenceToAuthor = currentMap.getDifferenceToMedal(Medal.AUTHOR, newBestTime);

	switch (medal) {
		case Medal.CHAMPION:
			return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} That's Champion Medal ${Emote.CHAMPION_MEDAL.name} ${Emote.BWOAH.name}`;
		case Medal.AUTHOR:
			if (differenceToChampion < 10) {
				const differenceAsTime = new TrackmaniaTime(Math.abs(differenceToChampion)).toString();
				return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} HAHAHAHAH HE MISSED CHAMPION BY ${differenceAsTime}s ${Emote.OMEGALUL.name} ${Emote.ICANT.name} ${Emote.PEPE_POINT.name}`;
			}
			return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} That's ${Emote.AUTHOR_MEDAL.name} ${Emote.OK.name}`;
		case Medal.GOLD:
			if (differenceToAuthor < 10) {
				const differenceAsTime = new TrackmaniaTime(Math.abs(differenceToAuthor)).toString();
				return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} HAHAHAHAH HE MISSED AUTHOR BY ${differenceAsTime}s ${Emote.OMEGALUL.name} ${Emote.ICANT.name} ${Emote.PEPE_POINT.name}`;
			}
			return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} That's only ${Emote.GOLD_MEDAL.name}`;
		case Medal.SILVER:
			return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} That's only ${Emote.SILVER_MEDAL.name}`;
		case Medal.BRONZE:
			return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} That's only ${Emote.BRONZE_MEDAL.name}`;
		case Medal.NONE:
			return `NEW PERSONAL BEST ${newBestTime.toString()} ${Emote.DINK_DONK.name} Not even §{Emote.BRO} ${Emote.PEPE_POINT.name}`;
	}
	return "";
};
