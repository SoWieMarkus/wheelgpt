import { getCommandArguments } from "./arguments";

describe("getCommandArguments", () => {
	it("returns null for empty message", () => {
		expect(getCommandArguments("")).toBeNull();
	});

	it("returns null if message does not start with command symbol", () => {
		expect(getCommandArguments("hello world")).toBeNull();
		expect(getCommandArguments(" test")).toBeNull();
	});

	it("parses command with no arguments", () => {
		expect(getCommandArguments("!foo")).toEqual({ key: "foo", args: [] });
	});

	it("parses command with arguments", () => {
		expect(getCommandArguments("!bar arg1 arg2")).toEqual({
			key: "bar",
			args: ["arg1", "arg2"],
		});
	});

	it("parses command with mixed case and extra spaces", () => {
		expect(getCommandArguments("!TeSt   a b  c")).toEqual({
			key: "test",
			args: ["a", "b", "c"],
		});
	});

	it("handles command symbol only", () => {
		expect(getCommandArguments("!")).toBeNull();
	});
});
