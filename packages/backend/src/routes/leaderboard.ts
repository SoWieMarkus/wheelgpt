import express from "express";
import { LeaderboardController } from "../controllers";

const router = express.Router();

router.get("/:channelId", LeaderboardController.getLeaderboardPage);
router.get("/:channelId/:displayName", LeaderboardController.getLeaderboardPositionByName);

export default router;
