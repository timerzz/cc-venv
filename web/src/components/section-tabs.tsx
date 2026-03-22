import { NavLink } from "react-router-dom";
import { useTranslation } from "react-i18next";

type TabItem = {
  to: string;
  labelKey: string;
  end?: boolean;
};

const tabs: TabItem[] = [
  { to: "/", labelKey: "common.overview", end: true },
  { to: "/llm", labelKey: "common.llm" },
  { to: "/mcp", labelKey: "common.mcp" },
  { to: "/skills", labelKey: "common.skills" },
  { to: "/agents", labelKey: "common.agents" },
  { to: "/commands", labelKey: "common.commands" },
  { to: "/rules", labelKey: "common.rules" },
  { to: "/env-vars", labelKey: "common.envVars" },
  { to: "/notes", labelKey: "common.notes" },
  { to: "/import-export", labelKey: "common.importExport" },
];

export function SectionTabs() {
  const { t } = useTranslation();

  return (
    <nav
      aria-label="Primary sections"
      className="flex flex-wrap gap-3 rounded-[1.25rem] bg-[#1e252a] p-3"
    >
      {tabs.map((tab) => (
        <NavLink
          key={tab.to}
          className={({ isActive }) =>
            isActive
              ? "rounded-xl bg-stone-100 px-4 py-2.5 text-sm font-semibold text-slate-900"
              : "px-3 py-2.5 text-sm text-stone-400"
          }
          end={tab.end}
          to={tab.to}
        >
          {t(tab.labelKey)}
        </NavLink>
      ))}
    </nav>
  );
}
