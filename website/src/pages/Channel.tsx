import { useState } from "react";
import { Link, useParams, useSearchParams } from "react-router";
import { ApiError } from "../api/client";
import { BOARD_PAGE_SIZE, useBoard, useChannel, useMap, usePeriods, useRoom } from "../api/hooks";
import { Avatar, ErrorNotice, LiveBadge, Loading, Medals, Notice } from "../components/ui";
import { periodLabel, plural } from "../lib/format";

const PER_PAGE = BOARD_PAGE_SIZE;

export function Channel() {
  const { channel = "" } = useParams();
  const [params, setParams] = useSearchParams();
  const [search, setSearch] = useState(params.get("q") ?? "");

  const info = useChannel(channel);
  const periods = usePeriods(channel);
  const map = useMap(channel);
  const room = useRoom(channel);

  const current = periods.data?.current;
  const period = params.get("period") ?? current;
  const page = Math.max(1, Number(params.get("page")) || 1);
  const q = params.get("q") ?? "";
  const board = useBoard(channel, period, q, (page - 1) * PER_PAGE);

  const set = (next: Record<string, string | undefined>) => {
    const p = new URLSearchParams(params);
    for (const [k, v] of Object.entries(next)) {
      if (v === undefined || v === "") p.delete(k);
      else p.set(k, v);
    }
    setParams(p, { replace: true });
  };

  if (info.error instanceof ApiError && info.error.status === 404) {
    return (
      <div className="wrap section narrow">
        <h1>Channel not found</h1>
        <p>
          <strong>{channel}</strong> doesn't use WheelGPT.{" "}
          <Link to="/#streamers">See who does.</Link>
        </p>
      </div>
    );
  }

  // This month and last month as tabs, all time last, older months in a menu.
  const months = (periods.data?.periods ?? []).filter((p) => p !== "all-time");
  const tabs = [...months.slice(0, 2), "all-time"];
  const older = months.slice(2);

  const ch = info.data;
  const total = board.data?.total ?? 0;
  const pages = Math.max(1, Math.ceil(total / PER_PAGE));

  return (
    <div className="wrap channel">
      <header className="channel-head">
        {ch ? (
          <>
            <Avatar src={ch.profileImage} name={ch.displayName} size={72} />
            <div>
              <p className="eyebrow">Guessing leaderboard</p>
              <h1>{ch.displayName}</h1>
              <p className="muted">
                <LiveBadge live={ch.isLive} /> ·{" "}
                <a href={`https://twitch.tv/${ch.login}`}>twitch.tv/{ch.login}</a>
              </p>
            </div>
          </>
        ) : (
          <Loading />
        )}
      </header>

      {map.data && (
        <aside className="now-driving" aria-label="Current map">
          <p className="small muted">Driving right now</p>
          <p className="map-name">{map.data.name}</p>
          <p className="muted small">by {map.data.author}</p>
          <Medals medals={map.data.medals} />
          <p className="small">
            {map.data.tmxUrl && (
              <>
                <a href={map.data.tmxUrl}>Trackmania Exchange</a> ·{" "}
              </>
            )}
            <a href={map.data.trackmaniaIoUrl}>trackmania.io</a>
            {room.data && (
              <>
                {" "}
                · Room <strong>{room.data.name}</strong> ({room.data.numberOfPlayers}/
                {room.data.maxPlayers})
              </>
            )}
          </p>
        </aside>
      )}

      <div className="board-controls">
        <div className="tabs" role="tablist" aria-label="Period">
          {tabs.map((p) => (
            <button
              key={p}
              type="button"
              role="tab"
              aria-selected={p === period}
              className="tab"
              onClick={() => set({ period: p === current ? undefined : p, page: undefined })}
            >
              {periodLabel(p, current)}
            </button>
          ))}
          {older.length > 0 && (
            <select
              aria-label="Older months"
              className="tab-select"
              value={period && older.includes(period) ? period : ""}
              onChange={(e) => set({ period: e.currentTarget.value || undefined, page: undefined })}
            >
              <option value="">Older…</option>
              {older.map((p) => (
                <option key={p} value={p}>
                  {periodLabel(p, current)}
                </option>
              ))}
            </select>
          )}
        </div>
        <form
          className="board-search"
          onSubmit={(e) => {
            e.preventDefault();
            set({ q: search.trim(), page: undefined });
          }}
        >
          <label className="visually-hidden" htmlFor="board-q">
            Find a viewer
          </label>
          <input
            id="board-q"
            type="search"
            placeholder="Find a viewer"
            value={search}
            maxLength={50}
            onChange={(e) => setSearch(e.currentTarget.value)}
          />
        </form>
      </div>

      {!board.data ? (
        board.error ? (
          <ErrorNotice error={board.error} retry={() => void board.refetch()} />
        ) : (
          <Loading />
        )
      ) : board.data.items.length === 0 ? (
        <Notice>
          {q
            ? `Nobody called “${q}” has points in this period.`
            : period === current
              ? "No points this month yet. The first personal best of the month starts the race."
              : "No points in this period."}
        </Notice>
      ) : (
        <>
          <p className="muted small">
            {plural(total, "guesser")} with points · perfect guess 5 points, closest guess 1 point
          </p>
          <div className="table-scroll">
            <table className="table board">
              <thead>
                <tr>
                  <th scope="col" className="num">
                    #
                  </th>
                  <th scope="col">Viewer</th>
                  <th scope="col" className="num">
                    Points
                  </th>
                  <th scope="col" className="num">
                    Perfect
                  </th>
                </tr>
              </thead>
              <tbody>
                {board.data.items.map((e) => (
                  <tr
                    key={e.userId}
                    className={e.position <= 3 ? `board__top board__top--${e.position}` : undefined}
                  >
                    <td className="num rank">{e.position}</td>
                    <td>{e.displayName}</td>
                    <td className="num">{e.points.toLocaleString("en")}</td>
                    <td className="num muted">{e.perfectGuesses || "–"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {pages > 1 && (
            <nav className="pager" aria-label="Pages">
              <button
                type="button"
                className="button button--secondary button--small"
                disabled={page <= 1}
                onClick={() => set({ page: String(page - 1) })}
              >
                Previous
              </button>
              <span className="muted small">
                Page {page} of {pages}
              </span>
              <button
                type="button"
                className="button button--secondary button--small"
                disabled={page >= pages}
                onClick={() => set({ page: String(page + 1) })}
              >
                Next
              </button>
            </nav>
          )}
        </>
      )}
    </div>
  );
}
