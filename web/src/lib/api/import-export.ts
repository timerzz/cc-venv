import { request } from "../http";
import type { ExportResponse, ImportResponse } from "../../types/import-export";

export function exportEnv(envName: string) {
  return request<ExportResponse>(`/api/envs/${encodeURIComponent(envName)}/export`, {
    method: "POST",
  });
}

export async function importEnv(file: File, force: boolean) {
  const formData = new FormData();
  formData.append("file", file);
  if (force) {
    formData.append("force", "true");
  }

  const response = await fetch("/api/envs/import", {
    method: "POST",
    body: formData,
  });

  const payload = (await response.json()) as { code: number; data?: ImportResponse; msg?: string };
  if (!response.ok || payload.code !== 0 || !payload.data) {
    throw new Error(payload.msg || `Request failed: ${response.status}`);
  }

  return payload.data;
}
