import { z } from 'zod'

export const authCredentialsSchema = z.object({
  clusterId: z
    .string()
    .min(1, 'クラスタIDを入力してください')
    .max(64, 'クラスタIDは64文字以内で入力してください'),
  password: z
    .string()
    .min(1, 'パスワードを入力してください')
    .max(128, 'パスワードは128文字以内で入力してください'),
})

export type AuthCredentials = z.infer<typeof authCredentialsSchema>

export const authRegistrationSchema = authCredentialsSchema
  .extend({
    confirmPassword: z
      .string()
      .min(1, '確認用パスワードを入力してください')
      .max(128, 'パスワードは128文字以内で入力してください'),
  })
  .superRefine((data, ctx) => {
    if (data.password !== data.confirmPassword) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['confirmPassword'],
        message: 'パスワードが一致しません',
      })
    }
  })

export type AuthRegistrationCredentials = z.infer<typeof authRegistrationSchema>

export const authResponseSchema = z.object({
  cluster_id: z.string(),
  token: z.string(),
  expires_at: z.number().optional(),
})

export type AuthResponseDto = z.infer<typeof authResponseSchema>
