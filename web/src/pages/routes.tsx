import { AgentsPage } from "./agents-page";
import { CommandsPage } from "./commands-page";
import { createBrowserRouter } from "react-router-dom";
import { AppShell } from "../layouts/app-shell";
import { DashboardPage } from "./dashboard-page";
import { EnvVarsPage } from "./env-vars-page";
import { ImportExportPage } from "./import-export-page";
import { LLMPage } from "./llm-page";
import { MCPPage } from "./mcp-page";
import { NotesPage } from "./notes-page";
import { RulesPage } from "./rules-page";
import { SkillsPage } from "./skills-page";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppShell />,
    children: [
      { index: true, element: <DashboardPage /> },
      {
        path: "llm",
        element: <LLMPage />,
      },
      {
        path: "mcp",
        element: <MCPPage />,
      },
      {
        path: "skills",
        element: <SkillsPage />,
      },
      {
        path: "agents",
        element: <AgentsPage />,
      },
      {
        path: "commands",
        element: <CommandsPage />,
      },
      {
        path: "rules",
        element: <RulesPage />,
      },
      {
        path: "env-vars",
        element: <EnvVarsPage />,
      },
      {
        path: "notes",
        element: <NotesPage />,
      },
      {
        path: "import-export",
        element: <ImportExportPage />,
      },
    ],
  },
]);
