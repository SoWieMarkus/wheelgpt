import type { RequestHandler } from "express";
import createHttpError from "http-errors";
import { wheelgpt } from "../wheelgpt";

export const requiresWheelGPTInitializedAuthentication: RequestHandler = (_, __, next) => {
	if (!wheelgpt.initialized) {
		const httpError = createHttpError(400, "WheelGPT is currently booting. Please try again later.");
		return next(httpError);
	}
	next();
};
