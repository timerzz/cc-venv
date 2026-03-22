import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { getEnv, updateEnv } from "../lib/api/env";
import type { AppShellContext } from "../types/app-shell";

export function NotesPage() {
  const { t } = useTranslation();
  const { activeEnv } = useOutletContext<AppShellContext>();
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!activeEnv) {
      return;
    }

    let cancelled = false;
    setLoading(true);
    setError(null);
    setMessage(null);

    void getEnv(activeEnv.name)
      .then((data) => {
        if (!cancelled) {
          setContent(data.claudeMd ?? "");
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("notes.loadError"));
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [activeEnv, t]);

  async function handleSave() {
    if (!activeEnv) {
      return;
    }

    setSaving(true);
    setError(null);
    setMessage(null);

    try {
      await updateEnv(activeEnv.name, { claudeMd: content });
      setMessage(t("notes.saveSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("notes.saveError"));
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(320px,0.9fr)]">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("common.notes")}</div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("notes.title")}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t("notes.description")}</p>

        {loading ? (
          <div className="mt-8 rounded-2xl border border-white/8 bg-white/[0.03] p-5 text-sm text-stone-400">
            {t("notes.loading")}
          </div>
        ) : (
          <>
            <textarea
              className="mt-8 min-h-[420px] w-full rounded-3xl border border-white/8 bg-[#151b1f] px-5 py-4 font-mono text-sm leading-7 text-ccv-ink outline-none placeholder:text-slate-500"
              placeholder={t("notes.placeholder")}
              value={content}
              onChange={(event) => setContent(event.target.value)}
            />

            {message ? (
              <div className="mt-6 rounded-2xl bg-[#7e8c67]/20 px-4 py-3 text-sm text-[#e3edd6]">{message}</div>
            ) : null}
            {error ? (
              <div className="mt-6 rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">{error}</div>
            ) : null}

            <div className="mt-6 flex justify-end">
              <button
                className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
                disabled={saving}
                type="button"
                onClick={() => void handleSave()}
              >
                {saving ? t("notes.saving") : t("notes.save")}
              </button>
            </div>
          </>
        )}
      </article>

      <aside className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("notes.guidance")}</div>
        <div className="mt-4 rounded-2xl bg-white/[0.03] p-5 text-sm leading-7 text-stone-300">
          <p>{t("notes.guidanceIntro")}</p>
          <p className="mt-3">{t("notes.guidanceExamples")}</p>
          <p className="mt-3">{t("notes.guidanceScope")}</p>
        </div>
      </aside>
    </section>
  );
}
