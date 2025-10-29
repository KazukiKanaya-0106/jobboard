import { useMutation } from '@tanstack/react-query'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import AuthForm from '../components/AuthForm'
import { register } from '../api'
import type { AuthCredentials } from '../schemas'
import { useAuth } from '../auth-context'

export default function RegisterPage() {
  const navigate = useNavigate()
  const { setAuth } = useAuth()
  const [apiError, setApiError] = useState<string | null>(null)

  const mutation = useMutation({
    mutationFn: (values: AuthCredentials) => register(values),
    onSuccess: (storedAuth) => {
      setApiError(null)
      setAuth(storedAuth)
      navigate('/', { replace: true })
    },
    onError: (error: unknown) => {
      setApiError(error instanceof Error ? error.message : '登録に失敗しました')
    },
  })

  return (
    <AuthForm
      mode="register"
      loading={mutation.isPending}
      apiError={apiError}
      onSubmit={(values) => mutation.mutate(values)}
    />
  )
}
