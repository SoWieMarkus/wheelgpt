import { describe, expect, it } from "vitest";
import { formatDiff, formatTime, monthLabel, periodLabel, plural } from "./format";

describe("formatTime", () => {
  it.each([
    [0, "0.000"],
    [1, "0.001"],
    [12_345, "12.345"],
    [62_345, "1:02.345"],
    [3_662_345, "1:01:02.345"],
  ])("%i → %s", (ms, want) => {
    expect(formatTime(ms)).toBe(want);
  });
});

describe("formatDiff", () => {
  it("signs differences", () => {
    expect(formatDiff(120)).toBe("+0.120");
    expect(formatDiff(-1_005)).toBe("-1.005");
    expect(formatDiff(0)).toBe("±0.000");
  });
});

describe("periods", () => {
  it("labels months", () => {
    expect(monthLabel("2026-11")).toBe("November 2026");
    expect(periodLabel("all-time")).toBe("All time");
    expect(periodLabel("2026-11", "2026-11")).toBe("This month");
    expect(periodLabel("2026-10", "2026-11")).toBe("October 2026");
  });
});

describe("plural", () => {
  it("picks the right form", () => {
    expect(plural(1, "point")).toBe("1 point");
    expect(plural(1200, "point")).toBe("1,200 points");
  });
});
