import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";
import { de } from "./de";
import { en, type Messages } from "./en";
import { fr } from "./fr";

export const languages = ["en", "de", "fr"] as const;
export type Lang = (typeof languages)[number];

const catalogs: Record<Lang, Messages> = { en, de, fr };
const STORAGE_KEY = "wgpt_lang";

function isLang(value: string | null | undefined): value is Lang {
  return (languages as readonly string[]).includes(value ?? "");
}

/** Saved choice first, then the browser's preferred languages, then English. */
export function detectLang(): Lang {
  try {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (isLang(saved)) return saved;
  } catch {
    // Storage may be blocked; fall through to the browser languages.
  }
  const preferred = typeof navigator === "undefined" ? [] : (navigator.languages ?? []);
  for (const l of preferred) {
    const base = l.toLowerCase().split("-")[0];
    if (isLang(base)) return base;
  }
  return "en";
}

interface I18n {
  lang: Lang;
  t: Messages;
  setLang: (lang: Lang) => void;
}

const I18nContext = createContext<I18n>({ lang: "en", t: en, setLang: () => {} });

export function I18nProvider({ children, initial }: { children: ReactNode; initial?: Lang }) {
  const [lang, setLangState] = useState<Lang>(() => initial ?? detectLang());

  useEffect(() => {
    document.documentElement.lang = lang;
    document.title = catalogs[lang].meta.title;
    document
      .querySelector('meta[name="description"]')
      ?.setAttribute("content", catalogs[lang].meta.description);
  }, [lang]);

  const setLang = useCallback((next: Lang) => {
    setLangState(next);
    try {
      localStorage.setItem(STORAGE_KEY, next);
    } catch {
      // Not remembered, but still switched for this visit.
    }
  }, []);

  return (
    <I18nContext.Provider value={{ lang, t: catalogs[lang], setLang }}>
      {children}
    </I18nContext.Provider>
  );
}

export function useI18n(): I18n {
  return useContext(I18nContext);
}

export type { Messages };
