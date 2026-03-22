import type { ReactNode } from "react";
import { LanguageSwitcher } from "../components/language-switcher";
import { useTranslation } from "react-i18next";

type ConsoleShellProps = {
  sidebar: ReactNode;
  children: ReactNode;
};

export function ConsoleShell({ sidebar, children }: ConsoleShellProps) {
  const { t } = useTranslation();

  return (
    <main className="min-h-[100svh] w-full px-2 py-3 sm:px-3 sm:py-4 lg:px-4">
      <header className="mb-4 flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <h1 className="font-serif text-4xl font-bold tracking-tight text-stone-900 sm:text-5xl">
            {t("common.appName")}
          </h1>
          <p className="mt-2 text-base text-stone-600 sm:text-lg">{t("common.appTagline")}</p>
        </div>

        <div className="flex flex-col items-start gap-3 lg:items-end">
          <div className="text-xs uppercase tracking-[0.42em] text-stone-500">
            {t("dashboard.caption")}
          </div>
          <LanguageSwitcher />
        </div>
      </header>

      <section className="grid min-h-[calc(100svh-7.5rem)] items-start gap-4 rounded-[2rem] bg-[linear-gradient(180deg,#171c20_0%,#111519_100%)] p-3 shadow-[0_24px_48px_rgba(57,44,30,0.18)] lg:min-h-[calc(100svh-8.25rem)] lg:grid-cols-[272px_minmax(0,1fr)] lg:p-4">
        {sidebar}
        <section className="min-h-0 rounded-[1.75rem] bg-white/[0.02] p-4 lg:p-5">{children}</section>
      </section>
    </main>
  );
}
