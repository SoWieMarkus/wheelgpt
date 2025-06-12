import express from "express";
import { TrackmaniaController } from "../controllers";
import { requiresPluginAuthentication } from "../middlewares/plugin-authentication";

const router = express.Router();

router.post("/update/map", requiresPluginAuthentication, TrackmaniaController.updateMap);
router.post("/update/pb", requiresPluginAuthentication, TrackmaniaController.updatePersonalBest);

export default router;
