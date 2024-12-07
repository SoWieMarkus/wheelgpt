import { z } from "zod";

export const GuessResultSchema = z.object({
    time: z.number().int().min(0),
});

export class TrackmaniaTime {
    private readonly hours: number;
    private readonly minutes: number;
    private readonly seconds: number;
    private readonly milliseconds: number;

    constructor(timeInMilliseconds: number) {
        this.milliseconds = Math.floor(timeInMilliseconds % 1000);
        this.seconds = Math.floor((timeInMilliseconds / 1000) % 60);
        this.minutes = Math.floor((timeInMilliseconds / (1000 * 60)) % 60);
        this.hours = Math.floor(timeInMilliseconds / (1000 * 60 * 60));
    }

    public getDifference(time: TrackmaniaTime) {
        const thisTime = this.getTotalInMilliSeconds();
        const compareTime = time.getTotalInMilliSeconds();
        return thisTime - compareTime;
    }

    public getTotalInMilliSeconds() {
        return (
            this.milliseconds +
            1000 * this.seconds +
            1000 * 60 * this.minutes +
            1000 * 60 * 60 * this.hours
        );
    }

    public toString() {
        const padZero = (num: number, size: number) =>
            num.toString().padStart(size, "0");

        const milliseconds = padZero(this.milliseconds, 3);
        const seconds = padZero(this.seconds, 2);
        const minutes = padZero(this.minutes, 2);

        if (this.hours > 0)
            return `${this.hours}:${minutes}:${seconds}.${milliseconds}`;
        if (this.minutes > 0) return `${this.minutes}:${seconds}.${milliseconds}`;
        return `${this.seconds}.${milliseconds}`;
    }
}

const regex = /^(?:(\d{1,5}):)?(?:(\d{1,5}):)?(\d{1,5})(?:\.(\d{1,3}))$/;

export const parseTrackmaniaTime = (value: string) => {
    const match = regex.test(value);
    if (!match) return null;

    const parts = value.split(".");
    const millisecondsString = parts[1]?.padEnd(3, "0") || "000";
    const milliseconds = Number(millisecondsString);

    const hhmmss = parts[0].split(":").reverse().map(Number);

    const seconds = hhmmss[0] || 0;
    const minutes = hhmmss[1] || 0;
    const hours = hhmmss[2] || 0;

    const totalMilliseconds =
        hours * 3600000 + // 1000 * 60 * 60
        minutes * 60000 + // 1000 * 60
        seconds * 1000 +
        milliseconds;

    return new TrackmaniaTime(totalMilliseconds);
};