import { TrackmaniaTime, parseTrackmaniaTime } from "./time";

describe("TrackmaniaTime class", () => {
    test("constructor should correctly parse time in milliseconds", () => {
        const time = new TrackmaniaTime(3723456);
        expect(time.toString()).toBe("1:02:03.456");
    });

    test("getTotalInMilliSeconds should return correct total milliseconds", () => {
        const time = new TrackmaniaTime(3723456);
        expect(time.getTotalInMilliSeconds()).toBe(3723456);
    });

    test("getDifference should return the difference in milliseconds between two TrackmaniaTime objects", () => {
        const time1 = new TrackmaniaTime(3723456);
        const time2 = new TrackmaniaTime(3600000);
        expect(time1.getDifference(time2)).toBe(123456);
    });

    test("toString should format time correctly", () => {
        const time1 = new TrackmaniaTime(123456);
        const time2 = new TrackmaniaTime(60000);
        const time3 = new TrackmaniaTime(3456);
        expect(time1.toString()).toBe("2:03.456");
        expect(time2.toString()).toBe("1:00.000");
        expect(time3.toString()).toBe("3.456");
    });
});

const validStrings = [
    "0.0",
    "0.001",
    "00:00:00.000",
    "100:100:12.3",
    "100:00.0",
    "1000:00:00.034",
    "1000:00.3",
    "1000.12",
];
const invalidStrings = [
    "1797693134862315700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001:34.2",
    "0.",
    ":0.1",
    "0.0001",
    "0,0",
    "10:10:10",
    ".234",
    "",
    "Hallo",
    "Yo mama",
    "1",
    ":10:0.0",
    ":10:10:10.0",
    "10:10:10:10:10.0",
    "10::10.0",
];

describe("parseTrackmaniaTime function", () => {
    for (const validString of validStrings) {
        test(`${validString} should be a valid string to parse.`, () => {
            expect(parseTrackmaniaTime(validString)).not.toBeNull();
        });
    }

    for (const invalidString of invalidStrings) {
        test(`${invalidString} should not be a valid string to parse.`, () => {
            expect(parseTrackmaniaTime(invalidString)).toBeNull();
        });
    }

    test("should correctly parse time string with hours, minutes, seconds and milliseconds", () => {
        const time = parseTrackmaniaTime("1:02:03.456");
        expect(time).not.toBeNull();
        expect(time?.getTotalInMilliSeconds()).toBe(3723456);
    });

    test("should correctly parse time string with minutes, seconds and milliseconds", () => {
        const time = parseTrackmaniaTime("2:03.456");
        expect(time).not.toBeNull();
        expect(time?.getTotalInMilliSeconds()).toBe(123456);
    });

    test("should correctly parse time string with seconds and milliseconds", () => {
        const time = parseTrackmaniaTime("3.456");
        expect(time).not.toBeNull();
        expect(time?.getTotalInMilliSeconds()).toBe(3456);
    });

    test("should correctly parse time string without milliseconds", () => {
        const time1 = parseTrackmaniaTime("1:02:03.000");
        const time2 = parseTrackmaniaTime("2:03.000");
        const time3 = parseTrackmaniaTime("3.000");
        expect(time1?.getTotalInMilliSeconds()).toBe(3723000);
        expect(time2?.getTotalInMilliSeconds()).toBe(123000);
        expect(time3?.getTotalInMilliSeconds()).toBe(3000);
    });

    test("should pad milliseconds if less than 3 digits", () => {
        const time = parseTrackmaniaTime("1:02:03.4");
        expect(time?.getTotalInMilliSeconds()).toBe(3723400);
    });
});