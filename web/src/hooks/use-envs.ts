import { useEffect, useMemo, useState } from "react";
import { createEnv, deleteEnv, listEnvs } from "../lib/api";
import type { EnvListItem } from "../types/env";

export function useEnvs() {
  const [envs, setEnvs] = useState<EnvListItem[]>([]);
  const [selectedEnvName, setSelectedEnvName] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  async function refreshEnvs() {
    setLoading(true);
    try {
      const data = await listEnvs();
      setEnvs(data.envs);
      setError(null);
      setSelectedEnvName((current) => {
        if (current && data.envs.some((env) => env.name === current)) {
          return current;
        }
        return data.envs[0]?.name ?? null;
      });
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to load environments");
      setEnvs([]);
      setSelectedEnvName(null);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void refreshEnvs();
  }, []);

  const filteredEnvs = useMemo(() => {
    const keyword = searchTerm.trim().toLowerCase();
    if (!keyword) {
      return envs;
    }

    return envs.filter((env) => env.name.toLowerCase().includes(keyword));
  }, [envs, searchTerm]);

  const activeEnv =
    filteredEnvs.find((env) => env.name === selectedEnvName) ??
    envs.find((env) => env.name === selectedEnvName) ??
    filteredEnvs[0] ??
    envs[0];

  async function createEnvironment(name: string) {
    await createEnv(name.trim());
    await refreshEnvs();
    setSelectedEnvName(name.trim());
  }

  async function deleteEnvironment(name: string) {
    await deleteEnv(name);
    await refreshEnvs();
  }

  return {
    envs: filteredEnvs,
    searchTerm,
    activeEnv,
    loading,
    error,
    setSearchTerm,
    selectEnv: setSelectedEnvName,
    createEnvironment,
    deleteEnvironment,
    refreshEnvs,
  };
}
