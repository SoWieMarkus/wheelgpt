import express from "express";
import { LandingController } from "../controllers";

const router = express.Router();

router.get("/", LandingController.getPublicChannels);

export default router;
