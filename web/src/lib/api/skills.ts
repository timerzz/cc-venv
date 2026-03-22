import { request } from "../http";
import type { SkillInfo, SkillListResponse } from "../../types/skill";

export function listSkills(envName: string) {
  return request<SkillListResponse>(`/api/envs/${encodeURIComponent(envName)}/skills`);
}

export function uploadSkill(envName: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);

  return request<SkillInfo>(`/api/envs/${encodeURIComponent(envName)}/skills`, {
    method: "POST",
    body: formData,
  });
}

export function deleteSkill(envName: string, skillName: string) {
  return request<{ name: string }>(
    `/api/envs/${encodeURIComponent(envName)}/skills/${encodeURIComponent(skillName)}`,
    {
      method: "DELETE",
    },
  );
}
