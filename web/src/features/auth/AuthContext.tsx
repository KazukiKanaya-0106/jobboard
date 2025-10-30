import type { PropsWithChildren } from "react";
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { loadAuth, saveAuth, type StoredAuth } from "../../lib/storage";
import { AUTH_INVALID_EVENT, FORCED_LOGOUT_MESSAGE_KEY } from "../../lib/apiCient";

type AuthState = StoredAuth | null;

type AuthContextValue = {
  auth: AuthState;
  isAuthenticated: boolean;
  setAuth: (auth: StoredAuth) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const [auth, setAuthState] = useState<AuthState>(() => loadAuth());

  const setAuth = useCallback((next: StoredAuth) => {
    setAuthState(next);
    saveAuth(next);
  }, []);

  const logout = useCallback(() => {
    setAuthState(null);
    saveAuth(null);
  }, []);

  useEffect(() => {
    const handler = (event: Event) => {
      const detail = (event as CustomEvent<string | undefined>).detail;
      if (detail) {
        window.sessionStorage.setItem(
          FORCED_LOGOUT_MESSAGE_KEY,
          detail || "セッションの有効期限が切れました。再ログインしてください。",
        );
      }
      logout();
    };

    window.addEventListener(AUTH_INVALID_EVENT, handler as EventListener);
    return () => window.removeEventListener(AUTH_INVALID_EVENT, handler as EventListener);
  }, [logout]);

  const value = useMemo<AuthContextValue>(
    () => ({
      auth,
      isAuthenticated: Boolean(auth?.token),
      setAuth,
      logout,
    }),
    [auth, logout, setAuth],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
