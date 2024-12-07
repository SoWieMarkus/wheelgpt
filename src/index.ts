import { env, Log } from "./utils";
import app from "./app";
import { WheelGPT } from "./bot/bot";

Log.info("Booting...");
app.listen(env.PORT, () => Log.complete(`Running on port "${env.PORT}"`));
WheelGPT.connect().catch(console.error);

