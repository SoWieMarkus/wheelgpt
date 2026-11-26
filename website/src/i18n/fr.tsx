import type { ReactNode } from "react";
import { Mail, type Messages } from "./en";

const n = (v: number) => v.toLocaleString("fr");

export const fr: Messages = {
  locale: "fr",
  meta: {
    title: "WheelGPT – un bot Twitch pour Trackmania",
    description:
      "WheelGPT est un bot Twitch pour Trackmania : les viewers devinent ton prochain record personnel, demandent la map en cours et grimpent dans un classement mensuel.",
  },

  layout: {
    skip: "Aller au contenu",
    home: "Accueil WheelGPT",
    mainNav: "Menu principal",
    login: "Se connecter avec Twitch",
    language: "Langue",
    footer: (author: ReactNode) => (
      <>Fait pour la communauté Trackmania par {author}. Non affilié à Ubisoft Nadeo ni à Twitch.</>
    ),
    source: "Code source",
    plugin: "Plugin Openplanet",
    api: "API",
    imprint: "Mentions légales",
    privacy: "Confidentialité",
  },

  ui: {
    offline: "hors ligne",
    live: "en direct",
    loading: "Chargement…",
    errorTitle: "Ça n'a pas marché.",
    retry: "Réessayer",
    hideToken: "Masquer le token",
    showToken: "Afficher le token",
    hide: "Masquer",
    show: "Afficher",
    copy: "Copier",
    copied: "Copié",
    medals: {
      champion: "champion",
      author: "auteur",
      gold: "or",
      silver: "argent",
      bronze: "bronze",
    },
  },

  periods: {
    allTime: "Depuis toujours",
    thisMonth: "Ce mois-ci",
  },

  landing: {
    eyebrow: "Un bot Twitch pour les streams Trackmania",
    lede: "Jaaa hello guys! My name is WheelGPT! And especially mothers love me!",
    intro: (plugin: ReactNode) => (
      <>
        Je suis un bot Twitch conçu spécialement pour les streamers Trackmania. Grâce à moi, tes
        viewers peuvent toujours voir sur quelle map tu joues. Ils peuvent aussi deviner ton
        prochain record personnel. Pour que tout fonctionne bien, installe le {plugin}. Curieux ?
        Alors ajoute WheelGPT à ta chaîne dès maintenant !
      </>
    ),
    pluginName: "plugin OpenPlanet",
    openDashboard: "Ouvrir ton tableau de bord",
    getPlugin: "Télécharger le plugin",
    commandsTitle: "Commandes du chat",
    commandsLede: (format: ReactNode) => (
      <>
        Les temps s'écrivent {format}. Les heures et les minutes sont facultatives, les
        millisecondes non. Une prédiction parfaite rapporte 5 points, la plus proche 1 point.
      </>
    ),
    streamersTitle: "Streamers qui utilisent WheelGPT",
    channelCount: (c: number) =>
      `${n(c)} ${c === 1 ? "chaîne utilise" : "chaînes utilisent"} WheelGPT.`,
    noStreamers: "Personne n'est encore listé. Tu pourrais être le premier.",
    badges: {
      mod: "Modérateur",
      vip: "VIP",
      sub: "Abonné",
      streamer: "Diffuseur",
    },
    chatLabel: "Exemple de chat Twitch avec WheelGPT",
    chat: {
      needsMs: "(sans ms)",
      format: "hh:mm:ss.xxx (les heures et les minutes sont facultatives)",
      pb: "NOUVEAU RECORD PERSONNEL 47.881",
      thats: "C'est",
      nobody: "Personne n'a trouvé le temps exact, mais @tinytires a proposé 47.850 (-0.031)",
      point: "(+1 point)",
    },
    table: {
      command: "Commande",
      what: "Ce qu'elle fait",
      who: "Qui",
      also: "aussi",
    },
    access: {
      everyone: "tout le monde",
      subscriber: "abonnés",
      vip: "VIP",
      moderator: "modérateurs",
      broadcaster: "diffuseur",
    },
    commands: {
      wheelgpt: "C'est quoi ce bot ?",
      "wgpt-emotes": "Liste les emotes 7TV utilisées par WheelGPT. Ajoute-les à ta chaîne.",
      format: "Explique le format des temps pour les prédictions.",
      guess: "Devine le prochain record personnel du streamer sur la map en cours.",
      "wgpt-leaderboard":
        "Les meilleurs pronostiqueurs du mois ou de tous les temps, ou la position d'un viewer.",
      map: "Affiche la map en cours avec ses médailles et ses liens.",
      myguess: "Affiche ta prédiction actuelle.",
      resetguesses: "Supprime toutes les prédictions en cours.",
      result: "Termine la manche à la main avec le temps final.",
      room: "Affiche le salon serveur dans lequel joue le streamer.",
    },
  },

  channel: {
    notFoundTitle: "Chaîne introuvable",
    notFound: (name: ReactNode, link: (text: string) => ReactNode) => (
      <>
        {name} n'utilise pas WheelGPT. {link("Voir qui l'utilise.")}
      </>
    ),
    eyebrow: "Classement des prédictions",
    currentMap: "Map en cours",
    driving: "En piste en ce moment",
    by: (author: string) => `par ${author}`,
    room: "Salon",
    period: "Période",
    olderMonths: "Mois précédents",
    older: "Plus ancien…",
    findViewer: "Chercher un viewer",
    nobodyCalled: (q: string) => `Personne nommé « ${q} » n'a de points sur cette période.`,
    noPointsMonth:
      "Pas encore de points ce mois-ci. Le premier record personnel du mois lance la course.",
    noPoints: "Aucun point sur cette période.",
    summary: (c: number) =>
      `${n(c)} ${c === 1 ? "pronostiqueur" : "pronostiqueurs"} avec des points · prédiction parfaite 5 points, la plus proche 1 point`,
    viewer: "Viewer",
    points: "Points",
    perfect: "Parfaites",
    pages: "Pages",
    previous: "Précédent",
    next: "Suivant",
    pageOf: (page: number, pages: number) => `Page ${page} sur ${pages}`,
  },

  dashboard: {
    publicBoard: "Classement public",
    loggedOutTitle: "Ton tableau de bord",
    loggedOutLede:
      "Connecte-toi avec ton compte Twitch pour connecter le bot, obtenir ton token de plugin et régler le jeu de prédictions sur ta chaîne.",
    logout: "Se déconnecter",

    devTitle: "Mode dev",
    devIntro:
      "Ce serveur tourne en mode dev. Fais comme si tu étais en direct pour tester le plugin sans streamer.",
    devLive: "En direct (simulé)",
    devOffline: "Hors ligne",
    devLiveHelp:
      "La map, le salon et les records personnels envoyés par le plugin sont acceptés. Terminer le stream les efface comme une vraie fin de stream.",
    devOfflineHelp: "Hors ligne, les données du plugin sont refusées.",
    endStream: "Terminer le stream",
    goLive: "Passer en direct",

    botTitle: "Bot",
    botIntro: "WheelGPT répond dans ton chat Twitch tant qu'il est connecté.",
    chatStatus: {
      enabled: "Lit ton chat",
      unknown: "Twitch n'a pas répondu",
      needs_authorization: "Twitch a refusé l'accès à ton chat",
      disabled: "Déconnecté",
      off: "Connecté, mais n'écoute pas le chat",
    },
    eventsubOff: (env: ReactNode) => (
      <>
        Ce serveur ne s'abonne pas aux chats Twitch ({env} est désactivé), donc les commandes du
        chat ne lui parviennent pas. C'est normal en développement local : Twitch ne peut envoyer le
        chat qu'à une adresse https publique.
      </>
    ),
    tryAgainLater: "Réessaie de te connecter dans un instant.",
    oneMoreStep: "Encore une étape",
    needsAuth: (cmd: ReactNode, link: (text: string) => ReactNode) => (
      <>
        Reconnecte-toi pour que Twitch laisse le bot entrer dans ton chat, ou tape {cmd} dans ton
        chat. {link("Se reconnecter")}
      </>
    ),
    legacyAuth: (link: (text: string) => ReactNode) => (
      <>
        Tu t'es inscrit avant la nouvelle connexion. Le bot continue de fonctionner tant qu'il est
        modérateur de ta chaîne ; {link("reconnecte-toi")} pour passer définitivement au nouveau
        système.
      </>
    ),
    working: "En cours…",
    disconnect: "Déconnecter le bot",
    connect: "Connecter le bot",

    tokenTitle: "Token du plugin",
    tokenIntro: (path: ReactNode) => (
      <>
        Le plugin Openplanet l'utilise pour envoyer ta map, ton salon et tes records personnels.
        Colle-le dans {path} et garde-le pour toi.
      </>
    ),
    tokenPath: "Settings → WheelGPT → Token",
    newToken: "Ton nouveau token (affiché une seule fois)",
    copyNow: "Copie-le maintenant. Nous n'en gardons qu'une empreinte.",
    tokenCreated: (created: string, lastUsed?: string) =>
      `Créé le ${created}${lastUsed ? `, utilisé pour la dernière fois le ${lastUsed}` : ", pas encore utilisé"}.`,
    legacyToken: "Ton token de l'ancien site fonctionne encore. Un nouveau token le remplace.",
    noToken: "Tu n'as pas encore de token.",
    confirmNewToken: "Le token actuel cessera de fonctionner. En créer un nouveau ?",
    createNewToken: "Créer un nouveau token",
    createToken: "Créer un token",

    settingsTitle: "Prédictions",
    settingsIntro: "Comment le jeu de prédictions fonctionne sur ta chaîne.",
    delay: "Délai d'annonce",
    seconds: "secondes",
    delayHelp:
      "Attendre ce délai avant d'annoncer un record, pour que le chat ne le voie pas avant ton stream. Règle-le sur le délai de ton stream.",
    minAge: "Prédictions de dernière seconde",
    minAgeHelp:
      "Ignorer les prédictions faites ou modifiées moins de ce délai avant le record. Empêche de deviner en te regardant franchir la ligne d'arrivée.",
    listPublic: "Lister ma chaîne sur wheelgpt.dev",
    listPublicHelp: "Ton classement reste accessible par lien dans tous les cas.",
    saving: "Enregistrement…",
    save: "Enregistrer les réglages",
    saved: "Enregistré.",

    deleteTitle: "Supprimer WheelGPT",
    deleteIntro:
      "Retire le bot de ta chaîne et supprime tes réglages, ton token et tes classements. Cette action est irréversible.",
    typeToConfirm: (login: ReactNode) => <>Tape {login} pour confirmer.</>,
    channelName: "Nom de la chaîne",
    deleteAll: "Tout supprimer",
  },

  legal: {
    imprint: () => (
      <>
        <h1>Mentions légales</h1>
        <p>Ce site est exploité par Markus Wieland.</p>
        <p>
          Contact : <Mail />
        </p>
        <p>
          Tout le contenu de ce site est fourni sans aucune garantie. L'auteur décline toute
          responsabilité quant à l'exactitude, l'exhaustivité et l'actualité du contenu.
          L'utilisation du site se fait à tes propres risques.
        </p>
      </>
    ),
    privacy: (dashboard: (text: string) => ReactNode) => (
      <>
        <h1>Confidentialité</h1>

        <h2>Comptes streamer</h2>
        <p>
          Quand tu te connectes sur wheelgpt.dev, nous enregistrons l'identifiant de ta chaîne
          Twitch, ton login, ton nom d'affichage et l'URL de ta photo de profil, ainsi que tes
          réglages. Nous demandons à Twitch une seule autorisation, <code>channel:bot</code>, pour
          que le compte bot WheelGPT puisse lire et écrire dans le chat de ta chaîne. Le jeton
          d'accès que Twitch nous donne lors de la connexion sert une seule fois à savoir qui tu es,
          puis il est supprimé ; il n'est jamais enregistré.
        </p>
        <p>
          La connexion crée un cookie, <code>wgpt_session</code>, qui te garde connecté jusqu'à 30
          jours. Il est strictement nécessaire et n'est pas utilisé pour le suivi. Ton token de
          plugin est enregistré uniquement sous forme d'empreinte (hash). Ton choix de langue est
          mémorisé dans le stockage local de ton navigateur et ne nous est jamais envoyé.
        </p>

        <h2>Viewers</h2>
        <p>
          Quand tu utilises une commande comme <code>!guess</code>, le bot enregistre ton
          identifiant Twitch, ton nom d'affichage et ta prédiction. Après un record personnel, ta
          prédiction, son résultat et tes points sont conservés pour les classements de la chaîne et
          ton historique de prédictions, qui sont publics sur ce site. Les messages qui ne sont pas
          des commandes ne sont pas enregistrés.
        </p>

        <h2>Journal d'accès</h2>
        <p>
          Nous conservons des journaux de trafic pouvant contenir ton adresse IP et l'agent
          utilisateur de ton navigateur, pour voir comment le site est utilisé et pour empêcher les
          abus. Ils font l'objet d'une rotation régulière et sont généralement conservés 24 heures
          au maximum.
        </p>

        <h2>Supprimer tes données</h2>
        <p>
          Les streamers peuvent supprimer leur chaîne et toutes ses données depuis le{" "}
          {dashboard("tableau de bord")}. Les viewers qui souhaitent faire supprimer leurs
          prédictions peuvent écrire à <Mail /> ou contacter @sowiemarkus sur Discord.
        </p>
      </>
    ),
    notFoundTitle: "Sortie de piste",
    notFound: (link: (text: string) => ReactNode) => (
      <>Il n'y a rien à cette adresse. {link("Retour au départ")}</>
    ),
  },
};
