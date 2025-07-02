import type { Guess } from "@prisma/client";
import { database } from "../../database";
import { Emote } from "./emotes";
import { Medal, type TrackmaniaMap } from "./map";
import { TrackmaniaTime } from "./time";
import { mentionUser } from "./user";

export const POINTS_PERFECT_GUESS = 5;
export const POINTS_CLOSEST_GUESS = 1;

export const updateGuessResultLeaderboard = async (winners: Guess[], result: TrackmaniaTime) => {
	if (winners.length === 0) {
		return;
	}

	const perfectGuess = winners[0].time === result.getTotalInMilliSeconds();
	const points = perfectGuess ? POINTS_PERFECT_GUESS : POINTS_CLOSEST_GUESS;

	for (const winner of winners) {
		await database.guesserLeaderboard.upsert({
			where: {
				channelId_userId: {
					channelId: winner.channelId,
					userId: winner.userId,
				},
			},
			update: {
				points: {
					increment: points,
				},
				perfectGuessCount: perfectGuess ? { increment: 1 } : undefined,
				displayName: winner.displayName,
			},
			create: {
				channelId: winner.channelId,
				userId: winner.userId,
				displayName: winner.displayName,
				points,
				perfectGuessCount: perfectGuess ? 1 : 0,
			},
		});
	}
};

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
		return `But no chatter participated ${Emote.REALLY_FUCKING_MAD.name}`;
	}

	const bestGuess = winners[0];
	const bestGuessTime = new TrackmaniaTime(bestGuess.time);
	const bestGuessAsString = bestGuessTime.toString();

	const difference = newBestTime.getDifference(bestGuessTime);
	const differenceAsString = `${difference < 0 ? "+" : "-"}${new TrackmaniaTime(Math.abs(difference)).toString()}`;
	const hasPerfectGuess = difference === 0;

	if (winners.length === 1) {
		const username = bestGuess.displayName;
		return hasPerfectGuess
			? `${mentionUser(username)} the ${Emote.GIGACHAD.name} guessed it correctly! ${Emote.BWOAH.name} (+${POINTS_PERFECT_GUESS} points)`
			: `Nobody guessed it correctly but ${mentionUser(username)} guessed ${bestGuessAsString} (${differenceAsString}) ${Emote.OK.name} (+${POINTS_CLOSEST_GUESS} point)`;
	}

	const usernames = winners.map((winner) => mentionUser(winner.displayName)).join(", ");
	return hasPerfectGuess
		? `${usernames} the ${Emote.GIGACHAD.name} 's guessed it correctly! ${Emote.BWOAH.name} (+${POINTS_PERFECT_GUESS} points)`
		: `Nobody guessed it correctly but ${mentionUser(usernames)} guessed ${bestGuessAsString} (${differenceAsString}) ${Emote.OK.name} (+${POINTS_CLOSEST_GUESS} point)`;
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
