import express from "express";
import { EmoteController } from "../controllers";

const router = express.Router();

router.get("/", EmoteController.getEmotes);

export default router;
