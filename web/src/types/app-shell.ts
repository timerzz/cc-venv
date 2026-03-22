import type { EnvListItem } from "./env";

export type AppShellContext = {
  envs: EnvListItem[];
  searchTerm: string;
  activeEnv?: EnvListItem;
  loading: boolean;
  error?: string | null;
  setSearchTerm: (value: string) => void;
  selectEnv: (name: string) => void;
  refreshEnvs: () => Promise<void>;
  openCreateEnvironmentModal: () => void;
  openDeleteEnvironmentModal: () => void;
};
