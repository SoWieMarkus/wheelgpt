import express from "express";
import { ChannelController } from "../controllers";

const router = express.Router();

router.get("/me", ChannelController.me);
router.post("/settings", ChannelController.updateSettings);

export default router;
