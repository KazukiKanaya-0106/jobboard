const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type ApiRequestOptions<TBody = unknown> = {
  method?: HttpMethod;
  body?: TBody;
  token?: string | null;
  headers?: Record<string, string>;
};

export async function apiRequest<TResponse = unknown, TBody = unknown>(
  path: string,
  { method = "GET", body, headers = {}, token }: ApiRequestOptions<TBody> = {},
): Promise<TResponse> {
  const url = `${API_BASE_URL}${path}`;
  const init: RequestInit = {
    method,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
  };

  if (token) {
    init.headers = {
      ...init.headers,
      Authorization: `Bearer ${token}`,
    };
  }

  if (body !== undefined) {
    (init as RequestInit).body = JSON.stringify(body);
  }

  const response = await fetch(url, init);
  const text = await response.text();
  const data = text ? JSON.parse(text) : null;
  if (!response.ok) {
    const errorMessage =
      (data && typeof data === "object" && "error" in data && (data as { error?: string }).error) ||
      response.statusText;
    throw new Error(errorMessage || "API request failed");
  }
  return data as TResponse;
}
