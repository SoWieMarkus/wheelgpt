import type { Problem } from "./types";

/** Error thrown for every non-2xx response, carrying the RFC 9457 problem. */
export class ApiError extends Error {
  readonly status: number;
  readonly problem: Problem | null;

  constructor(status: number, problem: Problem | null) {
    super(problem?.detail ?? `Request failed with status ${status}`);
    this.name = "ApiError";
    this.status = status;
    this.problem = problem;
  }

  get code(): string | undefined {
    return this.problem?.code;
  }

  /** Field level validation messages keyed by JSON field name. */
  get fieldErrors(): Record<string, string> {
    const out: Record<string, string> = {};
    for (const e of this.problem?.errors ?? []) out[e.field] = e.message;
    return out;
  }
}

type Method = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

interface RequestOptions {
  method?: Method;
  body?: unknown;
  signal?: AbortSignal;
}

/**
 * Calls the API. The session lives in an HttpOnly cookie, so there is no
 * token handling here at all: the browser sends the cookie same-origin.
 */
export async function api<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = { Accept: "application/json" };
  let body: BodyInit | undefined;
  if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
    body = JSON.stringify(options.body);
  }

  const response = await fetch(`/api/v1${path}`, {
    method: options.method ?? "GET",
    headers,
    body,
    credentials: "same-origin",
    signal: options.signal,
  });

  if (!response.ok) {
    const isProblem = response.headers.get("Content-Type")?.includes("problem+json") ?? false;
    const problem = isProblem ? ((await response.json()) as Problem) : null;
    throw new ApiError(response.status, problem);
  }
  if (response.status === 204) return undefined as T;
  return (await response.json()) as T;
}

/** Builds a query string, skipping empty values. */
export function query(params: Record<string, string | number | boolean | undefined>): string {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === "" || value === false) continue;
    search.set(key, String(value));
  }
  const s = search.toString();
  return s ? `?${s}` : "";
}

/** Where the Twitch login starts (a browser redirect, not a fetch). */
export function loginUrl(returnTo = "/dashboard"): string {
  return `/api/v1/auth/twitch?returnTo=${encodeURIComponent(returnTo)}`;
}
