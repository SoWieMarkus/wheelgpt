import type { ReactNode } from "react";
import { Mail, type Messages } from "./en";

const n = (v: number) => v.toLocaleString("de");

export const de: Messages = {
  locale: "de",
  meta: {
    title: "WheelGPT – ein Trackmania-Twitch-Bot",
    description:
      "WheelGPT ist ein Twitch-Bot für Trackmania: Zuschauer tippen deine nächste persönliche Bestzeit, fragen nach der aktuellen Map und klettern in einer monatlichen Rangliste nach oben.",
  },

  layout: {
    skip: "Zum Inhalt springen",
    home: "WheelGPT-Startseite",
    mainNav: "Hauptmenü",
    login: "Mit Twitch anmelden",
    language: "Sprache",
    footer: (author: ReactNode) => (
      <>
        Für die Trackmania-Community gemacht von {author}. Nicht verbunden mit Ubisoft Nadeo oder
        Twitch.
      </>
    ),
    source: "Quellcode",
    plugin: "Openplanet-Plugin",
    api: "API",
    imprint: "Impressum",
    privacy: "Datenschutz",
  },

  ui: {
    offline: "offline",
    live: "live",
    loading: "Lädt…",
    errorTitle: "Das hat nicht geklappt.",
    retry: "Nochmal versuchen",
    hideToken: "Token verbergen",
    showToken: "Token anzeigen",
    hide: "Verbergen",
    show: "Anzeigen",
    copy: "Kopieren",
    copied: "Kopiert",
    medals: {
      champion: "Champion",
      author: "Autor",
      gold: "Gold",
      silver: "Silber",
      bronze: "Bronze",
    },
  },

  periods: {
    allTime: "Gesamt",
    thisMonth: "Dieser Monat",
  },

  landing: {
    eyebrow: "Ein Twitch-Bot für Trackmania-Streams",
    lede: "Jaaa hello guys! My name is WheelGPT! And especially mothers love me!",
    intro: (plugin: ReactNode) => (
      <>
        Ich bin ein Twitch-Bot, speziell für Trackmania-Streamer gemacht. Mit mir können deine
        Zuschauer jederzeit nachsehen, welche Map du gerade fährst. Außerdem können sie deine
        nächste persönliche Bestzeit tippen. Damit alles reibungslos läuft, solltest du das {plugin}{" "}
        installieren. Neugierig? Dann hol dir WheelGPT jetzt in deinen Kanal!
      </>
    ),
    pluginName: "OpenPlanet-Plugin",
    openDashboard: "Zum Dashboard",
    getPlugin: "Plugin holen",
    commandsTitle: "Chat-Befehle",
    commandsLede: (format: ReactNode) => (
      <>
        Zeiten werden als {format} geschrieben. Stunden und Minuten sind optional, die Millisekunden
        nicht. Ein perfekter Tipp bringt 5 Punkte, der knappste Tipp 1 Punkt.
      </>
    ),
    streamersTitle: "Streamer mit WheelGPT",
    channelCount: (c: number) => `${n(c)} ${c === 1 ? "Kanal nutzt" : "Kanäle nutzen"} WheelGPT.`,
    noStreamers: "Noch ist niemand gelistet. Du könntest der Erste sein.",
    badges: {
      mod: "Moderator",
      vip: "VIP",
      sub: "Abonnent",
      streamer: "Streamer",
    },
    chatLabel: "Beispiel eines Twitch-Chats mit WheelGPT",
    chat: {
      needsMs: "(ohne ms)",
      format: "hh:mm:ss.xxx (Stunden und Minuten sind optional)",
      pb: "NEUE PERSÖNLICHE BESTZEIT 47.881",
      thats: "Das ist",
      nobody: "Niemand hat sie genau getroffen, aber @tinytires hat 47.850 getippt (-0.031)",
      point: "(+1 Punkt)",
    },
    table: {
      command: "Befehl",
      what: "Was er macht",
      who: "Wer",
      also: "auch",
    },
    access: {
      everyone: "alle",
      subscriber: "Abonnenten",
      vip: "VIPs",
      moderator: "Moderatoren",
      broadcaster: "Streamer",
    },
    commands: {
      wheelgpt: "Was ist das für ein Bot?",
      "wgpt-emotes": "Listet die 7TV-Emotes, die WheelGPT benutzt. Füg sie deinem Kanal hinzu.",
      format: "Erklärt das Zeitformat für Tipps.",
      guess: "Tippe die nächste persönliche Bestzeit des Streamers auf der aktuellen Map.",
      "wgpt-leaderboard":
        "Die besten Tipper des Monats oder aller Zeiten, oder die Platzierung eines Zuschauers.",
      map: "Zeigt die aktuelle Map mit Medaillen und Links.",
      myguess: "Zeigt deinen aktuellen Tipp.",
      resetguesses: "Löscht alle offenen Tipps.",
      result: "Beendet die Runde von Hand mit der Endzeit.",
      room: "Zeigt den Server-Raum, in dem der Streamer gerade spielt.",
    },
  },

  channel: {
    notFoundTitle: "Kanal nicht gefunden",
    notFound: (name: ReactNode, link: (text: string) => ReactNode) => (
      <>
        {name} nutzt WheelGPT nicht. {link("Hier siehst du, wer es nutzt.")}
      </>
    ),
    eyebrow: "Tipp-Rangliste",
    currentMap: "Aktuelle Map",
    driving: "Fährt gerade",
    by: (author: string) => `von ${author}`,
    room: "Raum",
    period: "Zeitraum",
    olderMonths: "Ältere Monate",
    older: "Älter…",
    findViewer: "Zuschauer suchen",
    nobodyCalled: (q: string) => `Niemand namens „${q}“ hat in diesem Zeitraum Punkte.`,
    noPointsMonth:
      "Diesen Monat noch keine Punkte. Die erste persönliche Bestzeit des Monats startet das Rennen.",
    noPoints: "Keine Punkte in diesem Zeitraum.",
    summary: (c: number) =>
      `${n(c)} Tipper mit Punkten · perfekter Tipp 5 Punkte, knappster Tipp 1 Punkt`,
    viewer: "Zuschauer",
    points: "Punkte",
    perfect: "Perfekt",
    pages: "Seiten",
    previous: "Zurück",
    next: "Weiter",
    pageOf: (page: number, pages: number) => `Seite ${page} von ${pages}`,
  },

  dashboard: {
    publicBoard: "Öffentliche Rangliste",
    loggedOutTitle: "Dein Dashboard",
    loggedOutLede:
      "Melde dich mit deinem Twitch-Konto an, um den Bot zu verbinden, deinen Plugin-Token zu holen und das Tippspiel in deinem Kanal einzustellen.",
    logout: "Abmelden",

    devTitle: "Dev-Modus",
    devIntro:
      "Dieser Server läuft im Dev-Modus. Tu so, als wärst du live, um das Plugin ohne Stream zu testen.",
    devLive: "Tut so, als wäre er live",
    devOffline: "Offline",
    devLiveHelp:
      "Map, Raum und persönliche Bestzeiten vom Plugin werden angenommen. Den Stream zu beenden löscht sie wie ein echtes Stream-Ende.",
    devOfflineHelp: "Solange du offline bist, werden die Daten des Plugins abgelehnt.",
    endStream: "Stream beenden",
    goLive: "Live gehen",

    botTitle: "Bot",
    botIntro: "WheelGPT antwortet in deinem Twitch-Chat, solange er verbunden ist.",
    chatStatus: {
      enabled: "Liest deinen Chat",
      unknown: "Twitch hat nicht geantwortet",
      needs_authorization: "Twitch verweigert den Zugriff auf deinen Chat",
      disabled: "Getrennt",
      off: "Verbunden, hört aber nicht auf den Chat",
    },
    eventsubOff: (env: ReactNode) => (
      <>
        Dieser Server abonniert keine Twitch-Chats ({env} ist aus), deshalb kommen Chat-Befehle
        nicht an. Bei lokaler Entwicklung ist das normal: Twitch kann den Chat nur an eine
        öffentliche https-Adresse schicken.
      </>
    ),
    tryAgainLater: "Versuch es gleich nochmal.",
    oneMoreStep: "Noch ein Schritt",
    needsAuth: (cmd: ReactNode, link: (text: string) => ReactNode) => (
      <>
        Melde dich erneut an, damit Twitch den Bot in deinen Chat lässt, oder schreib {cmd} in
        deinen Chat. {link("Erneut anmelden")}
      </>
    ),
    legacyAuth: (link: (text: string) => ReactNode) => (
      <>
        Du hast dich vor dem neuen Login registriert. Der Bot funktioniert weiter, solange er
        Moderator in deinem Kanal ist; {link("melde dich erneut an")}, um endgültig umzusteigen.
      </>
    ),
    working: "Moment…",
    disconnect: "Bot trennen",
    connect: "Bot verbinden",

    tokenTitle: "Plugin-Token",
    tokenIntro: (path: ReactNode) => (
      <>
        Damit schickt das Openplanet-Plugin deine Map, deinen Raum und deine persönlichen
        Bestzeiten. Füg ihn unter {path} ein und behalte ihn für dich.
      </>
    ),
    tokenPath: "Settings → WheelGPT → Token",
    newToken: "Dein neuer Token (wird nur einmal angezeigt)",
    copyNow: "Kopier ihn jetzt. Wir speichern nur einen Fingerabdruck davon.",
    tokenCreated: (created: string, lastUsed?: string) =>
      `Erstellt am ${created}${lastUsed ? `, zuletzt benutzt am ${lastUsed}` : ", noch nicht benutzt"}.`,
    legacyToken: "Dein Token von der alten Website funktioniert noch. Ein neuer Token ersetzt ihn.",
    noToken: "Du hast noch keinen Token.",
    confirmNewToken: "Der aktuelle Token funktioniert dann nicht mehr. Neuen erstellen?",
    createNewToken: "Neuen Token erstellen",
    createToken: "Token erstellen",

    settingsTitle: "Tippspiel",
    settingsIntro: "So funktioniert das Tippspiel in deinem Kanal.",
    delay: "Verzögerung der Ansage",
    seconds: "Sekunden",
    delayHelp:
      "So lange warten, bevor eine Bestzeit gepostet wird, damit der Chat sie nicht vor deinem Stream sieht. Stell es auf deine Stream-Verzögerung ein.",
    minAge: "Tipps in letzter Sekunde",
    minAgeHelp:
      "Tipps ignorieren, die weniger als so lange vor der Bestzeit abgegeben oder geändert wurden. Verhindert, dass jemand tippt, während er dich ins Ziel fahren sieht.",
    listPublic: "Meinen Kanal auf wheelgpt.dev auflisten",
    listPublicHelp: "Deine Rangliste bleibt so oder so über den Link erreichbar.",
    saving: "Speichert…",
    save: "Einstellungen speichern",
    saved: "Gespeichert.",

    deleteTitle: "WheelGPT entfernen",
    deleteIntro:
      "Entfernt den Bot aus deinem Kanal und löscht deine Einstellungen, deinen Token und deine Ranglisten. Das kann nicht rückgängig gemacht werden.",
    typeToConfirm: (login: ReactNode) => <>Gib {login} ein, um zu bestätigen.</>,
    channelName: "Kanalname",
    deleteAll: "Alles löschen",
  },

  legal: {
    imprint: () => (
      <>
        <h1>Impressum</h1>
        <p>Diese Website wird von Markus Wieland betrieben.</p>
        <p>
          Kontakt: <Mail />
        </p>
        <p>
          Alle Inhalte dieser Website werden ohne Gewähr bereitgestellt. Der Autor übernimmt keine
          Haftung für die Richtigkeit, Vollständigkeit und Aktualität der Inhalte. Die Nutzung der
          Website erfolgt auf eigene Gefahr.
        </p>
      </>
    ),
    privacy: (dashboard: (text: string) => ReactNode) => (
      <>
        <h1>Datenschutz</h1>

        <h2>Streamer-Konten</h2>
        <p>
          Wenn du dich auf wheelgpt.dev anmeldest, speichern wir deine Twitch-Kanal-ID, deinen
          Login, deinen Anzeigenamen und die URL deines Profilbilds, zusammen mit deinen
          Einstellungen. Wir bitten Twitch um eine einzige Berechtigung, <code>channel:bot</code>,
          damit das WheelGPT-Bot-Konto in deinem Kanal Chat lesen und schreiben darf. Den
          Access-Token, den Twitch uns beim Login gibt, benutzen wir einmal, um herauszufinden, wer
          du bist, und verwerfen ihn danach; er wird nie gespeichert.
        </p>
        <p>
          Beim Anmelden wird ein Cookie gesetzt, <code>wgpt_session</code>, das dich bis zu 30 Tage
          angemeldet hält. Es ist technisch notwendig und wird nicht zum Tracking verwendet. Dein
          Plugin-Token wird nur als Fingerabdruck (Hash) gespeichert. Deine Sprachauswahl wird im
          lokalen Speicher deines Browsers gemerkt und nie an uns gesendet.
        </p>

        <h2>Zuschauer</h2>
        <p>
          Wenn du einen Befehl wie <code>!guess</code> benutzt, speichert der Bot deine
          Twitch-Nutzer-ID, deinen Anzeigenamen und deinen Tipp. Nach einer persönlichen Bestzeit
          werden dein Tipp, sein Ergebnis und deine Punkte für die Ranglisten des Kanals und deinen
          Tipp-Verlauf aufbewahrt, die auf dieser Website öffentlich sind. Nachrichten, die keine
          Befehle sind, werden nicht gespeichert.
        </p>

        <h2>Zugriffsprotokoll</h2>
        <p>
          Wir führen Zugriffsprotokolle, die deine IP-Adresse und den User-Agent deines Browsers
          enthalten können, um zu sehen, wie die Seite genutzt wird, und um Missbrauch zu
          verhindern. Sie werden regelmäßig rotiert und in der Regel höchstens 24 Stunden
          aufbewahrt.
        </p>

        <h2>Deine Daten löschen</h2>
        <p>
          Streamer können ihren Kanal und alle zugehörigen Daten im {dashboard("Dashboard")}{" "}
          löschen. Zuschauer, die ihre Tipps entfernen lassen möchten, können an <Mail /> schreiben
          oder @sowiemarkus auf Discord fragen.
        </p>
      </>
    ),
    notFoundTitle: "Von der Strecke abgekommen",
    notFound: (link: (text: string) => ReactNode) => (
      <>Unter dieser Adresse gibt es nichts. {link("Zurück zum Start")}</>
    ),
  },
};
