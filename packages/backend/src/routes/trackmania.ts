import express from "express";
import { TrackmaniaController } from "../controllers";
import { requiresPluginAuthentication } from "../middlewares/plugin-authentication";

const router = express.Router();

router.post("/map", requiresPluginAuthentication, TrackmaniaController.updateMap);
router.post("/pb", requiresPluginAuthentication, TrackmaniaController.updatePersonalBest);
router.post("/room", requiresPluginAuthentication, TrackmaniaController.updateRoom);
router.delete("/map", requiresPluginAuthentication, TrackmaniaController.deleteMap);
router.delete("/room", requiresPluginAuthentication, TrackmaniaController.deleteRoom);

export default router;
