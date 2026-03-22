import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { deleteSkill, listSkills, uploadSkill } from "../lib/api/skills";
import type { AppShellContext } from "../types/app-shell";
import type { SkillInfo } from "../types/skill";

export function SkillsPage() {
  const { t } = useTranslation();
  const { activeEnv } = useOutletContext<AppShellContext>();
  const [skills, setSkills] = useState<SkillInfo[]>([]);
  const [file, setFile] = useState<File | null>(null);
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

    void listSkills(activeEnv.name)
      .then((data) => {
        if (!cancelled) {
          setSkills(Array.isArray(data.skills) ? data.skills : []);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : t("skills.loadError"));
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
    const form = event.currentTarget;
    if (!activeEnv || !file) {
      return;
    }

    setSaving(true);
    setError(null);
    setMessage(null);

    try {
      const skill = await uploadSkill(activeEnv.name, file);
      setSkills((current) =>
        [...current.filter((item) => item.name !== skill.name), skill].sort((a, b) => a.name.localeCompare(b.name)),
      );
      setFile(null);
      form.reset();
      setMessage(t("skills.addSuccess", { name: skill.name }));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("skills.addError"));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(skillName: string) {
    if (!activeEnv) {
      return;
    }

    setDeletingName(skillName);
    setError(null);
    setMessage(null);

    try {
      await deleteSkill(activeEnv.name, skillName);
      setSkills((current) => current.filter((skill) => skill.name !== skillName));
      setMessage(t("skills.deleteSuccess", { name: skillName }));
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : t("skills.deleteError"));
    } finally {
      setDeletingName(null);
    }
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[minmax(320px,0.9fr)_minmax(0,1.3fr)]">
      <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
        <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("common.skills")}</div>
        <h3 className="mt-3 font-serif text-4xl font-bold">{t("skills.title")}</h3>
        <p className="mt-3 text-base leading-7 text-slate-400">{t("skills.description")}</p>

        <form className="mt-8 grid gap-5" onSubmit={handleSubmit}>
          <label className="grid gap-2">
            <span className="text-xs uppercase tracking-[0.28em] text-stone-400">
              {t("skills.archive")}
            </span>
            <input
              accept=".zip,application/zip"
              className="w-full rounded-2xl border border-white/8 bg-[#151b1f] px-4 py-3 text-ccv-ink outline-none file:mr-4 file:rounded-xl file:border-0 file:bg-[#e17344] file:px-4 file:py-2.5 file:text-sm file:font-semibold file:text-white"
              type="file"
              onChange={(event) => setFile(event.target.files?.[0] ?? null)}
            />
            <span className="text-xs text-slate-500">{t("skills.archiveHint")}</span>
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
              disabled={saving || !file}
              type="submit"
            >
              {saving ? t("skills.adding") : t("skills.add")}
            </button>
          </div>
        </form>
      </article>

      <section className="grid gap-5">
        <article className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
          <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("skills.installed")}</div>
          <h3 className="mt-3 font-serif text-4xl font-bold">{t("skills.listTitle")}</h3>
          <p className="mt-3 text-base leading-7 text-slate-400">{t("skills.listDescription")}</p>

          {loading ? (
            <div className="mt-8 rounded-2xl border border-white/8 bg-white/[0.03] p-5 text-sm text-stone-400">
              {t("skills.loading")}
            </div>
          ) : skills.length === 0 ? (
            <div className="mt-8 rounded-2xl border border-dashed border-white/10 bg-white/[0.03] p-5 text-sm text-stone-400">
              {t("skills.empty")}
            </div>
          ) : (
            <div className="mt-8 grid gap-4">
              {skills.map((skill) => (
                <article key={skill.name} className="rounded-3xl bg-white/[0.03] p-5">
                  <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
                    <div className="min-w-0">
                      <div className="text-xs uppercase tracking-[0.28em] text-stone-500">
                        {t("skills.skill")}
                      </div>
                      <h4 className="mt-2 text-2xl font-semibold text-ccv-ink">{skill.name}</h4>
                      <p className="mt-3 break-all text-sm text-slate-400">
                        {t("skills.path")}: <span className="text-stone-200">{skill.path}</span>
                      </p>
                    </div>

                    <button
                      className="rounded-2xl bg-ccv-danger px-4 py-2.5 text-sm text-[#ffe5d8] disabled:cursor-not-allowed disabled:opacity-60"
                      disabled={deletingName === skill.name}
                      type="button"
                      onClick={() => void handleDelete(skill.name)}
                    >
                      {deletingName === skill.name ? t("skills.deleting") : t("skills.delete")}
                    </button>
                  </div>
                </article>
              ))}
            </div>
          )}
        </article>

        <aside className="rounded-[1.625rem] bg-ccv-panel-soft p-7 text-ccv-ink">
          <div className="text-xs uppercase tracking-[0.38em] text-stone-400">{t("skills.guidance")}</div>
          <div className="mt-4 rounded-2xl bg-white/[0.03] p-5 text-sm leading-7 text-stone-300">
            <p>{t("skills.guidanceIntro")}</p>
            <p className="mt-3">{t("skills.guidanceArchive")}</p>
          </div>
        </aside>
      </section>
    </section>
  );
}
