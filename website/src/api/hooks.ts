import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { ApiError, api, query } from "./client";
import type {
  BotConnection,
  Channel,
  Command,
  Entry,
  GameMap,
  List,
  Me,
  Page,
  Periods,
  Room,
  Settings,
  SettingsBody,
  Stats,
  TokenCreated,
  TokenInfo,
} from "./types";

export const keys = {
  me: ["me"] as const,
  token: ["me", "plugin-token"] as const,
  commands: ["commands"] as const,
  stats: ["stats"] as const,
  channels: ["channels"] as const,
  channel: (ref: string) => ["channel", ref] as const,
  map: (ref: string) => ["channel", ref, "map"] as const,
  room: (ref: string) => ["channel", ref, "room"] as const,
  periods: (ref: string) => ["channel", ref, "periods"] as const,
  board: (ref: string, period: string, q: string, offset: number) =>
    ["channel", ref, "board", period, q, offset] as const,
};

/** Resolves to null instead of throwing for the given status codes. */
async function orNull<T>(request: Promise<T>, ...statuses: number[]): Promise<T | null> {
  try {
    return await request;
  } catch (err) {
    if (err instanceof ApiError && statuses.includes(err.status)) return null;
    throw err;
  }
}

// --------------------------------------------------------------------- me --

/** The logged in streamer, or null when there is no session. */
export function useMe() {
  return useQuery({
    queryKey: keys.me,
    queryFn: () => orNull(api<Me>("/me"), 401),
    staleTime: 60_000,
    retry: false,
  });
}

export function useLogout() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api<undefined>("/auth/logout", { method: "POST" }),
    onSettled: () => {
      qc.clear();
      qc.setQueryData(keys.me, null);
    },
  });
}

export function useDeleteChannel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api<undefined>("/me", { method: "DELETE" }),
    onSuccess: () => {
      qc.clear();
      qc.setQueryData(keys.me, null);
    },
  });
}

export function useUpdateSettings() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: SettingsBody) => api<Settings>("/me/settings", { method: "PATCH", body }),
    onSuccess: (settings) =>
      qc.setQueryData<Me | null>(keys.me, (me) => (me ? { ...me, settings } : me)),
  });
}

export function useBotConnection() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (connect: boolean) =>
      api<BotConnection>("/me/bot-connection", { method: connect ? "PUT" : "DELETE" }),
    onSuccess: (bot) => qc.setQueryData<Me | null>(keys.me, (me) => (me ? { ...me, bot } : me)),
  });
}

/** Dev mode only: pretend the stream started or ended. */
export function useLiveSwitch() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (live: boolean) => api<Channel>("/me/live", { method: live ? "PUT" : "DELETE" }),
    onSuccess: (channel) => {
      qc.setQueryData<Me | null>(keys.me, (me) => (me ? { ...me, channel } : me));
      // Ending the stream clears map and room.
      void qc.invalidateQueries({ queryKey: keys.map(channel.id) });
      void qc.invalidateQueries({ queryKey: keys.room(channel.id) });
    },
  });
}

export function usePluginToken() {
  return useQuery({ queryKey: keys.token, queryFn: () => api<TokenInfo>("/me/plugin-token") });
}

export function useCreatePluginToken() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api<TokenCreated>("/me/plugin-token", { method: "POST" }),
    onSuccess: () => qc.invalidateQueries({ queryKey: keys.token }),
  });
}

// --------------------------------------------------------------- public --

export function useCommands() {
  return useQuery({
    queryKey: keys.commands,
    queryFn: () => api<List<Command>>("/commands").then((r) => r.items),
    staleTime: Infinity,
  });
}

export function useStats() {
  return useQuery({ queryKey: keys.stats, queryFn: () => api<Stats>("/stats") });
}

export function useChannels() {
  return useQuery({
    queryKey: keys.channels,
    queryFn: () => api<Page<Channel>>(`/channels${query({ limit: 100 })}`),
  });
}

export function useChannel(ref: string) {
  return useQuery({
    queryKey: keys.channel(ref),
    queryFn: () => api<Channel>(`/channels/${encodeURIComponent(ref)}`),
    retry: false,
  });
}

/** The map the streamer is on, or null. */
export function useMap(ref: string) {
  return useQuery({
    queryKey: keys.map(ref),
    queryFn: () => orNull(api<GameMap>(`/channels/${encodeURIComponent(ref)}/map`), 404),
  });
}

/** The room the streamer is in, or null. */
export function useRoom(ref: string) {
  return useQuery({
    queryKey: keys.room(ref),
    queryFn: () => orNull(api<Room>(`/channels/${encodeURIComponent(ref)}/room`), 404),
  });
}

export function usePeriods(ref: string) {
  return useQuery({
    queryKey: keys.periods(ref),
    queryFn: () => api<Periods>(`/channels/${encodeURIComponent(ref)}/leaderboards`),
  });
}

export const BOARD_PAGE_SIZE = 50;

export function useBoard(ref: string, period: string | undefined, q: string, offset: number) {
  return useQuery({
    queryKey: keys.board(ref, period ?? "", q, offset),
    queryFn: () =>
      api<Page<Entry>>(
        `/channels/${encodeURIComponent(ref)}/leaderboards/${period}${query({ q, limit: BOARD_PAGE_SIZE, offset: offset || undefined })}`,
      ),
    enabled: period !== undefined,
    placeholderData: (previous) => previous,
  });
}
