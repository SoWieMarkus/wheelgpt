import path from "node:path";
import cors from "cors";
import express, { json, type NextFunction, type Request, type Response } from "express";
import createHttpError, { isHttpError } from "http-errors";
import {
	AuthenticationRouter,
	ChannelRouter,
	LandingRouter,
	MetricsRouter,
	TrackmaniaRouter,
	TwitchWebhookRouter,
} from "./routes";
import { logger } from "./utils";

const app = express();

app.use(express.static(path.join(__dirname, "../../frontend/dist/frontend/browser")));
app.use(cors());
app.use(json({ limit: "7mb" }));

const apiRouter = express.Router();
apiRouter.use("/metrics", MetricsRouter);
apiRouter.use("/authentication", AuthenticationRouter);
apiRouter.use("/trackmania", TrackmaniaRouter);
apiRouter.use("/landing", LandingRouter);
apiRouter.use("/channel", ChannelRouter);
apiRouter.use("/twitch", TwitchWebhookRouter);

app.use("/api", apiRouter);
app.get("*name", (_, response) => {
	response.sendFile(path.join(__dirname, "../../frontend/dist/frontend/browser/index.html"));
});

// Handling of unknown endpoints
app.use((_, __, next) => {
	next(createHttpError(404, "Endpoint not found."));
});

// Error handling
app.use((error: unknown, _: Request, response: Response, __: NextFunction) => {
	const errorMessage = isHttpError(error) ? error.message : "An unknown error occured.";
	const errorStatus = isHttpError(error) ? error.status : 500;

	if (errorStatus >= 500) {
		logger.error(`Status ${errorStatus}: ${errorMessage}`);
		console.error(error);
	}
	response.status(errorStatus).json({ error: errorMessage });
});

export default app;
