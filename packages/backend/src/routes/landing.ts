import express from "express";
import { LandingController } from "../controllers";

const router = express.Router();

router.get("/channels", LandingController.getPublicChannels);

export default router;
