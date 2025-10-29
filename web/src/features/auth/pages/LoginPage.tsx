import { useMutation } from '@tanstack/react-query'
import { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import AuthForm from '../components/AuthForm'
import { login } from '../api'
import type { AuthCredentials } from '../schemas'
import { useAuth } from '../AuthContext'

type LocationState = {
  from?: {
    pathname: string
  }
}

export default function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const { setAuth } = useAuth()
  const [apiError, setApiError] = useState<string | null>(null)

  const mutation = useMutation({
    mutationFn: (values: AuthCredentials) => login(values),
    onSuccess: (storedAuth) => {
      setApiError(null)
      setAuth(storedAuth)
      const state = location.state as LocationState | undefined
      const redirectTo = state?.from?.pathname ?? '/'
      navigate(redirectTo, { replace: true })
    },
    onError: (error: unknown) => {
      setApiError(error instanceof Error ? error.message : 'ログインに失敗しました')
    },
  })

  return (
    <AuthForm mode="login" loading={mutation.isPending} apiError={apiError} onSubmit={(values) => mutation.mutate(values)} />
  )
}
