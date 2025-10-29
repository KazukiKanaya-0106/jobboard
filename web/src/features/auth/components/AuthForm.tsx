import Visibility from '@mui/icons-material/Visibility'
import VisibilityOff from '@mui/icons-material/VisibilityOff'
import { Alert, Button, IconButton, InputAdornment, Stack, TextField, Typography, Link } from '@mui/material'
import { useState } from 'react'
import { Link as RouterLink } from 'react-router-dom'
import type { AuthCredentials } from '../schemas'
import { authCredentialsSchema, authRegistrationSchema } from '../schemas'

type AuthFormProps = {
  mode: 'login' | 'register'
  onSubmit: (values: AuthCredentials) => void
  loading?: boolean
  apiError?: string | null
}

type FormValues = {
  clusterId: string
  password: string
  confirmPassword: string
}

type FormErrors = Partial<Record<keyof FormValues, string>>

export default function AuthForm({ mode, onSubmit, loading, apiError }: AuthFormProps) {
  const [values, setValues] = useState<FormValues>({ clusterId: '', password: '', confirmPassword: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [errors, setErrors] = useState<FormErrors>({})

  const handleChange = (field: keyof FormValues) => (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues((prev) => ({ ...prev, [field]: event.target.value }))
  }

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const baseValues = { clusterId: values.clusterId, password: values.password }
    const parseResult =
      mode === 'register'
        ? authRegistrationSchema.safeParse({ ...baseValues, confirmPassword: values.confirmPassword })
        : authCredentialsSchema.safeParse(baseValues)
    if (!parseResult.success) {
      const fieldErrors: FormErrors = {}
      parseResult.error.issues.forEach((issue) => {
        const field = issue.path[0]
        if (typeof field === 'string') {
          fieldErrors[field as keyof FormValues] = issue.message
        }
      })
      setErrors(fieldErrors)
      return
    }
    setErrors({})
    onSubmit({
      clusterId: parseResult.data.clusterId,
      password: parseResult.data.password,
    })
  }

  return (
    <Stack spacing={3} component="form" onSubmit={handleSubmit} noValidate>
      <Typography variant="h5">{mode === 'login' ? 'ログイン' : 'クラスタ登録'}</Typography>

      {apiError ? (
        <Alert severity="error" variant="filled">
          {apiError}
        </Alert>
      ) : null}

      <TextField
        label="クラスタID"
        value={values.clusterId}
        onChange={handleChange('clusterId')}
        error={Boolean(errors.clusterId)}
        helperText={errors.clusterId}
        autoFocus
        required
        fullWidth
        autoComplete="username"
      />

      <TextField
        label="パスワード"
        type={showPassword ? 'text' : 'password'}
        value={values.password}
        onChange={handleChange('password')}
        error={Boolean(errors.password)}
        helperText={errors.password}
        required
        fullWidth
        autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={() => setShowPassword((prev) => !prev)} edge="end">
                {showPassword ? <VisibilityOff /> : <Visibility />}
              </IconButton>
            </InputAdornment>
          ),
        }}
      />

      {mode === 'register' ? (
        <TextField
          label="パスワード（確認）"
          type={showPassword ? 'text' : 'password'}
          value={values.confirmPassword}
          onChange={handleChange('confirmPassword')}
          error={Boolean(errors.confirmPassword)}
          helperText={errors.confirmPassword}
          required
          fullWidth
          autoComplete="new-password"
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton onClick={() => setShowPassword((prev) => !prev)} edge="end">
                  {showPassword ? <VisibilityOff /> : <Visibility />}
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
      ) : null}

      <Button type="submit" variant="contained" size="large" disabled={loading}>
        {mode === 'login' ? 'ログイン' : '登録してログイン'}
      </Button>

      <Typography variant="body2" color="text.secondary">
        {mode === 'login' ? (
          <>
            クラスタをお持ちでない場合は{' '}
            <Link component={RouterLink} to="/auth/register">
              新規登録
            </Link>
            してください。
          </>
        ) : (
          <>
            すでにクラスタをお持ちですか？{' '}
            <Link component={RouterLink} to="/auth/login">
              ログイン
            </Link>
          </>
        )}
      </Typography>
    </Stack>
  )
}
