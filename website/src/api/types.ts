// Friendly names for the generated API schema. Run `make generate` after
// changing the Go handlers to refresh schema.gen.ts.
import type { components } from "./schema.gen";

type Schemas = components["schemas"];

export type Problem = Schemas["httpx.Problem"];
export type Channel = Schemas["channel.SummaryResponse"];
export type Me = Schemas["channel.MeResponse"];
export type Settings = Schemas["channel.SettingsResponse"];
export type SettingsBody = Schemas["channel.SettingsBody"];
export type BotConnection = Schemas["channel.BotResponse"];
export type Stats = Schemas["channel.StatsResponse"];
export type Command = Schemas["bot.CommandResponse"];
export type GameMap = Schemas["ingame.MapResponse"];
export type Medals = Schemas["ingame.Medals"];
export type Room = Schemas["ingame.RoomBody"];
export type Periods = Schemas["leaderboard.PeriodsResponse"];
export type Entry = Schemas["leaderboard.EntryResponse"];
export type TokenInfo = Schemas["plugin.TokenInfoResponse"];
export type TokenCreated = Schemas["plugin.TokenCreatedResponse"];

export interface List<T> {
  items: T[];
}

export interface Page<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}
