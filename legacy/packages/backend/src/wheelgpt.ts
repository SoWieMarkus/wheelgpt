import { WheelGPT } from "./bot";
import { env } from "./utils";

const username = env.BOT_USERNAME;
const password = env.BOT_OAUTH_TOKEN;

export const wheelgpt = new WheelGPT(username, password);
