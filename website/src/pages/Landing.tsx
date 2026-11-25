import type { ReactNode } from "react";
import { Link } from "react-router";
import { loginUrl } from "../api/client";
import { useChannels, useCommands, useMe, useStats } from "../api/hooks";
import { Avatar, Emote, ErrorNotice, LiveBadge, Loading } from "../components/ui";
import { plural } from "../lib/format";

export function Landing() {
  const { data: session } = useMe();
  return (
    <>
      <section className="intro wrap">
        <div className="intro__text">
          <p className="eyebrow">A Twitch bot for Trackmania streams</p>
          <h1 className="visually-hidden">WheelGPT</h1>
          <p className="lede">
            Jaaa hello guys! My name is WheelGPT! And especially mothers love me!{" "}
            <Emote name="HeyGuys" file="hey_guys" />
          </p>
          <p className="intro__more">
            I am a Twitch bot made especially for Trackmania streamers. With me, your viewers can
            always check which map you are currently playing. They can also guess your next personal
            best time. To make everything run smoothly, you should install the{" "}
            <a href="https://openplanet.dev/plugin/wheelgpt-plugin">OpenPlanet Plugin</a>. Curious?
            Then add WheelGPT to your channel now!
          </p>
          <div className="intro__actions">
            {session ? (
              <Link to="/dashboard" className="button">
                Open your dashboard
              </Link>
            ) : (
              <a href={loginUrl()} className="button button--twitch">
                Log in with Twitch
              </a>
            )}
            <a
              href="https://openplanet.dev/plugin/wheelgpt-plugin"
              className="button button--secondary"
            >
              Get the plugin
            </a>
          </div>
        </div>
        <ChatSample />
      </section>

      <section className="wrap section" aria-labelledby="commands">
        <h2 id="commands">Chat commands</h2>
        <p className="section__lede">
          Times are written as <code>h:mm:ss.xxx</code>. Hours and minutes are optional, the
          milliseconds are not. A perfect guess is worth 5 points, the closest guess 1 point.
        </p>
        <Commands />
      </section>

      <section className="wrap section" id="streamers" aria-labelledby="streamers-title">
        <h2 id="streamers-title">Streamers using WheelGPT</h2>
        <Streamers />
      </section>
    </>
  );
}

type Badge = "mod" | "vip" | "sub" | "streamer";

const badgeNames: Record<Badge, string> = {
  mod: "Moderator",
  vip: "VIP",
  sub: "Subscriber",
  streamer: "Broadcaster",
};

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
  return (
    <figure className="chat" aria-label="Example of a Twitch chat with WheelGPT">
      <ol className="chat__lines">
        <Line user="kazoo_k" color="#c0392b" badges={["sub"]}>
          !guess 47.912
        </Line>
        <Line user="marbled" color="#2471a3">
          !g 48.3 <span className="chat__hint">(needs ms)</span>
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          @marbled hh:mm:ss.xxx (hours and minutes are optional)
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
          NEW PERSONAL BEST 47.881 <Emote name="dinkDonk" file="dink_donk" /> That's{" "}
          <Emote name="AuthorTime" file="author_medal" /> <Emote name="ok" file="ok" />
        </Line>
        <Line user="WheelGPT" badges={["mod"]} bot>
          Nobody guessed it correctly but @tinytires guessed 47.850 (-0.031){" "}
          <Emote name="ok" file="ok" /> (+1 point)
        </Line>
      </ol>
    </figure>
  );
}

function Commands() {
  const { data, error, isPending, refetch } = useCommands();
  if (isPending) return <Loading />;
  if (error) return <ErrorNotice error={error} retry={() => void refetch()} />;
  return (
    <div className="table-scroll">
      <table className="table commands">
        <thead>
          <tr>
            <th scope="col">Command</th>
            <th scope="col">What it does</th>
            <th scope="col">Who</th>
          </tr>
        </thead>
        <tbody>
          {data.map((c) => (
            <tr key={c.name}>
              <td>
                <code>{c.usage}</code>
                <OtherNames command={c} />
              </td>
              <td>{c.description}</td>
              <td className="muted">{c.access === "everyone" ? "everyone" : `${c.access}s`}</td>
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
  const others = [command.name, ...command.aliases].filter((n) => n !== shown);
  if (others.length === 0) return null;
  return <span className="muted small"> also {others.map((n) => `!${n}`).join(", ")}</span>;
}

function Streamers() {
  const channels = useChannels();
  const stats = useStats();

  if (channels.isPending) return <Loading />;
  if (channels.error) {
    return <ErrorNotice error={channels.error} retry={() => void channels.refetch()} />;
  }
  const list = channels.data.items;
  return (
    <>
      {stats.data && (
        <p className="section__lede">{plural(stats.data.channels, "channel")} use WheelGPT.</p>
      )}
      {list.length === 0 ? (
        <p className="muted">Nobody is listed yet. You could be the first.</p>
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
