import { useTranslation } from "react-i18next";

type PlaceholderPageProps = {
  title: string;
  description: string;
};

export function PlaceholderPage({ title, description }: PlaceholderPageProps) {
  const { t } = useTranslation();

  return (
    <section className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
      <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{title}</div>
      <h3 className="mt-3 font-serif text-4xl font-bold">{title}</h3>
      <p className="mt-3 max-w-3xl text-base leading-7 text-slate-400">{t(description)}</p>
      <div className="mt-8 rounded-2xl border border-dashed border-white/10 bg-white/[0.03] p-5 text-sm text-stone-400">
        {t("pages.inProgress")}
      </div>
    </section>
  );
}
