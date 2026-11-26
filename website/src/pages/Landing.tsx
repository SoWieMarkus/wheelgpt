import type { ReactNode } from "react";
import { Link } from "react-router";
import { loginUrl } from "../api/client";
import { useChannels, useCommands, useMe, useStats } from "../api/hooks";
import { Avatar, Emote, ErrorNotice, LiveBadge, Loading } from "../components/ui";
import { useI18n } from "../i18n";

export function Landing() {
  const { data: session } = useMe();
  const { t } = useI18n();
  return (
    <>
      <section className="intro wrap">
        <div className="intro__text">
          <p className="eyebrow">{t.landing.eyebrow}</p>
          <h1 className="visually-hidden">WheelGPT</h1>
          <p className="lede">
            {t.landing.lede}
            <Emote name="HeyGuys" file="hey_guys" />
          </p>
          <p className="intro__more">
            {t.landing.intro(
              <a href="https://openplanet.dev/plugin/wheelgpt-plugin">{t.landing.pluginName}</a>,
            )}
          </p>
          <div className="intro__actions">
            {session ? (
              <Link to="/dashboard" className="button">
                {t.landing.openDashboard}
              </Link>
            ) : (
              <a href={loginUrl()} className="button button--twitch">
                {t.layout.login}
              </a>
            )}
            <a
              href="https://openplanet.dev/plugin/wheelgpt-plugin"
              className="button button--secondary"
            >
              {t.landing.getPlugin}
            </a>
          </div>
        </div>
        <ChatSample />
      </section>

      <section className="wrap section" aria-labelledby="commands">
        <h2 id="commands">{t.landing.commandsTitle}</h2>
        <p className="section__lede">{t.landing.commandsLede(<code>h:mm:ss.xxx</code>)}</p>
        <Commands />
      </section>

      <section className="wrap section" id="streamers" aria-labelledby="streamers-title">
        <h2 id="streamers-title">{t.landing.streamersTitle}</h2>
        <Streamers />
      </section>
    </>
  );
}

type Badge = "mod" | "vip" | "sub" | "streamer";

/** One chat line with Twitch badges, like v1 showed them. */
function Line({
  user,
  color,
  badges = [],
  bot = false,
  children,
}: {
  user: string;
  color?: string;
  badges?: Badge[];
  bot?: boolean;
  children: ReactNode;
}) {
  const badgeNames = useI18n().t.landing.badges;
  return (
    <li className={bot ? "chat__bot" : undefined}>
      {badges.map((b) => (
        <img
          key={b}
          className="chat__badge"
          src={`/badges/${b}.png`}
          alt={badgeNames[b]}
          title={badgeNames[b]}
          width={18}
          height={18}
        />
      ))}
      <b style={color ? { color } : undefined}>{user}</b> {children}
    </li>
  );
}

function ChatSample() {
  const { t } = useI18n();
  const c = t.landing.chat;
  return (
    <figure className="chat" aria-label={t.landing.chatLabel}>
      <ol className="chat__lines">
        <Line user="kazoo_k" color="#c0392b" badges={["sub"]}>
          !guess 47.912
        </Line>
        <Line user="marbled" color="#2471a3">
          !g 48.3 <span className="chat__hint">{c.needsMs}</span>
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          @marbled {c.format}
        </Line>
        <Line user="tinytires" color="#1e8449" badges={["vip"]}>
          !guess 47.850 <Emote name="YEK" file="yek" />
        </Line>
        <Line user="mod_otter" color="#7d3c98" badges={["mod", "sub"]}>
          !room
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          @mod_otter Your Mamas Room [5/100]
        </Line>
        <Line user="GranaDyy" color="#d35400" badges={["streamer"]}>
          !map
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          @GranaDyy "Spring 2026 - 07" by Nadeo | <Emote name="AuthorTime" file="author_medal" />{" "}
          47.904 | trackmania.io/…
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          {c.pb} <Emote name="dinkDonk" file="dink_donk" /> {c.thats}{" "}
          <Emote name="AuthorTime" file="author_medal" /> <Emote name="ok" file="ok" />
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          {c.nobody} <Emote name="ok" file="ok" /> {c.point}
        </Line>
      </ol>
    </figure>
  );
}

function Commands() {
  const { data, error, isPending, refetch } = useCommands();
  const { t } = useI18n();
  if (isPending) return <Loading />;
  if (error) return <ErrorNotice error={error} retry={() => void refetch()} />;
  return (
    <div className="table-scroll">
      <table className="table commands">
        <thead>
          <tr>
            <th scope="col">{t.landing.table.command}</th>
            <th scope="col">{t.landing.table.what}</th>
            <th scope="col">{t.landing.table.who}</th>
          </tr>
        </thead>
        <tbody>
          {data.map((c) => (
            <tr key={c.name}>
              <td>
                <code>{c.usage}</code>
                <OtherNames command={c} />
              </td>
              <td>{t.landing.commands[c.name] ?? c.description}</td>
              <td className="muted">{t.landing.access[c.access] ?? c.access}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

/** Names of a command that its usage example doesn't show already. */
function OtherNames({ command }: { command: { name: string; aliases: string[]; usage: string } }) {
  const shown = command.usage.split(/[\s,]+/)[0]?.slice(1);
  const { t } = useI18n();
  const others = [command.name, ...command.aliases].filter((n) => n !== shown);
  if (others.length === 0) return null;
  return (
    <span className="muted small">
      {" "}
      {t.landing.table.also} {others.map((n) => `!${n}`).join(", ")}
    </span>
  );
}

function Streamers() {
  const channels = useChannels();
  const stats = useStats();
  const { t } = useI18n();

  if (channels.isPending) return <Loading />;
  if (channels.error) {
    return <ErrorNotice error={channels.error} retry={() => void channels.refetch()} />;
  }
  const list = channels.data.items;
  return (
    <>
      {stats.data && <p className="section__lede">{t.landing.channelCount(stats.data.channels)}</p>}
      {list.length === 0 ? (
        <p className="muted">{t.landing.noStreamers}</p>
      ) : (
        <ul className="streamers">
          {list.map((c) => (
            <li key={c.id}>
              <Link to={`/c/${c.login}`} className="streamer">
                <Avatar src={c.profileImage} name={c.displayName} size={44} />
                <span className="streamer__name">{c.displayName}</span>
                <LiveBadge live={c.isLive} />
              </Link>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}
