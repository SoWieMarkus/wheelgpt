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

/** "2026-11" → "November 2026". */
export function monthLabel(key: string): string {
  const [year, month] = key.split("-").map(Number);
  if (!year || !month) return key;
  return new Date(Date.UTC(year, month - 1, 1)).toLocaleString("en", {
    month: "long",
    year: "numeric",
    timeZone: "UTC",
  });
}

export function periodLabel(key: string, current?: string): string {
  if (key === "all-time") return "All time";
  if (key === current) return "This month";
  return monthLabel(key);
}

const dateFormat = new Intl.DateTimeFormat("en", { dateStyle: "medium", timeStyle: "short" });
const dayFormat = new Intl.DateTimeFormat("en", { dateStyle: "medium" });

export function formatDate(iso: string, withTime = true): string {
  const d = new Date(iso);
  return withTime ? dateFormat.format(d) : dayFormat.format(d);
}

export function plural(n: number, one: string, many = `${one}s`): string {
  return `${n.toLocaleString("en")} ${n === 1 ? one : many}`;
}
