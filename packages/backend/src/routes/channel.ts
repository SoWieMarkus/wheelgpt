import express from "express";
import { ChannelController } from "../controllers";
import { requiresWebAuthentication } from "../middlewares/web-authentication";

const router = express.Router();

router.get("/me", requiresWebAuthentication, ChannelController.me);
router.post("/settings", requiresWebAuthentication, ChannelController.updateSettings);

export default router;
