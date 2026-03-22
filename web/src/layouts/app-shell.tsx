import { useState } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import { EnvSidebar } from "../components/env-sidebar";
import { Modal } from "../components/modal";
import { SectionTabs } from "../components/section-tabs";
import { ConsoleShell } from "./console-shell";
import { useEnvs } from "../hooks/use-envs";
import { useTranslation } from "react-i18next";
import type { AppShellContext } from "../types/app-shell";

export function AppShell() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const {
    envs,
    searchTerm,
    activeEnv,
    loading,
    error,
    setSearchTerm,
    selectEnv,
    createEnvironment,
    deleteEnvironment,
    refreshEnvs,
  } = useEnvs();
  const [createOpen, setCreateOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [newEnvName, setNewEnvName] = useState("");
  const [actionError, setActionError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleCreateEnvironment() {
    const name = newEnvName.trim();
    if (!name) {
      setActionError(t("modals.createNameRequired"));
      return;
    }

    setSubmitting(true);
    setActionError(null);
    try {
      await createEnvironment(name);
      setCreateOpen(false);
      setNewEnvName("");
    } catch (err: unknown) {
      setActionError(err instanceof Error ? err.message : t("modals.createFailed"));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDeleteEnvironment() {
    if (!activeEnv) {
      return;
    }

    setSubmitting(true);
    setActionError(null);
    try {
      await deleteEnvironment(activeEnv.name);
      setDeleteOpen(false);
    } catch (err: unknown) {
      setActionError(err instanceof Error ? err.message : t("modals.deleteFailed"));
    } finally {
      setSubmitting(false);
    }
  }

  const outletContext: AppShellContext = {
    envs,
    searchTerm,
    activeEnv,
    loading,
    error,
    setSearchTerm,
    selectEnv,
    refreshEnvs,
    openCreateEnvironmentModal: () => {
      setActionError(null);
      setCreateOpen(true);
    },
    openDeleteEnvironmentModal: () => {
      setActionError(null);
      setDeleteOpen(true);
    },
  };

  return (
    <ConsoleShell
      sidebar={
        <EnvSidebar
          activeEnvName={activeEnv?.name ?? ""}
          envs={envs}
          error={error}
          loading={loading}
          searchTerm={searchTerm}
          onCreateEnv={() => {
            setActionError(null);
            setCreateOpen(true);
          }}
          onSearchChange={setSearchTerm}
          onSelectEnv={selectEnv}
        />
      }
    >
      <div className="grid gap-6">
        <header className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <div className="text-xs uppercase tracking-[0.38em] text-stone-400">
              {t("common.currentEnvironment")}
            </div>
            <h2 className="mt-3 font-serif text-5xl font-bold leading-none text-ccv-ink sm:text-6xl">
              {activeEnv?.name ?? t("common.noEnvironment")}
            </h2>
            <p className="mt-3 text-sm text-slate-500 sm:text-base">
              {activeEnv?.path ?? t("common.noEnvironmentDescription")}
            </p>
            <p className="mt-3 max-w-3xl text-sm leading-7 text-stone-400 sm:text-base">
              {t("dashboard.heroDescription")}
            </p>
          </div>

          <div className="flex w-full gap-3 lg:w-auto">
            <button
              className="flex-1 rounded-2xl bg-ccv-danger px-5 py-3 text-sm text-[#ffe5d8] lg:flex-none"
              type="button"
              onClick={() => {
                setActionError(null);
                setDeleteOpen(true);
              }}
            >
              {t("dashboard.deleteEnvironment")}
            </button>
            <button
              className="flex-1 rounded-2xl bg-ccv-panel-strong px-5 py-3 text-sm text-ccv-ink lg:flex-none"
              type="button"
              onClick={() => navigate("/import-export")}
            >
              {t("common.export")}
            </button>
            <button
              className="flex-1 rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white lg:flex-none"
              type="button"
              onClick={() => navigate("/import-export")}
            >
              {t("common.importCreate")}
            </button>
          </div>
        </header>

        <SectionTabs />
        <Outlet context={outletContext} />
      </div>

      <Modal
        description={t("modals.createDescription")}
        footer={
          <>
            <button
              className="rounded-2xl bg-ccv-panel-strong px-5 py-3 text-sm text-ccv-ink"
              type="button"
              onClick={() => {
                setCreateOpen(false);
                setNewEnvName("");
                setActionError(null);
              }}
            >
              {t("common.cancel")}
            </button>
            <button
              className="rounded-2xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] px-5 py-3 text-sm font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
              disabled={submitting}
              type="button"
              onClick={() => void handleCreateEnvironment()}
            >
              {submitting ? t("modals.creating") : t("common.createEnvironment")}
            </button>
          </>
        }
        open={createOpen}
        title={t("modals.createTitle")}
        onClose={() => {
          setCreateOpen(false);
          setNewEnvName("");
          setActionError(null);
        }}
      >
        <label className="grid gap-2">
          <span className="text-xs uppercase tracking-[0.28em] text-stone-400">
            {t("modals.environmentName")}
          </span>
          <input
            className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none placeholder:text-slate-500"
            placeholder={t("modals.environmentNamePlaceholder")}
            type="text"
            value={newEnvName}
            onChange={(event) => setNewEnvName(event.target.value)}
          />
        </label>
        {actionError ? (
          <div className="mt-4 rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">
            {actionError}
          </div>
        ) : null}
      </Modal>

      <Modal
        description={t("modals.deleteDescription", { name: activeEnv?.name ?? "" })}
        footer={
          <>
            <button
              className="rounded-2xl bg-ccv-panel-strong px-5 py-3 text-sm text-ccv-ink"
              type="button"
              onClick={() => {
                setDeleteOpen(false);
                setActionError(null);
              }}
            >
              {t("common.cancel")}
            </button>
            <button
              className="rounded-2xl bg-ccv-danger px-5 py-3 text-sm font-semibold text-[#ffe5d8] disabled:cursor-not-allowed disabled:opacity-60"
              disabled={submitting || !activeEnv}
              type="button"
              onClick={() => void handleDeleteEnvironment()}
            >
              {submitting ? t("modals.deleting") : t("dashboard.deleteEnvironment")}
            </button>
          </>
        }
        open={deleteOpen}
        title={t("modals.deleteTitle")}
        onClose={() => {
          setDeleteOpen(false);
          setActionError(null);
        }}
      >
        {actionError ? (
          <div className="rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">
            {actionError}
          </div>
        ) : null}
      </Modal>
    </ConsoleShell>
  );
}
