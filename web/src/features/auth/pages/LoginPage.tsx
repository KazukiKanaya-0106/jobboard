import { useMutation } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import AuthForm from "../components/AuthForm";
import { login } from "../api";
import type { AuthCredentials } from "../schemas";
import { useAuth } from "../AuthContext";
import { FORCED_LOGOUT_MESSAGE_KEY } from "../../../lib/apiCient";
import { resolveErrorMessage } from "../../../lib/errorCatalog";

type LocationState = {
  from?: {
    pathname: string;
  };
};

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { setAuth } = useAuth();
  const [apiError, setApiError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: (values: AuthCredentials) => login(values),
    onSuccess: (storedAuth) => {
      setApiError(null);
      setAuth(storedAuth);
      const state = location.state as LocationState | undefined;
      const redirectTo = state?.from?.pathname ?? "/";
      navigate(redirectTo, { replace: true });
    },
    onError: (error: unknown) => {
      setApiError(resolveErrorMessage(error, "ログインに失敗しました"));
    },
  });

  useEffect(() => {
    const message = window.sessionStorage.getItem(FORCED_LOGOUT_MESSAGE_KEY);
    if (message) {
      setApiError(message);
      window.sessionStorage.removeItem(FORCED_LOGOUT_MESSAGE_KEY);
    }
  }, []);

  return (
    <AuthForm
      mode="login"
      loading={mutation.isPending}
      apiError={apiError}
      onSubmit={(values) => mutation.mutate(values)}
    />
  );
}
