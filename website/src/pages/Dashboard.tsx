import { useState, type FormEvent, type ReactNode } from "react";
import { Link, useNavigate } from "react-router";
import { ApiError, loginUrl } from "../api/client";
import {
  useBotConnection,
  useCreatePluginToken,
  useDeleteChannel,
  useLiveSwitch,
  useLogout,
  useMe,
  usePluginToken,
  useUpdateSettings,
} from "../api/hooks";
import type { BotConnection, Channel, Settings } from "../api/types";
import { Avatar, CopyField, ErrorNotice, LiveBadge, Loading, Notice } from "../components/ui";
import { formatDate } from "../lib/format";

export function Dashboard() {
  const me = useMe();

  if (me.isPending) {
    return (
      <div className="wrap section">
        <Loading />
      </div>
    );
  }
  if (me.error) {
    return (
      <div className="wrap section">
        <ErrorNotice error={me.error} retry={() => void me.refetch()} />
      </div>
    );
  }
  if (!me.data) return <LoggedOut />;

  const { channel, bot, settings, devMode } = me.data;
  return (
    <div className="wrap dashboard">
      <header className="dash-head">
        <Avatar src={channel.profileImage} name={channel.displayName} size={64} />
        <div>
          <h1>{channel.displayName}</h1>
          <p className="muted">
            <LiveBadge live={channel.isLive} /> ·{" "}
            <Link to={`/c/${channel.login}`}>Public leaderboard</Link> ·{" "}
            <a href={`https://twitch.tv/${channel.login}`}>twitch.tv/{channel.login}</a>
          </p>
        </div>
        <LogoutButton />
      </header>

      <div className="dash-sections">
        {devMode && <DevMode channel={channel} />}
        <BotSection bot={bot} />
        <TokenSection />
        <SettingsSection settings={settings} />
        <DangerZone login={channel.login} />
      </div>
    </div>
  );
}

function LoggedOut() {
  return (
    <div className="wrap section narrow">
      <h1>Your dashboard</h1>
      <p className="lede">
        Log in with your Twitch account to connect the bot, get your plugin token and change how the
        guessing game works in your channel.
      </p>
      <a href={loginUrl("/dashboard")} className="button button--twitch">
        Log in with Twitch
      </a>
    </div>
  );
}

function LogoutButton() {
  const logout = useLogout();
  const navigate = useNavigate();
  return (
    <button
      type="button"
      className="button button--secondary button--small dash-head__logout"
      disabled={logout.isPending}
      onClick={() => logout.mutate(undefined, { onSettled: () => void navigate("/") })}
    >
      Log out
    </button>
  );
}

/** One row of the dashboard: what it is on the left, the controls on the right. */
function Section({
  title,
  intro,
  tone,
  children,
}: {
  title: string;
  intro?: ReactNode;
  tone?: "dev" | "danger";
  children: ReactNode;
}) {
  return (
    <section className={`dash-section${tone ? ` dash-section--${tone}` : ""}`} aria-label={title}>
      <div className="dash-section__about">
        <h2>{title}</h2>
        {intro && <p>{intro}</p>}
      </div>
      <div className="dash-section__body">{children}</div>
    </section>
  );
}

// ---------------------------------------------------------------- dev mode

function DevMode({ channel }: { channel: Channel }) {
  const live = useLiveSwitch();
  return (
    <Section
      title="Dev mode"
      tone="dev"
      intro="This server runs in dev mode. Pretend to be live to test the plugin without streaming."
    >
      <p className="status-line">
        <span className={`dot dot--${channel.isLive ? "live" : "disabled"}`} aria-hidden="true" />
        {channel.isLive ? "Pretending to be live" : "Offline"}
      </p>
      <p className="muted small">
        {channel.isLive
          ? "The plugin's map, room and personal bests are accepted. Ending the stream clears them like a real stream end."
          : "While offline the plugin's data is rejected."}
      </p>
      {live.error && <ErrorNotice error={live.error} />}
      <button
        type="button"
        className={`button${channel.isLive ? " button--secondary" : ""}`}
        disabled={live.isPending}
        onClick={() => live.mutate(!channel.isLive)}
      >
        {channel.isLive ? "End stream" : "Go live"}
      </button>
    </Section>
  );
}

// ---------------------------------------------------------------- bot

const chatStatus: Record<string, string> = {
  enabled: "Reading your chat",
  unknown: "Twitch didn't answer",
  needs_authorization: "Twitch refused access to your chat",
  disabled: "Disconnected",
  off: "Connected, but not listening to chat",
};

function BotSection({ bot }: { bot: BotConnection }) {
  const toggle = useBotConnection();
  const status = bot.enabled ? bot.chatSubscription : "disabled";

  return (
    <Section title="Bot" intro="WheelGPT answers in your Twitch chat as long as it is connected.">
      <p className="status-line">
        <span className={`dot dot--${status}`} aria-hidden="true" />
        {chatStatus[status] ?? status}
      </p>
      {status === "off" && (
        <p className="muted small">
          This server doesn't subscribe to Twitch chats (<code>EVENTSUB_SYNC</code> is off), so chat
          commands don't reach it. That's normal for local development: Twitch can only deliver chat
          to a public https address.
        </p>
      )}
      {status === "unknown" && <p className="muted small">Try connecting again in a moment.</p>}
      {bot.enabled && bot.chatSubscription === "needs_authorization" && (
        <Notice tone="warn" title="One more step">
          Log in again so Twitch lets the bot into your chat, or type <code>/mod wheelgpt</code> in
          your chat. <a href={loginUrl("/dashboard")}>Log in again</a>
        </Notice>
      )}
      {!bot.authorized && bot.chatSubscription !== "needs_authorization" && (
        <p className="muted small">
          You signed up before the new login. The bot keeps working as long as it is a moderator in
          your channel; <a href={loginUrl("/dashboard")}>log in again</a> to switch over for good.
        </p>
      )}
      {toggle.error && <ErrorNotice error={toggle.error} />}
      <button
        type="button"
        className={`button${bot.enabled ? " button--secondary" : ""}`}
        disabled={toggle.isPending}
        onClick={() => toggle.mutate(!bot.enabled)}
      >
        {toggle.isPending ? "Working…" : bot.enabled ? "Disconnect the bot" : "Connect the bot"}
      </button>
    </Section>
  );
}

// ---------------------------------------------------------------- token

function TokenSection() {
  const info = usePluginToken();
  const create = useCreatePluginToken();
  const exists = info.data?.exists ?? false;
  const legacy = info.data?.legacyTokenActive ?? false;

  return (
    <Section
      title="Plugin token"
      intro={
        <>
          The Openplanet plugin sends your map, room and personal bests with it. Paste it in{" "}
          <em>Settings → WheelGPT → Token</em> and keep it to yourself.
        </>
      }
    >
      {create.data ? (
        <>
          <CopyField label="Your new token (shown once)" value={create.data.token} secret />
          <p className="muted small">Copy it now. We only store a fingerprint of it.</p>
        </>
      ) : info.isPending ? (
        <Loading />
      ) : exists ? (
        <p className="muted small">
          Created {info.data?.createdAt && formatDate(info.data.createdAt)}
          {info.data?.lastUsedAt
            ? `, last used ${formatDate(info.data.lastUsedAt)}`
            : ", not used yet"}
          .
        </p>
      ) : legacy ? (
        <p className="muted small">
          Your token from the old website still works. A new token replaces it.
        </p>
      ) : (
        <p className="muted small">You don't have a token yet.</p>
      )}
      {create.error && <ErrorNotice error={create.error} />}
      <button
        type="button"
        className="button button--secondary"
        disabled={create.isPending}
        onClick={() => {
          if (exists && !confirm("The current token stops working. Create a new one?")) return;
          create.mutate();
        }}
      >
        {exists || legacy ? "Create a new token" : "Create a token"}
      </button>
    </Section>
  );
}

// ---------------------------------------------------------------- settings

function SettingsSection({ settings }: { settings: Settings }) {
  const [draft, setDraft] = useState(settings);
  const update = useUpdateSettings();
  const dirty = JSON.stringify(draft) !== JSON.stringify(settings);
  const error = update.error instanceof ApiError ? update.error : undefined;
  const fields = error?.fieldErrors ?? {};

  const submit = (e: FormEvent) => {
    e.preventDefault();
    update.mutate(draft, { onSuccess: (next) => setDraft(next) });
  };

  return (
    <Section title="Guessing" intro="How the guessing game works in your channel.">
      <form className="form" onSubmit={submit}>
        <div className="field">
          <label htmlFor="delay">Announcement delay</label>
          <div className="field__input">
            <input
              id="delay"
              type="number"
              min={0}
              max={60}
              step={0.5}
              value={draft.guessDelaySeconds}
              onChange={(e) =>
                setDraft({ ...draft, guessDelaySeconds: e.currentTarget.valueAsNumber || 0 })
              }
              aria-describedby="delay-help"
            />
            <span>seconds</span>
          </div>
          <p id="delay-help" className="field__help">
            Wait this long before posting a PB, so chat doesn't see it before your stream does.
            Match it to your stream delay.
          </p>
          {fields.guessDelaySeconds && <p className="field__error">{fields.guessDelaySeconds}</p>}
        </div>

        <div className="field">
          <label htmlFor="minage">Last-second guesses</label>
          <div className="field__input">
            <input
              id="minage"
              type="number"
              min={0}
              max={60}
              step={1}
              value={draft.guessMinAgeSeconds}
              onChange={(e) =>
                setDraft({ ...draft, guessMinAgeSeconds: e.currentTarget.valueAsNumber || 0 })
              }
              aria-describedby="minage-help"
            />
            <span>seconds</span>
          </div>
          <p id="minage-help" className="field__help">
            Ignore guesses placed or changed less than this long before the PB. Stops people from
            guessing while they watch you cross the finish.
          </p>
          {fields.guessMinAgeSeconds && <p className="field__error">{fields.guessMinAgeSeconds}</p>}
        </div>

        <label className="check">
          <input
            type="checkbox"
            checked={draft.public}
            onChange={(e) => setDraft({ ...draft, public: e.currentTarget.checked })}
          />
          <span>
            <strong>List my channel on wheelgpt.dev</strong>
            <span className="field__help">
              Your leaderboard stays reachable by link either way.
            </span>
          </span>
        </label>

        {update.error && Object.keys(fields).length === 0 && <ErrorNotice error={update.error} />}
        <div className="row">
          <button type="submit" className="button" disabled={!dirty || update.isPending}>
            {update.isPending ? "Saving…" : "Save settings"}
          </button>
          {update.isSuccess && !dirty && <span className="muted small">Saved.</span>}
        </div>
      </form>
    </Section>
  );
}

function DangerZone({ login }: { login: string }) {
  const [confirmText, setConfirmText] = useState("");
  const remove = useDeleteChannel();
  const navigate = useNavigate();

  return (
    <Section
      title="Remove WheelGPT"
      tone="danger"
      intro="Removes the bot from your channel and deletes your settings, token and leaderboards. This can't be undone."
    >
      <p className="small">
        Type <strong>{login}</strong> to confirm.
      </p>
      <form
        className="row"
        onSubmit={(e) => {
          e.preventDefault();
          remove.mutate(undefined, { onSuccess: () => void navigate("/") });
        }}
      >
        <label className="visually-hidden" htmlFor="confirm-delete">
          Channel name
        </label>
        <input
          id="confirm-delete"
          value={confirmText}
          onChange={(e) => setConfirmText(e.currentTarget.value)}
          autoComplete="off"
          placeholder={login}
        />
        <button
          type="submit"
          className="button button--danger"
          disabled={confirmText !== login || remove.isPending}
        >
          Delete everything
        </button>
      </form>
      {remove.error && <ErrorNotice error={remove.error} />}
    </Section>
  );
}
