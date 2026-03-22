import { request } from "../http";
import type { EnvDetail, EnvListItem } from "../../types/env";

export function listEnvs() {
  return request<{ envs: EnvListItem[] }>("/api/envs", {
    allowEmptyBody: true,
  }).then((data) => data ?? { envs: [] });
}

export function createEnv(name: string) {
  return request<{ name: string; path: string }>("/api/envs", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export function deleteEnv(envName: string) {
  return request<{ name: string }>(`/api/envs/${encodeURIComponent(envName)}`, {
    method: "DELETE",
  });
}

export function getEnv(envName: string) {
  return request<EnvDetail>(`/api/envs/${encodeURIComponent(envName)}`);
}

export function updateEnv(
  envName: string,
  payload: {
    name?: string;
    claudeMd?: string;
    envVars?: Record<string, string>;
  },
) {
  return request<{ name: string; renamed?: boolean }>(`/api/envs/${encodeURIComponent(envName)}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}
