import { request } from "../http";
import type { AddMCPRequest, MCPListResponse } from "../../types/mcp";

export function listMCP(envName: string) {
  return request<MCPListResponse>(`/api/envs/${encodeURIComponent(envName)}/mcp`);
}

export function addMCP(envName: string, payload: AddMCPRequest) {
  return request<{ name: string; config: AddMCPRequest["config"] }>(
    `/api/envs/${encodeURIComponent(envName)}/mcp`,
    {
      method: "POST",
      body: JSON.stringify(payload),
    },
  );
}

export function deleteMCP(envName: string, serverName: string) {
  return request<{ name: string }>(
    `/api/envs/${encodeURIComponent(envName)}/mcp/${encodeURIComponent(serverName)}`,
    {
      method: "DELETE",
    },
  );
}
