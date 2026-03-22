import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import {
  deleteResourceFile,
  getResourceFile,
  listResourceFiles,
  upsertResourceFile,
} from "../lib/api/resource-files";
import type { AppShellContext } from "../types/app-shell";
import type { ResourceKind } from "../types/resource-file";

type ResourceFilesPageProps = {
  kind: ResourceKind;
  titleKey: string;
  descriptionKey: string;
  emptyKey: string;
  guidanceIntroKey: string;
  nameLabelKey: string;
  namePlaceholder: string;
};

export function ResourceFilesPage({
  kind,
  titleKey,
  descriptionKey,
  emptyKey,
  guidanceIntroKey,
  nameLabelKey,
  namePlaceholder,
}: ResourceFilesPageProps) {
  const { t } = useTranslation();
  const { activeEnv } = useOutletContext<AppShellContext>();
  const [items, setItems] = useState<string[]>([]);
  const [selectedName, setSelectedName] = useState("");
  const [draftName, setDraftName] = useState("");
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!activeEnv) {
      setItems([]);
      setSelectedName("");
      setDraftName("");
      setContent("");
      setLoading(false);
      return;
    }

    let cancelled = false;
    setLoading(true);
    setError(null);
    setMessage(null);

    void listResourceFiles(activeEnv.name, kind)
      .then((data) => {
        if (cancelled) {
          return;
        }
        setItems(data.items);
        const nextName = data.items[0] ?? "";
        setSelectedName(nextName);
        setDraftName(nextName);
        if (!nextName) {
          setContent("");
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("resources.loadError"));
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
  }, [activeEnv, kind, t]);

  useEffect(() => {
    if (!activeEnv || !selectedName) {
      return;
    }

    let cancelled = false;
    setError(null);
    setMessage(null);

    void getResourceFile(activeEnv.name, kind, selectedName)
      .then((file) => {
        if (!cancelled) {
          setDraftName(file.name);
          setContent(file.content);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("resources.loadError"));
        }
      });

    return () => {
      cancelled = true;
    };
  }, [activeEnv, kind, selectedName, t]);

  const hasSelection = useMemo(() => selectedName !== "", [selectedName]);

  async function handleSave() {
    if (!activeEnv || !draftName.trim()) {
      setError(t("resources.nameRequired"));
      return;
    }

    setSaving(true);
    setError(null);
    setMessage(null);
    try {
      const file = await upsertResourceFile(activeEnv.name, kind, {
        name: draftName.trim(),
        content,
      });

      setItems((current) => Array.from(new Set([...current.filter((item) => item !== selectedName), file.name])).sort());
      setSelectedName(file.name);
      setDraftName(file.name);
      setMessage(t("resources.saveSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("resources.saveError"));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete() {
    if (!activeEnv || !selectedName) {
      return;
    }

    setDeleting(true);
    setError(null);
    setMessage(null);
    try {
      await deleteResourceFile(activeEnv.name, kind, selectedName);
      const nextItems = items.filter((item) => item !== selectedName);
      setItems(nextItems);
      setSelectedName(nextItems[0] ?? "");
      setDraftName(nextItems[0] ?? "");
      setContent("");
      setMessage(t("resources.deleteSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("resources.deleteError"));
    } finally {
      setDeleting(false);
    }
  }

  function handleCreateNew() {
    setSelectedName("");
    setDraftName("");
    setContent("");
    setError(null);
    setMessage(null);
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[280px_minmax(0,1fr)]">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-6 text-ccv-ink">
        <div className="flex items-center justify-between gap-3">
          <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("resources.library")}</div>
          <button className="rounded-xl bg-ccv-panel-strong px-3 py-2 text-xs text-ccv-ink" type="button" onClick={handleCreateNew}>
            {t("resources.new")}
          </button>
        </div>

        {loading ? (
          <div className="mt-6 rounded-2xl border border-white/8 bg-white/[0.03] p-4 text-sm text-stone-400">
            {t("resources.loading")}
          </div>
        ) : items.length === 0 ? (
          <div className="mt-6 rounded-2xl border border-dashed border-white/10 bg-white/[0.03] p-4 text-sm text-stone-400">
            {t(emptyKey)}
          </div>
        ) : (
          <div className="mt-6 grid gap-2">
            {items.map((item) => (
              <button
                key={item}
                className={
                  item === selectedName
                    ? "rounded-2xl bg-stone-100 px-4 py-3 text-left text-sm font-semibold text-slate-900"
                    : "rounded-2xl bg-white/[0.03] px-4 py-3 text-left text-sm text-stone-300"
                }
                type="button"
                onClick={() => setSelectedName(item)}
              >
                {item}
              </button>
            ))}
          </div>
        )}
      </article>

      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("resources.editor")}</div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t(titleKey)}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t(descriptionKey)}</p>

        <div className="mt-8 grid gap-5">
          <label className="grid gap-2">
            <span className="text-xs uppercase tracking-[0.28em] text-stone-400">{t(nameLabelKey)}</span>
            <input
              className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
              placeholder={namePlaceholder}
              value={draftName}
              onChange={(event) => setDraftName(event.target.value)}
            />
          </label>

          <label className="grid gap-2">
            <span className="text-xs uppercase tracking-[0.28em] text-stone-400">{t("resources.content")}</span>
            <textarea
              className="min-h-[360px] w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
              placeholder={t("resources.contentPlaceholder")}
              value={content}
              onChange={(event) => setContent(event.target.value)}
            />
          </label>

          {message ? (
            <div className="rounded-2xl bg-[#7e8c67]/20 px-4 py-3 text-sm text-[#e3edd6]">{message}</div>
          ) : null}
          {error ? (
            <div className="rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">{error}</div>
          ) : null}

          <div className="flex flex-wrap justify-end gap-3">
            <button
              className="rounded-2xl bg-ccv-danger px-5 py-3 text-sm text-[#ffe5d8] disabled:cursor-not-allowed disabled:opacity-60"
              disabled={deleting || !hasSelection}
              type="button"
              onClick={() => void handleDelete()}
            >
              {deleting ? t("resources.deleting") : t("resources.delete")}
            </button>
            <button
              className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
              disabled={saving}
              type="button"
              onClick={() => void handleSave()}
            >
              {saving ? t("resources.saving") : t("resources.save")}
            </button>
          </div>
        </div>

        <aside className="mt-6 rounded-2xl bg-white/[0.03] p-5 text-sm leading-7 text-stone-300">
          <p>{t(guidanceIntroKey)}</p>
          <p className="mt-3">{t("resources.guidanceShared")}</p>
        </aside>
      </article>
    </section>
  );
}
