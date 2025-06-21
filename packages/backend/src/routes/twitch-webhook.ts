import express from "express";
import { TwitchWebhookController } from "../controllers";
import { verifyTwitchSignature } from "../middlewares/twitch-verification";

const router = express.Router();

router.post("/webhook/stream", verifyTwitchSignature, TwitchWebhookController.streamStateWebhook);

export default router;
