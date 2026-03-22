import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { getEnv, updateEnv } from "../lib/api/env";
import type { AppShellContext } from "../types/app-shell";

type Row = {
  id: string;
  key: string;
  value: string;
};

export function EnvVarsPage() {
  const { t } = useTranslation();
  const { activeEnv } = useOutletContext<AppShellContext>();
  const [rows, setRows] = useState<Row[]>([]);
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
          setRows(fromEnvVars(data.envVars));
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("envVars.loadError"));
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
      await updateEnv(activeEnv.name, { envVars: toEnvVars(rows) });
      setMessage(t("envVars.saveSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("envVars.saveError"));
    } finally {
      setSaving(false);
    }
  }

  function updateRow(id: string, patch: Partial<Row>) {
    setRows((current) => current.map((row) => (row.id === id ? { ...row, ...patch } : row)));
  }

  function removeRow(id: string) {
    setRows((current) => current.filter((row) => row.id !== id));
  }

  function addRow() {
    setRows((current) => [...current, createRow("", "")]);
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(320px,0.9fr)]">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("common.envVars")}</div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("envVars.title")}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t("envVars.description")}</p>

        {loading ? (
          <div className="mt-8 rounded-2xl border border-white/8 bg-white/[0.03] p-5 text-sm text-stone-400">
            {t("envVars.loading")}
          </div>
        ) : (
          <>
            <div className="mt-8 grid gap-4">
              {rows.length === 0 ? (
                <div className="rounded-2xl border border-dashed border-white/10 bg-white/[0.03] p-5 text-sm text-stone-400">
                  {t("envVars.empty")}
                </div>
              ) : (
                rows.map((row) => (
                  <div key={row.id} className="grid gap-3 rounded-3xl bg-white/[0.03] p-4 lg:grid-cols-[minmax(180px,0.9fr)_minmax(0,1.2fr)_auto]">
                    <input
                      className="rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                      placeholder={t("envVars.keyPlaceholder")}
                      value={row.key}
                      onChange={(event) => updateRow(row.id, { key: event.target.value })}
                    />
                    <input
                      className="rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                      placeholder={t("envVars.valuePlaceholder")}
                      value={row.value}
                      onChange={(event) => updateRow(row.id, { value: event.target.value })}
                    />
                    <button
                      className="rounded-2xl bg-ccv-danger px-4 py-3 text-sm text-[#ffe5d8]"
                      type="button"
                      onClick={() => removeRow(row.id)}
                    >
                      {t("envVars.remove")}
                    </button>
                  </div>
                ))
              )}
            </div>

            {message ? (
              <div className="mt-6 rounded-2xl bg-[#7e8c67]/20 px-4 py-3 text-sm text-[#e3edd6]">{message}</div>
            ) : null}
            {error ? (
              <div className="mt-6 rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">{error}</div>
            ) : null}

            <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:justify-between">
              <button
                className="rounded-2xl bg-ccv-panel-strong px-5 py-3 text-sm text-ccv-ink"
                type="button"
                onClick={addRow}
              >
                {t("envVars.add")}
              </button>
              <button
                className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
                disabled={saving}
                type="button"
                onClick={() => void handleSave()}
              >
                {saving ? t("envVars.saving") : t("envVars.save")}
              </button>
            </div>
          </>
        )}
      </article>

      <aside className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("envVars.guidance")}</div>
        <div className="mt-4 rounded-2xl bg-white/[0.03] p-5 text-sm leading-7 text-stone-300">
          <p>{t("envVars.guidanceIntro")}</p>
          <p className="mt-3">{t("envVars.guidanceSecrets")}</p>
          <p className="mt-3">{t("envVars.guidanceScope")}</p>
        </div>
      </aside>
    </section>
  );
}

function fromEnvVars(envVars?: Record<string, string>) {
  return Object.entries(envVars ?? {})
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([key, value]) => createRow(key, value));
}

function toEnvVars(rows: Row[]) {
  const envVars: Record<string, string> = {};

  for (const row of rows) {
    const key = row.key.trim();
    if (!key) {
      continue;
    }
    envVars[key] = row.value;
  }

  return envVars;
}

function createRow(key: string, value: string): Row {
  return {
    id: `${key}-${value}-${Math.random().toString(36).slice(2, 10)}`,
    key,
    value,
  };
}
