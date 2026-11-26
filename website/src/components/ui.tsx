import { useState, type ReactNode } from "react";
import type { Medals as MedalTimes } from "../api/types";
import { useI18n } from "../i18n";
import { formatTime } from "../lib/format";

export function Avatar({ src, name, size = 40 }: { src?: string; name: string; size?: number }) {
  if (!src) {
    return (
      <span
        className="avatar avatar--empty"
        style={{ width: size, height: size }}
        aria-hidden="true"
      >
        {name.slice(0, 1).toUpperCase()}
      </span>
    );
  }
  return <img className="avatar" src={src} alt="" width={size} height={size} loading="lazy" />;
}

export function LiveBadge({ live }: { live: boolean }) {
  const { t } = useI18n();
  if (!live) return <span className="live live--off">{t.ui.offline}</span>;
  return <span className="live">{t.ui.live}</span>;
}

export function Time({ ms, className = "" }: { ms: number; className?: string }) {
  return <span className={`time ${className}`}>{formatTime(ms)}</span>;
}

const medalNames = ["champion", "author", "gold", "silver", "bronze"] as const;

export function Medals({ medals }: { medals: MedalTimes }) {
  const { t } = useI18n();
  return (
    <dl className="medals">
      {medalNames.map((name) => {
        const ms = medals[name];
        if (ms === undefined) return null;
        return (
          <div key={name} className="medals__item">
            <dt>
              <span className={`medal medal--${name}`} aria-hidden="true" />
              {t.ui.medals[name]}
            </dt>
            <dd>
              <Time ms={ms} />
            </dd>
          </div>
        );
      })}
    </dl>
  );
}

export function Notice({
  tone = "info",
  title,
  children,
}: {
  tone?: "info" | "warn" | "error" | "ok";
  title?: string;
  children: ReactNode;
}) {
  return (
    <div className={`notice notice--${tone}`} role={tone === "error" ? "alert" : "status"}>
      {title && <strong className="notice__title">{title}</strong>}
      <div>{children}</div>
    </div>
  );
}

export function Loading({ label }: { label?: string }) {
  const { t } = useI18n();
  return (
    <p className="loading" aria-live="polite">
      {label ?? t.ui.loading}
    </p>
  );
}

export function ErrorNotice({ error, retry }: { error: Error; retry?: () => void }) {
  const { t } = useI18n();
  return (
    <Notice tone="error" title={t.ui.errorTitle}>
      {error.message}{" "}
      {retry && (
        <button type="button" className="link-button" onClick={retry}>
          {t.ui.retry}
        </button>
      )}
    </Notice>
  );
}

/** A read-only value with a copy button. Secrets stay masked until the eye is clicked. */
export function CopyField({
  value,
  label,
  secret = false,
}: {
  value: string;
  label: string;
  secret?: boolean;
}) {
  const { t } = useI18n();
  const [copied, setCopied] = useState(false);
  const [shown, setShown] = useState(!secret);
  return (
    <div className="copy">
      <label className="copy__label">
        {label}
        <span className="copy__field">
          <input
            className={`copy__input${secret ? " copy__input--secret" : ""}`}
            type={shown ? "text" : "password"}
            readOnly
            value={value}
            autoComplete="off"
            spellCheck={false}
            onFocus={(e) => shown && e.currentTarget.select()}
          />
          {secret && (
            <button
              type="button"
              className="copy__eye"
              aria-label={shown ? t.ui.hideToken : t.ui.showToken}
              aria-pressed={shown}
              title={shown ? t.ui.hide : t.ui.show}
              onClick={() => setShown((v) => !v)}
            >
              <EyeIcon crossed={shown} />
            </button>
          )}
        </span>
      </label>
      <button
        type="button"
        className="button button--secondary"
        onClick={() => {
          void navigator.clipboard.writeText(value).then(() => {
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
          });
        }}
      >
        {copied ? t.ui.copied : t.ui.copy}
      </button>
    </div>
  );
}

function EyeIcon({ crossed }: { crossed: boolean }) {
  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7S2 12 2 12Z" />
      <circle cx="12" cy="12" r="3" />
      {crossed && <path d="M4 4l16 16" />}
    </svg>
  );
}

export function Emote({ name, file }: { name: string; file: string }) {
  return (
    <img
      className="emote"
      src={`/emotes/${file}.avif`}
      alt={name}
      title={name}
      width={28}
      height={28}
    />
  );
}
