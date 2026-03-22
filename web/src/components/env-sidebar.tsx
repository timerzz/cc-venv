import type { EnvListItem } from "../types/env";
import { useTranslation } from "react-i18next";

type EnvSidebarProps = {
  envs: EnvListItem[];
  searchTerm: string;
  activeEnvName: string;
  loading?: boolean;
  error?: string | null;
  onSearchChange: (value: string) => void;
  onSelectEnv: (name: string) => void;
  onCreateEnv: () => void | Promise<void>;
};

export function EnvSidebar({
  envs,
  searchTerm,
  activeEnvName,
  loading = false,
  error = null,
  onSearchChange,
  onSelectEnv,
  onCreateEnv,
}: EnvSidebarProps) {
  const { t } = useTranslation();

  return (
    <aside className="flex h-fit min-h-0 flex-col gap-4 rounded-[1.75rem] bg-white/[0.02] p-5 lg:p-6">
      <div className="text-xs uppercase tracking-[0.38em] text-stone-400">
        {t("dashboard.environments")}
      </div>
      <input
        className="rounded-2xl bg-ccv-panel-soft px-4 py-3 text-sm text-ccv-ink outline-none placeholder:text-slate-500"
        placeholder={t("dashboard.searchPlaceholder")}
        type="text"
        value={searchTerm}
        onChange={(event) => onSearchChange(event.target.value)}
      />

      <div className="grid gap-3">
        {loading ? (
          <div className="rounded-[1.25rem] bg-ccv-panel-strong p-5 text-sm text-stone-300">
            {t("dashboard.loadingEnvs")}
          </div>
        ) : error ? (
          <div className="rounded-[1.25rem] bg-[#e17344]/15 p-5 text-sm text-[#ffddcf]">
            {t("dashboard.loadEnvsError", { message: error })}
          </div>
        ) : envs.length === 0 ? (
          <div className="rounded-[1.25rem] border border-dashed border-white/10 bg-white/[0.03] p-5 text-sm text-stone-400">
            {t("dashboard.emptyEnvs")}
          </div>
        ) : (
          envs.map((env) => {
            const active = env.name === activeEnvName;
            return (
              <button
                key={env.name}
                className={
                  active
                    ? "rounded-[1.25rem] bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] p-5 text-left text-white"
                    : "rounded-[1.25rem] bg-ccv-panel-strong p-5 text-left text-ccv-ink"
                }
                type="button"
                onClick={() => onSelectEnv(env.name)}
              >
                <div className="mb-2 text-[0.72rem] uppercase tracking-[0.22em] text-white/80">
                  {active ? t("dashboard.activeEnv") : t("dashboard.envLabel")}
                </div>
                <h3 className="mb-2 text-[1.8rem] font-semibold leading-none">{env.name}</h3>
                <p className="text-sm text-white/75">
                  {t("dashboard.resourcesConfigured", { count: countResources(env.resources) })}
                </p>
              </button>
            );
          })
        )}
      </div>

      <button
        className="mt-auto rounded-[1.125rem] bg-[linear-gradient(180deg,#8c9c72_0%,#7e8c67_100%)] px-4 py-4 font-semibold text-stone-50"
        type="button"
        onClick={() => void onCreateEnv()}
      >
        + {t("common.createEnvironment")}
      </button>
    </aside>
  );
}

function countResources(resources: EnvListItem["resources"]) {
  return (
    resources.skills +
    resources.agents +
    resources.commands +
    resources.rules +
    resources.hooks +
    resources.mcpServers
  );
}
