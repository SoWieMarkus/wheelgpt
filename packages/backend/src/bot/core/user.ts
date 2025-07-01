import type { ChatUserstate } from "tmi.js";

export const AccessLevel = {
	STREAMER: 4,
	MOD: 3,
	VIP: 2,
	SUBSCRIBER: 1,
	USER: 0,
};

const getAccessLevel = (chatUser: ChatUserstate, channelId: string) => {
	if (`#${chatUser.username}` === channelId) return AccessLevel.STREAMER;
	if (chatUser.mod) return AccessLevel.MOD;
	if (chatUser.vip !== undefined) return AccessLevel.VIP;
	if (chatUser.subscriber) return AccessLevel.SUBSCRIBER;
	return AccessLevel.USER;
};

export type User = {
	name: string;
	id: string;
	displayName: string;
	accessLevel: number;
	channelId: string;
};

export const getUser = (chatUser: ChatUserstate, channelId: string): User | null => {
	const name = chatUser.username;
	const id = chatUser["user-id"];
	const displayName = chatUser["display-name"];

	if (name === undefined || id === undefined || displayName === undefined) {
		return null;
	}

	const accessLevel = getAccessLevel(chatUser, channelId);
	return { name, id, displayName, accessLevel, channelId };
};

export const mentionUser = (username: string): string => {
	if (!username) return "";
	const usernameWithoutAt = username.replace(/^@/, "");
	return `@${usernameWithoutAt} `;
}