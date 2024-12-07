import { z } from "zod";
import { TrackmaniaTime } from "../guess/time";
import { assertIsDefined } from "../../utils";

export enum Medal {
    CHAMPION = 0,
    AUTHOR = 1,
    GOLD = 2,
    SILVER = 3,
    BRONZE = 4,
    NONE = 5,
}

export const TrackmaniaMapSchema = z.object({
    name: z.string(),
    uid: z.string(),
    author: z.string(),
    authorTime: z.number().min(0),
    goldTime: z.number().min(0),
    silverTime: z.number().min(0),
    bronzeTime: z.number().min(0),
    championTime: z.number().min(0),
    tmxId: z.string().optional(),
});

export const TrackmaniaMapPostSchema = z.object({
    mapAvailable: z.boolean(),
    map: TrackmaniaMapSchema.optional(),
});

const stylingChars = [
    "w",
    "n",
    "o",
    "i",
    "t",
    "s",
    "g",
    "z",
    "$",
    "W",
    "N",
    "O",
    "I",
    "T",
    "S",
    "G",
    "Z",
];

const stylingRegex = new RegExp(
    `[$]([${stylingChars.join("")}]|[0-9A-Fa-f]{3})`,
    "g",
);

export class TrackmaniaMap {
    constructor(public readonly data: z.infer<typeof TrackmaniaMapSchema>) {
        this.data.name = this.data.name.replace(stylingRegex, "").trim();
    }

    public hasChampionTime() {
        return this.data.championTime > 0;
    }

    public hasTrackmaniaExchangeId() {
        return this.data.tmxId !== undefined;
    }

    public setTrackmaniaExchangeId(id: string | undefined) {
        this.data.tmxId = id;
    }

    public getMedal(value: number | TrackmaniaTime) {
        const time =
            typeof value === "number" ? value : value.getTotalInMilliSeconds();

        const { championTime, authorTime, goldTime, silverTime, bronzeTime } =
            this.data;

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
        const time =
            typeof value === "number" ? value : value.getTotalInMilliSeconds();
        const { championTime, authorTime, goldTime, silverTime, bronzeTime } =
            this.data;

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
        const { championTime, authorTime, goldTime, silverTime, bronzeTime } =
            this.data;

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

    public getTrackmaniaExchangeLink() {
        const { tmxId } = this.data;
        if (tmxId === undefined) return undefined;
        return `https://trackmania.exchange/maps/${tmxId}`;
    }

    public toString() {
        const hasChampionTime = this.hasChampionTime();
        const hasTrackmaniaExchangeId = this.hasTrackmaniaExchangeId();

        const { author, name } = this.data;

        const tmx = this.getTrackmaniaExchangeLink();
        const champion = this.getMedalTime(Medal.CHAMPION).toString();
        const authorMedal = this.getMedalTime(Medal.AUTHOR).toString();

        if (hasTrackmaniaExchangeId && hasChampionTime) {
            assertIsDefined(tmx, "TrackmaniaExchange Link");
            return `"${name}" by ${author} | champion_medal  ${champion} | AuthorTime  ${authorMedal} | ${tmx}`;
        }

        if (!hasTrackmaniaExchangeId && hasChampionTime) {
            return `"${name}" by ${author} | champion_medal  ${champion} | AuthorTime  ${authorMedal}`;
        }

        if (hasTrackmaniaExchangeId && !hasChampionTime) {
            assertIsDefined(tmx, "TrackmaniaExchange Link");
            return `"${name}" by ${author} | AuthorTime  ${authorMedal} | ${tmx}`;
        }

        return `"${name}" by ${author} | AuthorTime  ${authorMedal}`
    }
}