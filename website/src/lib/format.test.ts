import { describe, expect, it } from "vitest";
import { formatDiff, formatNumber, formatTime, monthLabel } from "./format";

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

describe("monthLabel", () => {
  it("labels months in the given locale", () => {
    expect(monthLabel("2026-11")).toBe("November 2026");
    expect(monthLabel("2026-03", "de")).toBe("März 2026");
    expect(monthLabel("all-time")).toBe("all-time");
  });
});

describe("formatNumber", () => {
  it("groups digits per locale", () => {
    expect(formatNumber(1200)).toBe("1,200");
    expect(formatNumber(1200, "de")).toBe("1.200");
  });
});
