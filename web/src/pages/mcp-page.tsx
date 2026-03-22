import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { addMCP, deleteMCP, listMCP } from "../lib/api/mcp";
import type { AppShellContext } from "../types/app-shell";
import type { MCPServer } from "../types/mcp";

type MCPFormState = {
  name: string;
  type: "stdio" | "http";
  url: string;
  headersText: string;
  command: string;
  argsText: string;
  envText: string;
};

const emptyForm: MCPFormState = {
  name: "",
  type: "stdio",
  url: "",
  headersText: "",
  command: "",
  argsText: "",
  envText: "",
};

export function MCPPage() {
  const { t } = useTranslation();
  const { activeEnv } = useOutletContext<AppShellContext>();
  const [servers, setServers] = useState<Record<string, MCPServer>>({});
  const [form, setForm] = useState<MCPFormState>(emptyForm);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [deletingName, setDeletingName] = useState<string | null>(null);
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

    void listMCP(activeEnv.name)
      .then((data) => {
        if (!cancelled) {
          setServers(data.servers);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("mcp.loadError"));
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
      const headers = parseEnvLines(form.headersText);
      const env = parseEnvLines(form.envText);
      const payload = {
        name: form.name.trim(),
        config:
          form.type === "http"
            ? {
                type: "http",
                url: form.url.trim(),
                headers,
              }
            : {
                type: "stdio",
                command: form.command.trim(),
                args: splitLines(form.argsText),
                env,
              },
      };

      await addMCP(activeEnv.name, payload);
      setServers((current) => ({ ...current, [payload.name]: payload.config }));
      setForm(emptyForm);
      setMessage(t("mcp.addSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("mcp.addError"));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(serverName: string) {
    if (!activeEnv) {
      return;
    }

    setDeletingName(serverName);
    setError(null);
    setMessage(null);

    try {
      await deleteMCP(activeEnv.name, serverName);
      setServers((current) => {
        const next = { ...current };
        delete next[serverName];
        return next;
      });
      setMessage(t("mcp.deleteSuccess", { name: serverName }));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("mcp.deleteError"));
    } finally {
      setDeletingName(null);
    }
  }

  const serverEntries = Object.entries(servers);

  return (
    <section className="grid gap-5 xl:grid-cols-[minmax(320px,0.9fr)_minmax(0,1.3fr)]">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("common.mcp")}</div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("mcp.title")}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t("mcp.description")}</p>

        <form className="mt-8 grid gap-5" onSubmit={handleSubmit}>
          <FormField label={t("mcp.serverName")}>
            <input
              className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
              placeholder="filesystem"
              value={form.name}
              onChange={(event) => setForm({ ...form, name: event.target.value })}
            />
          </FormField>

          <FormField label={t("mcp.type")}>
            <select
              className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none"
              value={form.type}
              onChange={(event) => setForm({ ...form, type: event.target.value as MCPFormState["type"] })}
            >
              <option value="stdio">stdio</option>
              <option value="http">http</option>
            </select>
          </FormField>

          {form.type === "http" ? (
            <>
              <FormField label={t("mcp.url")}>
                <input
                  className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="https://mcp.example.com/mcp"
                  value={form.url}
                  onChange={(event) => setForm({ ...form, url: event.target.value })}
                />
              </FormField>

              <FormField hint={t("mcp.headersHint")} label={t("mcp.headers")}>
                <textarea
                  className="min-h-28 w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder={"AUTH_TOKEN=xxx\nX_API_KEY=yyy"}
                  value={form.headersText}
                  onChange={(event) => setForm({ ...form, headersText: event.target.value })}
                />
              </FormField>
            </>
          ) : (
            <>
              <FormField label={t("mcp.command")}>
                <input
                  className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="npx"
                  value={form.command}
                  onChange={(event) => setForm({ ...form, command: event.target.value })}
                />
              </FormField>

              <FormField hint={t("mcp.argsHint")} label={t("mcp.args")}>
                <textarea
                  className="min-h-28 w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder="-y&#10;@anthropic-ai/mcp-server-filesystem&#10;/path"
                  value={form.argsText}
                  onChange={(event) => setForm({ ...form, argsText: event.target.value })}
                />
              </FormField>

              <FormField hint={t("mcp.envHint")} label={t("mcp.env")}>
                <textarea
                  className="min-h-28 w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
                  placeholder={"API_KEY=xxx\nROOT=/workspace"}
                  value={form.envText}
                  onChange={(event) => setForm({ ...form, envText: event.target.value })}
                />
              </FormField>
            </>
          )}

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
              {saving ? t("mcp.adding") : t("mcp.add")}
            </button>
          </div>
        </form>
      </article>

      <section className="grid gap-5">
        <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
          <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("mcp.registeredServers")}</div>
          <h3 className="mt-3 font-serif text-4xl font-bold">{t("mcp.serverListTitle")}</h3>
          <p className="mt-3 text-base leading-7 text-slate-400">{t("mcp.serverListDescription")}</p>

          {loading ? (
            <div className="mt-8 rounded-2xl border border-white/8 bg-white/[0.03] p-5 text-sm text-stone-400">
              {t("mcp.loading")}
            </div>
          ) : serverEntries.length === 0 ? (
            <div className="mt-8 rounded-2xl border border-dashed border-white/10 bg-white/[0.03] p-5 text-sm text-stone-400">
              {t("mcp.empty")}
            </div>
          ) : (
            <div className="mt-8 grid gap-4">
              {serverEntries.map(([name, server]) => (
                <article key={name} className="rounded-2xl bg-white/[0.03] px-4 py-4">
                  <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
                    <div className="min-w-0">
                      <div className="text-xs uppercase tracking-[0.28em] text-stone-500">{t("mcp.server")}</div>
                      <h4 className="mt-1.5 text-xl font-semibold text-ccv-ink">{name}</h4>
                      <p className="mt-2 text-sm text-slate-400">
                        {t("mcp.typeLabel")}:{" "}
                        <span className="text-stone-200">{server.type || inferServerType(server)}</span>
                      </p>
                      <p className="mt-1.5 break-all text-sm text-slate-400">
                        {t("mcp.commandLabel")}:{" "}
                        <span className="text-stone-200">
                          {server.command || server.url || t("mcp.none")}
                        </span>
                      </p>
                      <p className="mt-1.5 text-sm text-slate-400">
                        {t("mcp.argsLabel")}:{" "}
                        <span className="text-stone-200">
                          {(server.args ?? []).length > 0 ? server.args?.join(" ") : t("mcp.none")}
                        </span>
                      </p>
                    </div>

                    <button
                      className="rounded-2xl bg-ccv-danger px-4 py-2.5 text-sm text-[#ffe5d8] disabled:cursor-not-allowed disabled:opacity-60"
                      disabled={deletingName === name}
                      type="button"
                      onClick={() => void handleDelete(name)}
                    >
                      {deletingName === name ? t("mcp.deleting") : t("mcp.delete")}
                    </button>
                  </div>
                </article>
              ))}
            </div>
          )}
        </article>

        <aside className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
          <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("mcp.guidance")}</div>
          <div className="mt-4 rounded-2xl bg-white/[0.03] p-5 text-sm leading-7 text-stone-300">
            <p>{t("mcp.guidanceIntro")}</p>
            <p className="mt-3">{t("mcp.guidanceArgs")}</p>
            <p className="mt-3">{t("mcp.guidanceEnv")}</p>
          </div>
        </aside>
      </section>
    </section>
  );
}

type FormFieldProps = {
  label: string;
  hint?: string;
  children: React.ReactNode;
};

function FormField({ label, hint, children }: FormFieldProps) {
  return (
    <label className="grid gap-2">
      <span className="text-xs uppercase tracking-[0.28em] text-stone-400">{label}</span>
      {children}
      {hint ? <span className="text-xs text-slate-500">{hint}</span> : null}
    </label>
  );
}

function splitLines(value: string) {
  return value
    .split("\n")
    .map((line) => line.trim())
    .filter(Boolean);
}

function parseEnvLines(value: string) {
  const env: Record<string, string> = {};

  for (const line of splitLines(value)) {
    const index = line.indexOf("=");
    if (index <= 0) {
      continue;
    }

    const key = line.slice(0, index).trim();
    const val = line.slice(index + 1).trim();
    if (key) {
      env[key] = val;
    }
  }

  return env;
}

function inferServerType(server: MCPServer) {
  if (server.type) {
    return server.type;
  }
  if (server.url) {
    return "http";
  }
  if (server.command) {
    return "stdio";
  }
  return "unknown";
}
