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
import { useI18n } from "../i18n";
import { formatDate } from "../lib/format";

export function Dashboard() {
  const me = useMe();
  const { t } = useI18n();

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
            <Link to={`/c/${channel.login}`}>{t.dashboard.publicBoard}</Link> ·{" "}
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
  const { t } = useI18n();
  return (
    <div className="wrap section narrow">
      <h1>{t.dashboard.loggedOutTitle}</h1>
      <p className="lede">{t.dashboard.loggedOutLede}</p>
      <a href={loginUrl("/dashboard")} className="button button--twitch">
        {t.layout.login}
      </a>
    </div>
  );
}

function LogoutButton() {
  const logout = useLogout();
  const navigate = useNavigate();
  const { t } = useI18n();
  return (
    <button
      type="button"
      className="button button--secondary button--small dash-head__logout"
      disabled={logout.isPending}
      onClick={() => logout.mutate(undefined, { onSettled: () => void navigate("/") })}
    >
      {t.dashboard.logout}
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
  const { t } = useI18n();
  return (
    <Section title={t.dashboard.devTitle} tone="dev" intro={t.dashboard.devIntro}>
      <p className="status-line">
        <span className={`dot dot--${channel.isLive ? "live" : "disabled"}`} aria-hidden="true" />
        {channel.isLive ? t.dashboard.devLive : t.dashboard.devOffline}
      </p>
      <p className="muted small">
        {channel.isLive ? t.dashboard.devLiveHelp : t.dashboard.devOfflineHelp}
      </p>
      {live.error && <ErrorNotice error={live.error} />}
      <button
        type="button"
        className={`button${channel.isLive ? " button--secondary" : ""}`}
        disabled={live.isPending}
        onClick={() => live.mutate(!channel.isLive)}
      >
        {channel.isLive ? t.dashboard.endStream : t.dashboard.goLive}
      </button>
    </Section>
  );
}

// ---------------------------------------------------------------- bot

function relogin(text: string) {
  return <a href={loginUrl("/dashboard")}>{text}</a>;
}

function BotSection({ bot }: { bot: BotConnection }) {
  const toggle = useBotConnection();
  const status = bot.enabled ? bot.chatSubscription : "disabled";
  const { t } = useI18n();

  return (
    <Section title={t.dashboard.botTitle} intro={t.dashboard.botIntro}>
      <p className="status-line">
        <span className={`dot dot--${status}`} aria-hidden="true" />
        {t.dashboard.chatStatus[status] ?? status}
      </p>
      {status === "off" && (
        <p className="muted small">{t.dashboard.eventsubOff(<code>EVENTSUB_SYNC</code>)}</p>
      )}
      {status === "unknown" && <p className="muted small">{t.dashboard.tryAgainLater}</p>}
      {bot.enabled && bot.chatSubscription === "needs_authorization" && (
        <Notice tone="warn" title={t.dashboard.oneMoreStep}>
          {t.dashboard.needsAuth(<code>/mod wheelgpt</code>, relogin)}
        </Notice>
      )}
      {!bot.authorized && bot.chatSubscription !== "needs_authorization" && (
        <p className="muted small">{t.dashboard.legacyAuth(relogin)}</p>
      )}
      {toggle.error && <ErrorNotice error={toggle.error} />}
      <button
        type="button"
        className={`button${bot.enabled ? " button--secondary" : ""}`}
        disabled={toggle.isPending}
        onClick={() => toggle.mutate(!bot.enabled)}
      >
        {toggle.isPending
          ? t.dashboard.working
          : bot.enabled
            ? t.dashboard.disconnect
            : t.dashboard.connect}
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
  const { t, lang } = useI18n();

  return (
    <Section
      title={t.dashboard.tokenTitle}
      intro={t.dashboard.tokenIntro(<em>{t.dashboard.tokenPath}</em>)}
    >
      {create.data ? (
        <>
          <CopyField label={t.dashboard.newToken} value={create.data.token} secret />
          <p className="muted small">{t.dashboard.copyNow}</p>
        </>
      ) : info.isPending ? (
        <Loading />
      ) : exists ? (
        <p className="muted small">
          {t.dashboard.tokenCreated(
            info.data?.createdAt ? formatDate(info.data.createdAt, lang) : "",
            info.data?.lastUsedAt ? formatDate(info.data.lastUsedAt, lang) : undefined,
          )}
        </p>
      ) : legacy ? (
        <p className="muted small">{t.dashboard.legacyToken}</p>
      ) : (
        <p className="muted small">{t.dashboard.noToken}</p>
      )}
      {create.error && <ErrorNotice error={create.error} />}
      <button
        type="button"
        className="button button--secondary"
        disabled={create.isPending}
        onClick={() => {
          if (exists && !confirm(t.dashboard.confirmNewToken)) return;
          create.mutate();
        }}
      >
        {exists || legacy ? t.dashboard.createNewToken : t.dashboard.createToken}
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
  const { t } = useI18n();

  const submit = (e: FormEvent) => {
    e.preventDefault();
    update.mutate(draft, { onSuccess: (next) => setDraft(next) });
  };

  return (
    <Section title={t.dashboard.settingsTitle} intro={t.dashboard.settingsIntro}>
      <form className="form" onSubmit={submit}>
        <div className="field">
          <label htmlFor="delay">{t.dashboard.delay}</label>
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
            <span>{t.dashboard.seconds}</span>
          </div>
          <p id="delay-help" className="field__help">
            {t.dashboard.delayHelp}
          </p>
          {fields.guessDelaySeconds && <p className="field__error">{fields.guessDelaySeconds}</p>}
        </div>

        <div className="field">
          <label htmlFor="minage">{t.dashboard.minAge}</label>
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
            <span>{t.dashboard.seconds}</span>
          </div>
          <p id="minage-help" className="field__help">
            {t.dashboard.minAgeHelp}
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
            <strong>{t.dashboard.listPublic}</strong>
            <span className="field__help">{t.dashboard.listPublicHelp}</span>
          </span>
        </label>

        {update.error && Object.keys(fields).length === 0 && <ErrorNotice error={update.error} />}
        <div className="row">
          <button type="submit" className="button" disabled={!dirty || update.isPending}>
            {update.isPending ? t.dashboard.saving : t.dashboard.save}
          </button>
          {update.isSuccess && !dirty && <span className="muted small">{t.dashboard.saved}</span>}
        </div>
      </form>
    </Section>
  );
}

function DangerZone({ login }: { login: string }) {
  const [confirmText, setConfirmText] = useState("");
  const remove = useDeleteChannel();
  const navigate = useNavigate();
  const { t } = useI18n();

  return (
    <Section title={t.dashboard.deleteTitle} tone="danger" intro={t.dashboard.deleteIntro}>
      <p className="small">{t.dashboard.typeToConfirm(<strong>{login}</strong>)}</p>
      <form
        className="row"
        onSubmit={(e) => {
          e.preventDefault();
          remove.mutate(undefined, { onSuccess: () => void navigate("/") });
        }}
      >
        <label className="visually-hidden" htmlFor="confirm-delete">
          {t.dashboard.channelName}
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
          {t.dashboard.deleteAll}
        </button>
      </form>
      {remove.error && <ErrorNotice error={remove.error} />}
    </Section>
  );
}
