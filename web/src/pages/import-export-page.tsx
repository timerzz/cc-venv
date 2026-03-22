import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { exportEnv, importEnv } from "../lib/api/import-export";
import type { AppShellContext } from "../types/app-shell";

export function ImportExportPage() {
  const { t } = useTranslation();
  const { activeEnv, refreshEnvs, selectEnv } = useOutletContext<AppShellContext>();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [forceImport, setForceImport] = useState(false);
  const [exportUrl, setExportURL] = useState<string | null>(null);
  const [exporting, setExporting] = useState(false);
  const [importing, setImporting] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleExport() {
    if (!activeEnv) {
      return;
    }

    setExporting(true);
    setError(null);
    setMessage(null);
    setExportURL(null);

    try {
      const result = await exportEnv(activeEnv.name);
      setExportURL(result.downloadUrl);
      setMessage(t("importExport.exportSuccess"));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("importExport.exportError"));
    } finally {
      setExporting(false);
    }
  }

  async function handleImport(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedFile) {
      return;
    }

    setImporting(true);
    setError(null);
    setMessage(null);

    try {
      const result = await importEnv(selectedFile, forceImport);
      await refreshEnvs();
      selectEnv(result.name);
      setMessage(t("importExport.importSuccess", { name: result.name }));
      setSelectedFile(null);
      setForceImport(false);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("importExport.importError"));
    } finally {
      setImporting(false);
    }
  }

  return (
    <section className="grid gap-5 xl:grid-cols-2">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">
          {t("common.importExport")}
        </div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("importExport.exportTitle")}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t("importExport.exportDescription")}</p>

        <div className="mt-8 rounded-3xl bg-white/[0.03] p-5">
          <p className="text-sm leading-7 text-stone-300">
            {t("importExport.currentEnv")}:{" "}
            <span className="text-ccv-ink">{activeEnv?.name ?? t("importExport.noEnv")}</span>
          </p>

          <div className="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center">
            <button
              className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
              disabled={!activeEnv || exporting}
              type="button"
              onClick={() => void handleExport()}
            >
              {exporting ? t("importExport.exporting") : t("importExport.exportAction")}
            </button>

            {exportUrl ? (
              <a
                className="rounded-2xl bg-ccv-panel-strong px-5 py-3 text-sm text-ccv-ink"
                href={exportUrl}
              >
                {t("importExport.downloadArchive")}
              </a>
            ) : null}
          </div>
        </div>
      </article>

      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">
          {t("importExport.importTitle")}
        </div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("importExport.importHeading")}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t("importExport.importDescription")}</p>

        <form className="mt-8 grid gap-5" onSubmit={handleImport}>
          <label className="grid gap-2">
            <span className="text-xs uppercase tracking-[0.28em] text-stone-400">
              {t("importExport.archiveFile")}
            </span>
            <input
              accept=".tar.gz,.tgz"
              className="rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink file:mr-4 file:rounded-xl file:border-0 file:bg-ccv-panel-strong file:px-3 file:py-2 file:text-sm file:text-ccv-ink"
              type="file"
              onChange={(event) => setSelectedFile(event.target.files?.[0] ?? null)}
            />
            <span className="text-xs text-slate-500">{t("importExport.archiveHint")}</span>
          </label>

          <label className="flex items-center gap-3 rounded-2xl bg-white/[0.03] px-4 py-3 text-sm text-stone-300">
            <input
              checked={forceImport}
              className="h-4 w-4"
              type="checkbox"
              onChange={(event) => setForceImport(event.target.checked)}
            />
            {t("importExport.force")}
          </label>

          {message ? (
            <div className="rounded-2xl bg-[#7e8c67]/20 px-4 py-3 text-sm text-[#e3edd6]">{message}</div>
          ) : null}
          {error ? (
            <div className="rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">{error}</div>
          ) : null}

          <div className="flex justify-end">
            <button
              className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
              disabled={!selectedFile || importing}
              type="submit"
            >
              {importing ? t("importExport.importing") : t("importExport.importAction")}
            </button>
          </div>
        </form>
      </article>
    </section>
  );
}
