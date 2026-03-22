import { useTranslation } from "react-i18next";

const languages = [
  { code: "zh-CN", labelKey: "common.chinese" },
  { code: "en", labelKey: "common.english" },
] as const;

export function LanguageSwitcher() {
  const { i18n, t } = useTranslation();

  return (
    <div className="flex items-center gap-1.5 rounded-2xl border border-white/12 bg-white/6 p-1.5 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]">
      <span className="px-3 text-xs font-medium uppercase tracking-[0.3em] text-stone-300">
        {t("common.language")}
      </span>
      {languages.map((language) => {
        const active = i18n.language === language.code;
        return (
          <button
            key={language.code}
            className={
              active
                ? "rounded-xl bg-stone-100 px-3.5 py-2 text-sm font-semibold text-slate-900 shadow-sm"
                : "rounded-xl border border-transparent bg-white/[0.03] px-3.5 py-2 text-sm font-medium text-stone-100 transition hover:border-white/10 hover:bg-white/10"
            }
            type="button"
            onClick={() => void i18n.changeLanguage(language.code)}
          >
            {t(language.labelKey)}
          </button>
        );
      })}
    </div>
  );
}
