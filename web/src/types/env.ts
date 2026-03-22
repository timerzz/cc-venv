export type ResourceCounts = {
  skills: number;
  agents: number;
  commands: number;
  rules: number;
  hooks: number;
  mcpServers: number;
};

export type EnvListItem = {
  name: string;
  path: string;
  resources: ResourceCounts;
};

export type EnvDetail = {
  name: string;
  path: string;
  claudeMd?: string;
  settings?: Record<string, unknown>;
  envVars?: Record<string, string>;
  mcpServers?: Record<
    string,
    {
      command: string;
      args?: string[];
      env?: Record<string, string>;
    }
  >;
  resources: {
    skills: string[];
    agents: string[];
    commands: string[];
    rules: string[];
    hooks: string[];
    plugins: string[];
  };
};
