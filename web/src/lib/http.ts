export type ApiResponse<T> = {
  code: number;
  data?: T;
  msg?: string;
};

type RequestOptions = RequestInit & {
  allowEmptyBody?: boolean;
};

export async function request<T>(path: string, init?: RequestOptions): Promise<T> {
  const { allowEmptyBody = false, ...requestInit } = init ?? {};
  const headers = new Headers(requestInit.headers ?? {});

  if (!(requestInit.body instanceof FormData) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(path, {
    headers,
    ...requestInit,
  });

  const text = await response.text();
  if (!text.trim()) {
    if (allowEmptyBody && response.ok) {
      return undefined as T;
    }
    throw new Error(`Empty response from server (${response.status})`);
  }

  let payload: ApiResponse<T>;
  try {
    payload = JSON.parse(text) as ApiResponse<T>;
  } catch {
    throw new Error(`Invalid server response (${response.status})`);
  }

  if (!response.ok || payload.code !== 0 || !payload.data) {
    throw new Error(payload.msg || `Request failed: ${response.status}`);
  }

  return payload.data;
}
