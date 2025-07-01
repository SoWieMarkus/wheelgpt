import type * as Table from "@prisma/client";
import { assertIsDefined } from "./assert-defined";
import { Emote } from "./emotes";
import { TrackmaniaTime } from "./time";

export enum Medal {
	CHAMPION = 0,
	AUTHOR = 1,
	GOLD = 2,
	SILVER = 3,
	BRONZE = 4,
	NONE = 5,
}

const stylingChars = ["w", "n", "o", "i", "t", "s", "g", "z", "$", "W", "N", "O", "I", "T", "S", "G", "Z"];
const stylingRegex = new RegExp(`[$]([${stylingChars.join("")}]|[0-9A-Fa-f]{3})`, "g");

export class TrackmaniaMap {
	constructor(public readonly data: Table.TrackmaniaMap) {
		this.data.name = this.data.name.replace(stylingRegex, "").trim();
	}

	public hasChampionTime() {
		return this.data.championTime > 0;
	}

	public hasTrackmaniaExchangeId() {
		return this.data.tmxId !== null;
	}

	public getMedal(value: number | TrackmaniaTime) {
		const time = typeof value === "number" ? value : value.getTotalInMilliSeconds();

		const { championTime, authorTime, goldTime, silverTime, bronzeTime } = this.data;

		if (this.hasChampionTime() && time <= championTime) {
			return Medal.CHAMPION;
		}

		if (time <= authorTime) {
			return Medal.AUTHOR;
		}

		if (time <= goldTime) {
			return Medal.GOLD;
		}

		if (time <= silverTime) {
			return Medal.SILVER;
		}

		if (time <= bronzeTime) {
			return Medal.BRONZE;
		}

		return Medal.NONE;
	}

	public getDifferenceToMedal(medal: Medal, value: number | TrackmaniaTime) {
		const time = typeof value === "number" ? value : value.getTotalInMilliSeconds();
		const { championTime, authorTime, goldTime, silverTime, bronzeTime } = this.data;

		switch (medal) {
			case Medal.CHAMPION:
				return time - championTime;
			case Medal.AUTHOR:
				return time - authorTime;
			case Medal.GOLD:
				return time - goldTime;
			case Medal.SILVER:
				return time - silverTime;
			case Medal.BRONZE:
				return time - bronzeTime;
			default:
				throw new Error("Not possible to get the difference to this medal.");
		}
	}

	public getMedalTime(medal: Medal) {
		const { championTime, authorTime, goldTime, silverTime, bronzeTime } = this.data;

		switch (medal) {
			case Medal.CHAMPION:
				return new TrackmaniaTime(championTime);
			case Medal.AUTHOR:
				return new TrackmaniaTime(authorTime);
			case Medal.GOLD:
				return new TrackmaniaTime(goldTime);
			case Medal.SILVER:
				return new TrackmaniaTime(silverTime);
			case Medal.BRONZE:
				return new TrackmaniaTime(bronzeTime);
			default:
				throw new Error("Not possible to get the time of this medal.");
		}
	}

	public get tmxLink() {
		const { tmxId } = this.data;
		if (tmxId === null) return undefined;
		return `https://trackmania.exchange/maps/${tmxId}`;
	}

	public get trackmaniaIOLink() {
		const { uid } = this.data;
		if (uid === null) return undefined;
		return `https://trackmania.io/#/leaderboard/${uid}`;
	}

	public toString() {
		const hasChampionTime = this.hasChampionTime();
		const hasTrackmaniaExchangeId = this.hasTrackmaniaExchangeId();

		const { author, name } = this.data;

		const tmx = this.tmxLink;
		const champion = this.getMedalTime(Medal.CHAMPION).toString();
		const authorMedal = this.getMedalTime(Medal.AUTHOR).toString();

		if (hasTrackmaniaExchangeId && hasChampionTime) {
			assertIsDefined(tmx, "TrackmaniaExchange Link");
			return `"${name}" by ${author} | ${Emote.CHAMPION_MEDAL.name}  ${champion} | ${Emote.AUTHOR_MEDAL.name}  ${authorMedal} | ${tmx} | ${this.trackmaniaIOLink}`;
		}

		if (!hasTrackmaniaExchangeId && hasChampionTime) {
			return `"${name}" by ${author} | ${Emote.CHAMPION_MEDAL.name}  ${champion} | ${Emote.AUTHOR_MEDAL.name}  ${authorMedal} | ${this.trackmaniaIOLink}`;
		}

		if (hasTrackmaniaExchangeId && !hasChampionTime) {
			assertIsDefined(tmx, "TrackmaniaExchange Link");
			return `"${name}" by ${author} | ${Emote.AUTHOR_MEDAL.name}  ${authorMedal} | ${tmx} | ${this.trackmaniaIOLink}`;
		}

		return `"${name}" by ${author} | ${Emote.AUTHOR_MEDAL.name}  ${authorMedal} | ${this.trackmaniaIOLink}`;
	}
}
