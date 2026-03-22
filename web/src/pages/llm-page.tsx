import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { getLLMConfig, updateLLMConfig } from "../lib/api/llm";
import type { AppShellContext } from "../types/app-shell";
import type { LLMConfig } from "../types/llm";

const emptyConfig: LLMConfig = {
  apiKey: "",
  baseUrl: "",
  models: {
    default: "",
    sonnet: "",
    opus: "",
    haiku: "",
  },
};

export function LLMPage() {
  const { t } = useTranslation();
  const { activeEnv } = useOutletContext<AppShellContext>();
  const [config, setConfig] = useState<LLMConfig>(emptyConfig);
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

    void getLLMConfig(activeEnv.name)
      .then((llmConfig) => {
        if (!cancelled) {
          setConfig(llmConfig);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("llm.loadError"));
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

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!activeEnv) {
      return;
    }

    setSaving(true);
    setError(null);
    setMessage(null);

    try {
      await updateLLMConfig(activeEnv.name, config);
      setMessage(t("llm.saveSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("llm.saveError"));
    } finally {
      setSaving(false);
    }
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[minmax(0,1.5fr)_minmax(300px,0.9fr)]">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("common.llm")}</div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("llm.title")}</h3>
        <p className="mt-3 max-w-3xl text-base leading-7 text-slate-400">{t("llm.description")}</p>

        {loading ? (
          <div className="mt-8 rounded-2xl border border-white/8 bg-white/[0.03] p-5 text-sm text-stone-400">
            {t("llm.loading")}
          </div>
        ) : (
          <form className="mt-8 grid gap-5" onSubmit={handleSubmit}>
            <Field label={t("llm.baseUrl")} value={config.baseUrl}>
              <input
                className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                placeholder="https://api.example.com"
                type="url"
                value={config.baseUrl}
                onChange={(event) => setConfig({ ...config, baseUrl: event.target.value })}
              />
            </Field>

            <Field label={t("llm.apiKey")} value={config.apiKey}>
              <input
                className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                placeholder="sk-..."
                type="password"
                value={config.apiKey}
                onChange={(event) => setConfig({ ...config, apiKey: event.target.value })}
              />
            </Field>

            <div className="grid gap-4 md:grid-cols-2">
              <Field label={t("llm.defaultModel")} value={config.models.default}>
                <input
                  className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="model-name"
                  value={config.models.default}
                  onChange={(event) =>
                    setConfig({
                      ...config,
                      models: { ...config.models, default: event.target.value },
                    })
                  }
                />
              </Field>
              <Field label={t("llm.sonnetModel")} value={config.models.sonnet}>
                <input
                  className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="model-name"
                  value={config.models.sonnet}
                  onChange={(event) =>
                    setConfig({
                      ...config,
                      models: { ...config.models, sonnet: event.target.value },
                    })
                  }
                />
              </Field>
              <Field label={t("llm.opusModel")} value={config.models.opus}>
                <input
                  className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="model-name"
                  value={config.models.opus}
                  onChange={(event) =>
                    setConfig({
                      ...config,
                      models: { ...config.models, opus: event.target.value },
                    })
                  }
                />
              </Field>
              <Field label={t("llm.haikuModel")} value={config.models.haiku}>
                <input
                  className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="model-name"
                  value={config.models.haiku}
                  onChange={(event) =>
                    setConfig({
                      ...config,
                      models: { ...config.models, haiku: event.target.value },
                    })
                  }
                />
              </Field>
            </div>

            {message ? (
              <div className="rounded-2xl bg-[#7e8c67]/20 px-4 py-3 text-sm text-[#e3edd6]">{message}</div>
            ) : null}
            {error ? (
              <div className="rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">{error}</div>
            ) : null}

            <div className="flex justify-end">
              <button
                className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
                disabled={saving}
                type="submit"
              >
                {saving ? t("llm.saving") : t("llm.save")}
              </button>
            </div>
          </form>
        )}
      </article>

      <aside className="flex flex-col gap-4 rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("llm.guidance")}</div>
        <div className="rounded-2xl bg-white/[0.03] p-5 text-sm leading-7 text-stone-300">
          <p>{t("llm.guidanceIntro")}</p>
          <p className="mt-3">{t("llm.guidanceMasked")}</p>
          <p className="mt-3">{t("llm.guidanceEnvScope")}</p>
        </div>
      </aside>
    </section>
  );
}

type FieldProps = {
  label: string;
  value?: string;
  children: React.ReactNode;
};

function Field({ label, value, children }: FieldProps) {
  return (
    <label className="grid gap-2">
      <span className="text-xs uppercase tracking-[0.28em] text-stone-400">{label}</span>
      {children}
      {value ? <span className="text-xs text-slate-500">{value}</span> : null}
    </label>
  );
}
