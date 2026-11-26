import { Link, NavLink, Outlet, useLocation } from "react-router";
import { useEffect } from "react";
import { loginUrl } from "../api/client";
import { useMe } from "../api/hooks";
import { languages, useI18n } from "../i18n";
import { Avatar } from "./ui";

export function Layout() {
  const { data: session, isPending } = useMe();
  const { pathname, hash } = useLocation();
  const { t } = useI18n();

  useEffect(() => scrollToTarget(pathname, hash), [pathname, hash]);

  return (
    <div className="page">
      <a className="skip" href="#main">
        {t.layout.skip}
      </a>
      <header className="site-header">
        <div className="wrap site-header__inner">
          <Link to="/" className="wordmark" aria-label={t.layout.home}>
            WheelGPT
          </Link>
          <nav className="site-nav" aria-label={t.layout.mainNav}>
            <LanguageSwitch />
            {session ? (
              <NavLink to="/dashboard" className="site-nav__me">
                <Avatar
                  src={session.channel.profileImage}
                  name={session.channel.displayName}
                  size={24}
                />
                {session.channel.displayName}
              </NavLink>
            ) : !isPending ? (
              <a href={loginUrl()} className="button button--small button--twitch">
                {t.layout.login}
              </a>
            ) : null}
          </nav>
        </div>
      </header>

      <main id="main" className="site-main">
        <Outlet />
      </main>

      <footer className="site-footer">
        <div className="wrap site-footer__inner">
          <p>{t.layout.footer(<a href="https://github.com/SoWieMarkus">SoWieMarkus</a>)}</p>
          <ul className="site-footer__links">
            <li>
              <a href="https://github.com/SoWieMarkus/wheelgpt">{t.layout.source}</a>
            </li>
            <li>
              <a href="https://openplanet.dev/plugin/wheelgpt-plugin">{t.layout.plugin}</a>
            </li>
            <li>
              <a href="/api/docs/">{t.layout.api}</a>
            </li>
            <li>
              <Link to="/imprint">{t.layout.imprint}</Link>
            </li>
            <li>
              <Link to="/privacy">{t.layout.privacy}</Link>
            </li>
          </ul>
        </div>
      </footer>
    </div>
  );
}

function LanguageSwitch() {
  const { lang, setLang, t } = useI18n();
  return (
    <fieldset className="lang-switch">
      <legend className="visually-hidden">{t.layout.language}</legend>
      {languages.map((l) => (
        <button
          key={l}
          type="button"
          lang={l}
          className="lang-switch__option"
          aria-pressed={l === lang}
          onClick={() => setLang(l)}
        >
          {l.toUpperCase()}
        </button>
      ))}
    </fieldset>
  );
}

/** Scrolls to the #hash target (which may render late) or to the top. */
function scrollToTarget(_pathname: string, hash: string): (() => void) | undefined {
  if (!hash) {
    window.scrollTo(0, 0);
    return undefined;
  }
  let tries = 0;
  const timer = setInterval(() => {
    const el = document.getElementById(hash.slice(1));
    if (el || ++tries > 20) {
      el?.scrollIntoView({ behavior: "smooth", block: "start" });
      clearInterval(timer);
    }
  }, 50);
  return () => clearInterval(timer);
}
