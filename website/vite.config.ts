import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

// In development Vite proxies the API, so the site and the API share one
// origin (session cookie, CSRF checks), like nginx does in production.
const backend = process.env.VITE_BACKEND_URL ?? "http://localhost:8080";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api": { target: backend, changeOrigin: false },
    },
  },
  build: {
    sourcemap: true,
  },
  test: {
    environment: "jsdom",
    include: ["src/**/*.test.{ts,tsx}"],
  },
});
