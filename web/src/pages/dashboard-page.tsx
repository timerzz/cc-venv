import { useEffect, useState } from "react";
import { useOutletContext } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { getEnv } from "../lib/api/env";
import { getLLMConfig } from "../lib/api/llm";
import { StatusCard } from "../components/status-card";
import type { AppShellContext } from "../types/app-shell";
import type { EnvDetail } from "../types/env";
import type { LLMConfig } from "../types/llm";

export function DashboardPage() {
  const { t } = useTranslation();
  const { activeEnv, loading, error } = useOutletContext<AppShellContext>();
  const [envDetail, setEnvDetail] = useState<EnvDetail | null>(null);
  const [llmConfig, setLLMConfig] = useState<LLMConfig | null>(null);
  const [summaryError, setSummaryError] = useState<string | null>(null);

  useEffect(() => {
    if (!activeEnv) {
      setEnvDetail(null);
      setLLMConfig(null);
      setSummaryError(null);
      return;
    }

    let cancelled = false;
    setSummaryError(null);
    setEnvDetail(null);
    setLLMConfig(null);

    void Promise.all([getEnv(activeEnv.name), getLLMConfig(activeEnv.name)])
      .then(([detail, llm]) => {
        if (!cancelled) {
          setEnvDetail(detail);
          setLLMConfig(llm);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setSummaryError(err instanceof Error ? err.message : t("dashboard.summaryLoadError"));
        }
      });

    return () => {
      cancelled = true;
    };
  }, [activeEnv, t]);

  const llmTitle = getLLMTitle(llmConfig, t);
  const llmDescription = getLLMDescription(llmConfig, t);
  const mcpCount = Object.keys(envDetail?.mcpServers ?? {}).length;
  const skillCount = envDetail?.resources.skills.length ?? 0;
  const agentCount = envDetail?.resources.agents.length ?? 0;
  const commandCount = envDetail?.resources.commands.length ?? 0;
  const ruleCount = envDetail?.resources.rules.length ?? 0;
  const envVarCount = Object.keys(envDetail?.envVars ?? {}).length;

  return (
    <section className="grid gap-5">
      <div className="grid gap-5">
        {loading ? (
          <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-sm text-stone-300">
            {t("dashboard.loadingOverview")}
          </article>
        ) : !activeEnv ? (
          <article className="rounded-[1.625rem] border border-dashed border-white/10 bg-white/[0.03] p-7 text-sm text-stone-400">
            {t("dashboard.emptyOverview")}
          </article>
        ) : (
          <>
            <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <StatusCard
                label={t("common.llm")}
                title={llmTitle}
                description={llmDescription}
                badge={llmConfig?.models.default || llmConfig?.baseUrl ? t("common.healthy") : t("dashboard.notConfigured")}
              />
              <StatusCard
                label={t("common.mcp")}
                title={t("dashboard.status.mcpTitle", { count: mcpCount })}
                description={t("dashboard.status.mcpDescription", { count: mcpCount })}
                badge={mcpCount > 0 ? t("common.ready") : t("dashboard.notConfigured")}
              />
              <StatusCard
                label={t("common.skills")}
                title={t("dashboard.status.skillsTitle", { count: skillCount })}
                description={t("dashboard.status.skillsDescription", { count: skillCount })}
                badge={skillCount > 0 ? t("common.needsReview") : t("dashboard.notConfigured")}
              />
              <StatusCard
                label={t("common.envVars")}
                title={t("dashboard.status.envVarsTitle", { count: envVarCount })}
                description={t("dashboard.status.envVarsDescription", { count: envVarCount })}
                badge={envVarCount > 0 ? t("common.ready") : t("dashboard.notConfigured")}
              />
              <StatusCard
                label={t("common.agents")}
                title={t("dashboard.status.agentsTitle", { count: agentCount })}
                description={t("dashboard.status.agentsDescription", { count: agentCount })}
                badge={agentCount > 0 ? t("common.ready") : t("dashboard.notConfigured")}
              />
              <StatusCard
                label={t("common.commands")}
                title={t("dashboard.status.commandsTitle", { count: commandCount })}
                description={t("dashboard.status.commandsDescription", { count: commandCount })}
                badge={commandCount > 0 ? t("common.ready") : t("dashboard.notConfigured")}
              />
              <StatusCard
                label={t("common.rules")}
                title={t("dashboard.status.rulesTitle", { count: ruleCount })}
                description={t("dashboard.status.rulesDescription", { count: ruleCount })}
                badge={ruleCount > 0 ? t("common.ready") : t("dashboard.notConfigured")}
              />
            </section>

            {error || summaryError ? (
              <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
                <div className="rounded-2xl bg-[#e17344]/15 px-4 py-3 text-sm text-[#ffddcf]">
                  {t("dashboard.apiUnavailable", { message: error ?? summaryError })}
                </div>
              </article>
            ) : null}
          </>
        )}
      </div>
    </section>
  );
}

function getLLMTitle(llmConfig: LLMConfig | null, t: (key: string, options?: Record<string, unknown>) => string) {
  if (!llmConfig) {
    return t("dashboard.notConfigured");
  }

  if (llmConfig.models.default) {
    return llmConfig.models.default;
  }

  if (llmConfig.baseUrl) {
    try {
      return new URL(llmConfig.baseUrl).host;
    } catch {
      return llmConfig.baseUrl;
    }
  }

  return t("dashboard.notConfigured");
}

function getLLMDescription(
  llmConfig: LLMConfig | null,
  t: (key: string, options?: Record<string, unknown>) => string,
) {
  if (!llmConfig) {
    return t("dashboard.status.llmEmptyDescription");
  }

  if (llmConfig.baseUrl) {
    return t("dashboard.status.llmDescription", { value: llmConfig.baseUrl });
  }

  if (llmConfig.models.sonnet || llmConfig.models.opus || llmConfig.models.haiku) {
    return t("dashboard.status.llmModelAliases");
  }

  return t("dashboard.status.llmEmptyDescription");
}
