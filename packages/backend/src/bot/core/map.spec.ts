import { TrackmaniaTime } from "../guess/time";
import { Medal, TrackmaniaMap } from "./map";

describe("TrackmaniaMap", () => {
	const validMapData = {
		name: "MapName$w",
		uid: "abc123",
		author: "AuthorName",
		authorTime: 1000,
		goldTime: 2000,
		silverTime: 3000,
		bronzeTime: 4000,
		championTime: 500,
		tmxId: 12345,
	};

	const noChampionMapData = {
		name: "MapWithoutChampion",
		uid: "def456",
		author: "AnotherAuthor",
		authorTime: 1000,
		goldTime: 2000,
		silverTime: 3000,
		bronzeTime: 4000,
		championTime: 0,
		tmxId: 54321,
	};

	it("should clean the map name of styling characters", () => {
		const map = new TrackmaniaMap(validMapData);
		expect(map.data.name).toBe("MapName");
	});

	it("should detect if a map has a champion time", () => {
		const mapWithChampion = new TrackmaniaMap(validMapData);
		const mapWithoutChampion = new TrackmaniaMap(noChampionMapData);

		expect(mapWithChampion.hasChampionTime()).toBe(true);
		expect(mapWithoutChampion.hasChampionTime()).toBe(false);
	});

	it("should detect if a map has a Trackmania Exchange ID", () => {
		const mapWithTmx = new TrackmaniaMap(validMapData);
		const mapWithoutTmx = new TrackmaniaMap({
			...validMapData,
			tmxId: undefined,
		});

		expect(mapWithTmx.hasTrackmaniaExchangeId()).toBe(true);
		expect(mapWithoutTmx.hasTrackmaniaExchangeId()).toBe(false);
	});

	it("should set a new Trackmania Exchange ID", () => {
		const map = new TrackmaniaMap(validMapData);
		map.setTrackmaniaExchangeId(67890);
		expect(map.data.tmxId).toBe(67890);
	});

	it("should return the correct Trackmania Exchange link", () => {
		const mapWithTmx = new TrackmaniaMap(validMapData);
		mapWithTmx.setTrackmaniaExchangeId(67890);
		const mapWithoutTmx = new TrackmaniaMap({
			...validMapData,
			tmxId: undefined,
		});

		expect(mapWithTmx.getTrackmaniaExchangeLink()).toBe("https://trackmania.exchange/maps/67890");
		expect(mapWithoutTmx.getTrackmaniaExchangeLink()).toBeUndefined();
	});

	it("should return the correct medal based on the time", () => {
		const map = new TrackmaniaMap(validMapData);

		expect(map.getMedal(450)).toBe(Medal.CHAMPION);
		expect(map.getMedal(500)).toBe(Medal.CHAMPION);
		expect(map.getMedal(999)).toBe(Medal.AUTHOR);
		expect(map.getMedal(1000)).toBe(Medal.AUTHOR);
		expect(map.getMedal(1001)).toBe(Medal.GOLD);
		expect(map.getMedal(1500)).toBe(Medal.GOLD);
		expect(map.getMedal(1999)).toBe(Medal.GOLD);
		expect(map.getMedal(2000)).toBe(Medal.GOLD);
		expect(map.getMedal(2001)).toBe(Medal.SILVER);
		expect(map.getMedal(2500)).toBe(Medal.SILVER);
		expect(map.getMedal(2999)).toBe(Medal.SILVER);
		expect(map.getMedal(3000)).toBe(Medal.SILVER);
		expect(map.getMedal(3001)).toBe(Medal.BRONZE);
		expect(map.getMedal(3500)).toBe(Medal.BRONZE);
		expect(map.getMedal(3999)).toBe(Medal.BRONZE);
		expect(map.getMedal(4000)).toBe(Medal.BRONZE);
		expect(map.getMedal(4001)).toBe(Medal.NONE);
		expect(map.getMedal(6000)).toBe(Medal.NONE);

		const mapWithoutChampion = new TrackmaniaMap(noChampionMapData);

		expect(mapWithoutChampion.getMedal(450)).toBe(Medal.AUTHOR);
		expect(mapWithoutChampion.getMedal(500)).toBe(Medal.AUTHOR);
		expect(mapWithoutChampion.getMedal(1000)).toBe(Medal.AUTHOR);
		expect(mapWithoutChampion.getMedal(1001)).toBe(Medal.GOLD);
	});

	it("should return the correct difference to a given medal", () => {
		const map = new TrackmaniaMap(validMapData);

		expect(map.getDifferenceToMedal(Medal.CHAMPION, 450)).toBe(-50);
		expect(map.getDifferenceToMedal(Medal.AUTHOR, 1200)).toBe(200);
		expect(map.getDifferenceToMedal(Medal.GOLD, 1500)).toBe(-500);
	});

	it("should return the correct medal time", () => {
		const map = new TrackmaniaMap(validMapData);

		expect(map.getMedalTime(Medal.CHAMPION)).toEqual(new TrackmaniaTime(500));
		expect(map.getMedalTime(Medal.AUTHOR)).toEqual(new TrackmaniaTime(1000));
		expect(map.getMedalTime(Medal.GOLD)).toEqual(new TrackmaniaTime(2000));
		expect(map.getMedalTime(Medal.SILVER)).toEqual(new TrackmaniaTime(3000));
		expect(map.getMedalTime(Medal.BRONZE)).toEqual(new TrackmaniaTime(4000));
	});
});
