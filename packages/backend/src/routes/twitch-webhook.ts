import express from "express";
import { TwitchWebhookController } from "../controllers";
import { verifyTwitchSignature } from "../middlewares/twitch-verification";

const router = express.Router();

router.post(
	"/webhook/stream",
	// the raw body parser is needed for signature verification
	express.raw({
		type: "application/json",
	}),
	verifyTwitchSignature,
	express.json(),
	TwitchWebhookController.streamStateWebhook,
);

export default router;
