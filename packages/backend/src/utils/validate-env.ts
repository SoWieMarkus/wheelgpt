import "dotenv/config";
import dotenv from "dotenv";

dotenv.config({ path: "../../.env" });

import { cleanEnv, port, str, url } from "envalid";

export default cleanEnv(process.env, {
	PORT: port(),
	DATABASE_URL: str(),
	JWT_SECRET_CHANNEL: str(),
	JWT_SECRET_WEB: str(),
	TWITCH_CLIENT_ID: str(),
	TWITCH_CLIENT_SECRET: str(),
	TWITCH_REDIRECT_URL: str(),
	BOT_USERNAME: str(),
	BOT_OAUTH_TOKEN: str(),
	TWITCH_EVENTSUB_SECRET: str(),
	TWITCH_STREAM_WEB_HOOK_URL: url(),
});
