import express from "express";
import { GameClientController } from "../controllers";
import { requiresChannelAuthentication } from "../middleware/channel-authentication";

const router = express.Router();

router.post("/update/map", requiresChannelAuthentication, GameClientController.updateMap);
router.post("/update/pb", requiresChannelAuthentication, GameClientController.updatePersonalBest);

export default router;