/** Formats a race time like the game: 12.345, 1:02.345 or 1:01:02.345. */
export function formatTime(ms: number): string {
  const v = Math.abs(Math.trunc(ms));
  const millis = String(v % 1000).padStart(3, "0");
  const s = Math.floor(v / 1000) % 60;
  const m = Math.floor(v / 60_000) % 60;
  const h = Math.floor(v / 3_600_000);
  if (h > 0) return `${h}:${pad(m)}:${pad(s)}.${millis}`;
  if (m > 0) return `${m}:${pad(s)}.${millis}`;
  return `${s}.${millis}`;
}

/** Formats a signed difference: +0.120, -1.005, ±0.000. */
export function formatDiff(ms: number): string {
  if (ms === 0) return "±0.000";
  return `${ms > 0 ? "+" : "-"}${formatTime(ms)}`;
}

function pad(n: number): string {
  return String(n).padStart(2, "0");
}

/** "2026-11" → "November 2026" (or the given locale's equivalent). */
export function monthLabel(key: string, locale = "en"): string {
  const [year, month] = key.split("-").map(Number);
  if (!year || !month) return key;
  return new Date(Date.UTC(year, month - 1, 1)).toLocaleString(locale, {
    month: "long",
    year: "numeric",
    timeZone: "UTC",
  });
}

export function formatDate(iso: string, locale = "en", withTime = true): string {
  const options: Intl.DateTimeFormatOptions = withTime
    ? { dateStyle: "medium", timeStyle: "short" }
    : { dateStyle: "medium" };
  return new Intl.DateTimeFormat(locale, options).format(new Date(iso));
}

export function formatNumber(n: number, locale = "en"): string {
  return n.toLocaleString(locale);
}
