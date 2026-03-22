export type MCPServer = {
  type?: string;
  url?: string;
  headers?: Record<string, string>;
  command?: string;
  args?: string[];
  env?: Record<string, string>;
};

export type MCPListResponse = {
  servers: Record<string, MCPServer>;
};

export type AddMCPRequest = {
  name: string;
  config: MCPServer;
};
