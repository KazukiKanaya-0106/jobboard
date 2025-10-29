import type { PropsWithChildren } from "react";
import { createContext, useCallback, useContext, useMemo, useState } from "react";
import { loadAuth, saveAuth, type StoredAuth } from "../../lib/storage";

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
