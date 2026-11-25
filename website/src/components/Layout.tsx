import { Link, NavLink, Outlet, useLocation } from "react-router";
import { useEffect } from "react";
import { loginUrl } from "../api/client";
import { useMe } from "../api/hooks";
import { Avatar } from "./ui";

export function Layout() {
  const { data: session, isPending } = useMe();
  const { pathname, hash } = useLocation();

  useEffect(() => scrollToTarget(pathname, hash), [pathname, hash]);

  return (
    <div className="page">
      <a className="skip" href="#main">
        Skip to content
      </a>
      <header className="site-header">
        <div className="wrap site-header__inner">
          <Link to="/" className="wordmark" aria-label="WheelGPT home">
            WheelGPT
          </Link>
          <nav className="site-nav" aria-label="Main">
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
                Log in with Twitch
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
          <p>
            Made for the Trackmania community by{" "}
            <a href="https://github.com/SoWieMarkus">SoWieMarkus</a>. Not affiliated with Ubisoft
            Nadeo or Twitch.
          </p>
          <ul className="site-footer__links">
            <li>
              <a href="https://github.com/SoWieMarkus/wheelgpt">Source code</a>
            </li>
            <li>
              <a href="https://openplanet.dev/plugin/wheelgpt-plugin">Openplanet plugin</a>
            </li>
            <li>
              <a href="/api/docs/">API</a>
            </li>
            <li>
              <Link to="/imprint">Imprint</Link>
            </li>
            <li>
              <Link to="/privacy">Privacy</Link>
            </li>
          </ul>
        </div>
      </footer>
    </div>
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
