import { request } from "../http";
import type { ResourceFile, ResourceKind, ResourceListResponse } from "../../types/resource-file";

export function listResourceFiles(envName: string, kind: ResourceKind) {
  return request<ResourceListResponse>(
    `/api/envs/${encodeURIComponent(envName)}/resources/${encodeURIComponent(kind)}`,
  );
}

export function getResourceFile(envName: string, kind: ResourceKind, name: string) {
  return request<ResourceFile>(
    `/api/envs/${encodeURIComponent(envName)}/resources/${encodeURIComponent(kind)}/content?name=${encodeURIComponent(name)}`,
  );
}

export function upsertResourceFile(envName: string, kind: ResourceKind, file: ResourceFile) {
  return request<ResourceFile>(
    `/api/envs/${encodeURIComponent(envName)}/resources/${encodeURIComponent(kind)}/content`,
    {
      method: "PUT",
      body: JSON.stringify(file),
    },
  );
}

export function deleteResourceFile(envName: string, kind: ResourceKind, name: string) {
  return request<{ name: string }>(
    `/api/envs/${encodeURIComponent(envName)}/resources/${encodeURIComponent(kind)}/content?name=${encodeURIComponent(name)}`,
    {
      method: "DELETE",
    },
  );
}
