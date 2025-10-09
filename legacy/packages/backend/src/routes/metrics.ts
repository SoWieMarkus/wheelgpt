import express from "express";
import { MetricsController } from "../controllers";

const router = express.Router();

router.get("/", MetricsController.getMetrics);

export default router;
