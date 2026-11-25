import { Link } from "react-router";

export function Imprint() {
  return (
    <article className="wrap section narrow prose">
      <h1>Imprint</h1>
      <p>This website is operated by Markus Wieland.</p>
      <p>
        Contact: <a href="mailto:markus.u.wieland@gmail.com">markus.u.wieland@gmail.com</a>
      </p>
      <p>
        All content on this website is provided without any warranty. The author assumes no
        liability for the correctness, completeness, and topicality of the content. The use of the
        website is at your own risk.
      </p>
    </article>
  );
}

export function Privacy() {
  return (
    <article className="wrap section narrow prose">
      <h1>Privacy</h1>

      <h2>Streamer accounts</h2>
      <p>
        When you log in on wheelgpt.dev we store your Twitch channel id, login, display name and the
        URL of your profile picture, together with your settings. We ask Twitch for one permission,{" "}
        <code>channel:bot</code>, so the WheelGPT bot account may read and write chat in your
        channel. The access token Twitch gives us during login is used once to find out who you are
        and is then discarded; it is never stored.
      </p>
      <p>
        Logging in sets one cookie, <code>wgpt_session</code>, that keeps you logged in for up to 30
        days. It is strictly necessary and not used for tracking. Your plugin token is stored only
        as a fingerprint (hash).
      </p>

      <h2>Viewers</h2>
      <p>
        When you use a command like <code>!guess</code>, the bot stores your Twitch user id, display
        name and your guess. After a personal best your guess, its result and your points are kept
        for the channel's leaderboards and your guess history, which are public on this website.
        Messages that are not commands are not stored.
      </p>

      <h2>Access log</h2>
      <p>
        We keep traffic logs that may contain your IP address and browser user agent, to see how the
        site is used and to stop abuse. They are rotated regularly and usually kept for at most 24
        hours.
      </p>

      <h2>Deleting your data</h2>
      <p>
        Streamers can delete their channel and all of its data on the{" "}
        <Link to="/dashboard">dashboard</Link>. Viewers who want their guesses removed can write to{" "}
        <a href="mailto:markus.u.wieland@gmail.com">markus.u.wieland@gmail.com</a> or ask
        @sowiemarkus on Discord.
      </p>
    </article>
  );
}

export function NotFound() {
  return (
    <div className="wrap section narrow">
      <h1>Off the track</h1>
      <p>
        There's nothing at this address. <Link to="/">Back to the start</Link>
      </p>
    </div>
  );
}
