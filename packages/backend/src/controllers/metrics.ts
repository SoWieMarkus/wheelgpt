import type { RequestHandler } from "express";
import { prometheus } from "../prometheus";

export const getMetrics: RequestHandler = async (request, response) => {
	response.set("Content-Type", prometheus.contentType);
	const metrics = await prometheus.metrics();
	response.end(metrics);
};
