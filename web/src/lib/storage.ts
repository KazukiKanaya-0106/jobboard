const AUTH_STORAGE_KEY = "jobboard_auth";

export type StoredAuth = {
  token: string;
  clusterId: string;
};

export function loadAuth(): StoredAuth | null {
  try {
    const raw = window.localStorage.getItem(AUTH_STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as StoredAuth;
    if (typeof parsed?.token === "string" && typeof parsed?.clusterId === "string") {
      return parsed;
    }
    return null;
  } catch {
    return null;
  }
}

export function saveAuth(auth: StoredAuth | null) {
  if (!auth) {
    window.localStorage.removeItem(AUTH_STORAGE_KEY);
    return;
  }
  window.localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(auth));
}
