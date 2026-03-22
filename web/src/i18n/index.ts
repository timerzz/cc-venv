import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { en } from "./locales/en";
import { zhCN } from "./locales/zh-cn";

export const supportedLanguages = ["zh-CN", "en"] as const;
export type SupportedLanguage = (typeof supportedLanguages)[number];

const LANGUAGE_STORAGE_KEY = "ccv.web.language";

function detectLanguage(): SupportedLanguage {
  if (typeof window === "undefined") {
    return "zh-CN";
  }

  const saved = window.localStorage.getItem(LANGUAGE_STORAGE_KEY);
  if (saved === "zh-CN" || saved === "en") {
    return saved;
  }

  const browserLanguage = window.navigator.language.toLowerCase();
  return browserLanguage.startsWith("zh") ? "zh-CN" : "en";
}

void i18n.use(initReactI18next).init({
  lng: detectLanguage(),
  fallbackLng: "en",
  resources: {
    en: { translation: en },
    "zh-CN": { translation: zhCN },
  },
  interpolation: {
    escapeValue: false,
  },
});

i18n.on("languageChanged", (language) => {
  if (typeof window !== "undefined") {
    window.localStorage.setItem(LANGUAGE_STORAGE_KEY, language);
  }
});

export { i18n, LANGUAGE_STORAGE_KEY };
