import { useState } from "react";
import { Link, useParams, useSearchParams } from "react-router";
import { ApiError } from "../api/client";
import { BOARD_PAGE_SIZE, useBoard, useChannel, useMap, usePeriods, useRoom } from "../api/hooks";
import { Avatar, ErrorNotice, LiveBadge, Loading, Medals, Notice } from "../components/ui";
import { useI18n } from "../i18n";
import { formatNumber, monthLabel } from "../lib/format";

const PER_PAGE = BOARD_PAGE_SIZE;

export function Channel() {
  const { channel = "" } = useParams();
  const [params, setParams] = useSearchParams();
  const [search, setSearch] = useState(params.get("q") ?? "");
  const { t, lang } = useI18n();

  const info = useChannel(channel);
  const periods = usePeriods(channel);
  const map = useMap(channel);
  const room = useRoom(channel);

  const current = periods.data?.current;
  const periodLabel = (key: string) =>
    key === "all-time"
      ? t.periods.allTime
      : key === current
        ? t.periods.thisMonth
        : monthLabel(key, lang);
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
        <h1>{t.channel.notFoundTitle}</h1>
        <p>
          {t.channel.notFound(<strong>{channel}</strong>, (text) => (
            <Link to="/#streamers">{text}</Link>
          ))}
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
              <p className="eyebrow">{t.channel.eyebrow}</p>
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
        <aside className="now-driving" aria-label={t.channel.currentMap}>
          <p className="small muted">{t.channel.driving}</p>
          <p className="map-name">{map.data.name}</p>
          <p className="muted small">{t.channel.by(map.data.author)}</p>
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
                · {t.channel.room} <strong>{room.data.name}</strong> ({room.data.numberOfPlayers}/
                {room.data.maxPlayers})
              </>
            )}
          </p>
        </aside>
      )}

      <div className="board-controls">
        <div className="tabs" role="tablist" aria-label={t.channel.period}>
          {tabs.map((p) => (
            <button
              key={p}
              type="button"
              role="tab"
              aria-selected={p === period}
              className="tab"
              onClick={() => set({ period: p === current ? undefined : p, page: undefined })}
            >
              {periodLabel(p)}
            </button>
          ))}
          {older.length > 0 && (
            <select
              aria-label={t.channel.olderMonths}
              className="tab-select"
              value={period && older.includes(period) ? period : ""}
              onChange={(e) => set({ period: e.currentTarget.value || undefined, page: undefined })}
            >
              <option value="">{t.channel.older}</option>
              {older.map((p) => (
                <option key={p} value={p}>
                  {periodLabel(p)}
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
            {t.channel.findViewer}
          </label>
          <input
            id="board-q"
            type="search"
            placeholder={t.channel.findViewer}
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
            ? t.channel.nobodyCalled(q)
            : period === current
              ? t.channel.noPointsMonth
              : t.channel.noPoints}
        </Notice>
      ) : (
        <>
          <p className="muted small">{t.channel.summary(total)}</p>
          <div className="table-scroll">
            <table className="table board">
              <thead>
                <tr>
                  <th scope="col" className="num">
                    #
                  </th>
                  <th scope="col">{t.channel.viewer}</th>
                  <th scope="col" className="num">
                    {t.channel.points}
                  </th>
                  <th scope="col" className="num">
                    {t.channel.perfect}
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
                    <td className="num">{formatNumber(e.points, lang)}</td>
                    <td className="num muted">{e.perfectGuesses || "–"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {pages > 1 && (
            <nav className="pager" aria-label={t.channel.pages}>
              <button
                type="button"
                className="button button--secondary button--small"
                disabled={page <= 1}
                onClick={() => set({ page: String(page - 1) })}
              >
                {t.channel.previous}
              </button>
              <span className="muted small">{t.channel.pageOf(page, pages)}</span>
              <button
                type="button"
                className="button button--secondary button--small"
                disabled={page >= pages}
                onClick={() => set({ page: String(page + 1) })}
              >
                {t.channel.next}
              </button>
            </nav>
          )}
        </>
      )}
    </div>
  );
}
