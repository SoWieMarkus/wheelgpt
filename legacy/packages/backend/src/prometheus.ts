import { collectDefaultMetrics, Registry } from "prom-client";

export const prometheus = new Registry();

// Collect default Node.js metrics (CPU, memory, etc.)
collectDefaultMetrics({ register: prometheus, prefix: "wheelgpt_" });
