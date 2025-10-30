const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

export const AUTH_INVALID_EVENT = "jobboard:auth-invalid";
export const FORCED_LOGOUT_MESSAGE_KEY = "jobboard:forced-logout-message";

export class ApiError extends Error {
  status?: number;

  constructor(message: string, status?: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

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

  let data: unknown = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = null;
    }
  }

  if (!response.ok) {
    let extractedError: unknown;
    if (data && typeof data === "object" && "error" in data) {
      extractedError = (data as { error?: unknown }).error;
    }

    const errorMessage =
      (typeof extractedError === "string" && extractedError.trim() !== "" ? extractedError : null) ??
      (response.statusText && response.statusText.trim() !== "" ? response.statusText : null) ??
      "API request failed";

    const error = new ApiError(errorMessage, response.status);

    if (error.status === 401) {
      const message = error.message || "セッションの有効期限が切れました。再ログインしてください。";
      window.sessionStorage.setItem(FORCED_LOGOUT_MESSAGE_KEY, message);
      window.dispatchEvent(new CustomEvent(AUTH_INVALID_EVENT, { detail: message }));
    }

    throw error;
  }

  return data as TResponse;
}
