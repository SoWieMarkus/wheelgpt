import type { ReactNode } from "react";

const n = (v: number) => v.toLocaleString("en");

/**
 * English texts. This object defines the shape every other language must
 * match; TypeScript reports a missing or extra key in the other catalogs.
 */
export const en = {
  locale: "en",
  meta: {
    title: "WheelGPT – a Trackmania Twitch bot",
    description:
      "WheelGPT is a Trackmania Twitch bot: viewers guess your next personal best, ask for the current map and climb a monthly leaderboard.",
  },

  layout: {
    skip: "Skip to content",
    home: "WheelGPT home",
    mainNav: "Main",
    login: "Log in with Twitch",
    language: "Language",
    footer: (author: ReactNode) => (
      <>
        Made for the Trackmania community by {author}. Not affiliated with Ubisoft Nadeo or Twitch.
      </>
    ),
    source: "Source code",
    plugin: "Openplanet plugin",
    api: "API",
    imprint: "Imprint",
    privacy: "Privacy",
  },

  ui: {
    offline: "offline",
    live: "live",
    loading: "Loading…",
    errorTitle: "That didn't work.",
    retry: "Try again",
    hideToken: "Hide token",
    showToken: "Show token",
    hide: "Hide",
    show: "Show",
    copy: "Copy",
    copied: "Copied",
    medals: {
      champion: "champion",
      author: "author",
      gold: "gold",
      silver: "silver",
      bronze: "bronze",
    },
  },

  periods: {
    allTime: "All time",
    thisMonth: "This month",
  },

  landing: {
    eyebrow: "A Twitch bot for Trackmania streams",
    lede: "Jaaa hello guys! My name is WheelGPT! And especially mothers love me!",
    intro: (plugin: ReactNode) => (
      <>
        I am a Twitch bot made especially for Trackmania streamers. With me, your viewers can always
        check which map you are currently playing. They can also guess your next personal best time.
        To make everything run smoothly, you should install the {plugin}. Curious? Then add WheelGPT
        to your channel now!
      </>
    ),
    pluginName: "OpenPlanet Plugin",
    openDashboard: "Open your dashboard",
    getPlugin: "Get the plugin",
    commandsTitle: "Chat commands",
    commandsLede: (format: ReactNode) => (
      <>
        Times are written as {format}. Hours and minutes are optional, the milliseconds are not. A
        perfect guess is worth 5 points, the closest guess 1 point.
      </>
    ),
    streamersTitle: "Streamers using WheelGPT",
    channelCount: (c: number) => `${n(c)} ${c === 1 ? "channel uses" : "channels use"} WheelGPT.`,
    noStreamers: "Nobody is listed yet. You could be the first.",
    badges: {
      mod: "Moderator",
      vip: "VIP",
      sub: "Subscriber",
      streamer: "Broadcaster",
    },
    chatLabel: "Example of a Twitch chat with WheelGPT",
    chat: {
      needsMs: "(needs ms)",
      format: "hh:mm:ss.xxx (hours and minutes are optional)",
      pb: "NEW PERSONAL BEST 47.881",
      thats: "That's",
      nobody: "Nobody guessed it correctly but @tinytires guessed 47.850 (-0.031)",
      point: "(+1 point)",
    },
    table: {
      command: "Command",
      what: "What it does",
      who: "Who",
      also: "also",
    },
    access: {
      everyone: "everyone",
      subscriber: "subscribers",
      vip: "VIPs",
      moderator: "moderators",
      broadcaster: "broadcaster",
    } as Record<string, string>,
    /** Keyed by command name; unknown commands fall back to the API text. */
    commands: {
      wheelgpt: "What is this bot?",
      "wgpt-emotes": "List the 7TV emotes WheelGPT uses. Add them to your channel.",
      format: "Explain the time format for guesses.",
      guess: "Guess the streamer's next personal best on the current map.",
      "wgpt-leaderboard": "Top guessers of this month or of all time, or one viewer's position.",
      map: "Show the current map with its medals and links.",
      myguess: "Show your current guess.",
      resetguesses: "Delete all open guesses.",
      result: "End the round by hand with the final time.",
      room: "Show the server room the streamer is playing in.",
    } as Record<string, string>,
  },

  channel: {
    notFoundTitle: "Channel not found",
    notFound: (name: ReactNode, link: (text: string) => ReactNode) => (
      <>
        {name} doesn't use WheelGPT. {link("See who does.")}
      </>
    ),
    eyebrow: "Guessing leaderboard",
    currentMap: "Current map",
    driving: "Driving right now",
    by: (author: string) => `by ${author}`,
    room: "Room",
    period: "Period",
    olderMonths: "Older months",
    older: "Older…",
    findViewer: "Find a viewer",
    nobodyCalled: (q: string) => `Nobody called “${q}” has points in this period.`,
    noPointsMonth:
      "No points this month yet. The first personal best of the month starts the race.",
    noPoints: "No points in this period.",
    summary: (c: number) =>
      `${n(c)} ${c === 1 ? "guesser" : "guessers"} with points · perfect guess 5 points, closest guess 1 point`,
    viewer: "Viewer",
    points: "Points",
    perfect: "Perfect",
    pages: "Pages",
    previous: "Previous",
    next: "Next",
    pageOf: (page: number, pages: number) => `Page ${page} of ${pages}`,
  },

  dashboard: {
    publicBoard: "Public leaderboard",
    loggedOutTitle: "Your dashboard",
    loggedOutLede:
      "Log in with your Twitch account to connect the bot, get your plugin token and change how the guessing game works in your channel.",
    logout: "Log out",

    devTitle: "Dev mode",
    devIntro:
      "This server runs in dev mode. Pretend to be live to test the plugin without streaming.",
    devLive: "Pretending to be live",
    devOffline: "Offline",
    devLiveHelp:
      "The plugin's map, room and personal bests are accepted. Ending the stream clears them like a real stream end.",
    devOfflineHelp: "While offline the plugin's data is rejected.",
    endStream: "End stream",
    goLive: "Go live",

    botTitle: "Bot",
    botIntro: "WheelGPT answers in your Twitch chat as long as it is connected.",
    chatStatus: {
      enabled: "Reading your chat",
      unknown: "Twitch didn't answer",
      needs_authorization: "Twitch refused access to your chat",
      disabled: "Disconnected",
      off: "Connected, but not listening to chat",
    } as Record<string, string>,
    eventsubOff: (env: ReactNode) => (
      <>
        This server doesn't subscribe to Twitch chats ({env} is off), so chat commands don't reach
        it. That's normal for local development: Twitch can only deliver chat to a public https
        address.
      </>
    ),
    tryAgainLater: "Try connecting again in a moment.",
    oneMoreStep: "One more step",
    needsAuth: (cmd: ReactNode, link: (text: string) => ReactNode) => (
      <>
        Log in again so Twitch lets the bot into your chat, or type {cmd} in your chat.{" "}
        {link("Log in again")}
      </>
    ),
    legacyAuth: (link: (text: string) => ReactNode) => (
      <>
        You signed up before the new login. The bot keeps working as long as it is a moderator in
        your channel; {link("log in again")} to switch over for good.
      </>
    ),
    working: "Working…",
    disconnect: "Disconnect the bot",
    connect: "Connect the bot",

    tokenTitle: "Plugin token",
    tokenIntro: (path: ReactNode) => (
      <>
        The Openplanet plugin sends your map, room and personal bests with it. Paste it in {path}{" "}
        and keep it to yourself.
      </>
    ),
    tokenPath: "Settings → WheelGPT → Token",
    newToken: "Your new token (shown once)",
    copyNow: "Copy it now. We only store a fingerprint of it.",
    tokenCreated: (created: string, lastUsed?: string) =>
      `Created ${created}${lastUsed ? `, last used ${lastUsed}` : ", not used yet"}.`,
    legacyToken: "Your token from the old website still works. A new token replaces it.",
    noToken: "You don't have a token yet.",
    confirmNewToken: "The current token stops working. Create a new one?",
    createNewToken: "Create a new token",
    createToken: "Create a token",

    settingsTitle: "Guessing",
    settingsIntro: "How the guessing game works in your channel.",
    delay: "Announcement delay",
    seconds: "seconds",
    delayHelp:
      "Wait this long before posting a PB, so chat doesn't see it before your stream does. Match it to your stream delay.",
    minAge: "Last-second guesses",
    minAgeHelp:
      "Ignore guesses placed or changed less than this long before the PB. Stops people from guessing while they watch you cross the finish.",
    listPublic: "List my channel on wheelgpt.dev",
    listPublicHelp: "Your leaderboard stays reachable by link either way.",
    saving: "Saving…",
    save: "Save settings",
    saved: "Saved.",

    deleteTitle: "Remove WheelGPT",
    deleteIntro:
      "Removes the bot from your channel and deletes your settings, token and leaderboards. This can't be undone.",
    typeToConfirm: (login: ReactNode) => <>Type {login} to confirm.</>,
    channelName: "Channel name",
    deleteAll: "Delete everything",
  },

  legal: {
    imprint: () => (
      <>
        <h1>Imprint</h1>
        <p>This website is operated by Markus Wieland.</p>
        <p>
          Contact: <Mail />
        </p>
        <p>
          All content on this website is provided without any warranty. The author assumes no
          liability for the correctness, completeness, and topicality of the content. The use of the
          website is at your own risk.
        </p>
      </>
    ),
    privacy: (dashboard: (text: string) => ReactNode) => (
      <>
        <h1>Privacy</h1>

        <h2>Streamer accounts</h2>
        <p>
          When you log in on wheelgpt.dev we store your Twitch channel id, login, display name and
          the URL of your profile picture, together with your settings. We ask Twitch for one
          permission, <code>channel:bot</code>, so the WheelGPT bot account may read and write chat
          in your channel. The access token Twitch gives us during login is used once to find out
          who you are and is then discarded; it is never stored.
        </p>
        <p>
          Logging in sets one cookie, <code>wgpt_session</code>, that keeps you logged in for up to
          30 days. It is strictly necessary and not used for tracking. Your plugin token is stored
          only as a fingerprint (hash). Your language choice is remembered in your browser's local
          storage and never sent to us.
        </p>

        <h2>Viewers</h2>
        <p>
          When you use a command like <code>!guess</code>, the bot stores your Twitch user id,
          display name and your guess. After a personal best your guess, its result and your points
          are kept for the channel's leaderboards and your guess history, which are public on this
          website. Messages that are not commands are not stored.
        </p>

        <h2>Access log</h2>
        <p>
          We keep traffic logs that may contain your IP address and browser user agent, to see how
          the site is used and to stop abuse. They are rotated regularly and usually kept for at
          most 24 hours.
        </p>

        <h2>Deleting your data</h2>
        <p>
          Streamers can delete their channel and all of its data on the {dashboard("dashboard")}.
          Viewers who want their guesses removed can write to <Mail /> or ask @sowiemarkus on
          Discord.
        </p>
      </>
    ),
    notFoundTitle: "Off the track",
    notFound: (link: (text: string) => ReactNode) => (
      <>There's nothing at this address. {link("Back to the start")}</>
    ),
  },
};

export type Messages = typeof en;

export function Mail() {
  return <a href="mailto:markus.u.wieland@gmail.com">markus.u.wieland@gmail.com</a>;
}
