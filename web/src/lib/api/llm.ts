import { request } from "../http";
import type { LLMConfig, LLMProvider } from "../../types/llm";

export function getLLMConfig(envName: string) {
  return request<LLMConfig>(`/api/envs/${encodeURIComponent(envName)}/llm`);
}

export function updateLLMConfig(envName: string, config: LLMConfig) {
  return request<{ name: string }>(`/api/envs/${encodeURIComponent(envName)}/llm`, {
    method: "PUT",
    body: JSON.stringify(config),
  });
}

export function listLLMProviders() {
  return request<{ providers: LLMProvider[] }>("/api/llm/providers");
}
