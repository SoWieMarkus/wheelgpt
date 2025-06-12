import app from "./app";
import { env, logger } from "./utils";
import { wheelgpt } from "./wheelgpt";

// TODO start bot
// check what channels are online

// register webhooks

logger.info("Starting the backend server...");
app.listen(env.PORT, () => logger.info(`Running on port "${env.PORT}"`));
