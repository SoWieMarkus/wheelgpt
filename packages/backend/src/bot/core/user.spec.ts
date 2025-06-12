import type { ChatUserstate } from "tmi.js";
import { AccessLevel, getUser } from "./user";

describe("getUser", () => {
	const baseChatUser: ChatUserstate = {
		username: "testuser",
		"user-id": "123",
		"display-name": "TestUser",
		mod: false,
		vip: undefined,
		subscriber: false,
	};

	test("returns null if username is missing", () => {
		const user = getUser({ ...baseChatUser, username: undefined }, "#testuser");
		expect(user).toBeNull();
	});

	test("returns null if user-id is missing", () => {
		const user = getUser({ ...baseChatUser, "user-id": undefined }, "#testuser");
		expect(user).toBeNull();
	});

	test("returns null if display-name is missing", () => {
		const user = getUser({ ...baseChatUser, "display-name": undefined }, "#testuser");
		expect(user).toBeNull();
	});

	test("returns a user object with correct fields", () => {
		const user = getUser(baseChatUser, "#testuser");
		expect(user).toEqual({
			name: "testuser",
			id: "123",
			displayName: "TestUser",
			accessLevel: AccessLevel.STREAMER,
			channelId: "#testuser",
		});
	});

	test("assigns correct access levels", () => {
		expect(getUser({ ...baseChatUser, username: "mychannel" }, "#mychannel")?.accessLevel).toBe(AccessLevel.STREAMER);
		expect(getUser({ ...baseChatUser, mod: true }, "#otherchannel")?.accessLevel).toBe(AccessLevel.MOD);
		expect(getUser({ ...baseChatUser, vip: true }, "#otherchannel")?.accessLevel).toBe(AccessLevel.VIP);
		expect(getUser({ ...baseChatUser, subscriber: true }, "#otherchannel")?.accessLevel).toBe(AccessLevel.SUBSCRIBER);
		expect(getUser({ ...baseChatUser }, "#otherchannel")?.accessLevel).toBe(AccessLevel.USER);
	});
});
