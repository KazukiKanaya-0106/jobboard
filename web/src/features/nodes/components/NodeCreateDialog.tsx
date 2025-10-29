import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Stack,
} from '@mui/material'
import { useState } from 'react'
import { createNodeRequestSchema, type CreateNodeRequest } from '../schemas'

type NodeCreateDialogProps = {
  open: boolean
  onClose: () => void
  onSubmit: (values: CreateNodeRequest) => Promise<void>
  loading?: boolean
  apiError?: string | null
}

type FormErrors = Partial<Record<keyof CreateNodeRequest, string>>

export default function NodeCreateDialog({ open, onClose, onSubmit, loading, apiError }: NodeCreateDialogProps) {
  const [nodeName, setNodeName] = useState('')
  const [errors, setErrors] = useState<FormErrors>({})

  const handleClose = () => {
    setErrors({})
    setNodeName('')
    onClose()
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const parseResult = createNodeRequestSchema.safeParse({ nodeName })
    if (!parseResult.success) {
      const fieldErrors: FormErrors = {}
      parseResult.error.issues.forEach((issue) => {
        const field = issue.path[0]
        if (typeof field === 'string') {
          fieldErrors[field as keyof CreateNodeRequest] = issue.message
        }
      })
      setErrors(fieldErrors)
      return
    }
    setErrors({})
    await onSubmit(parseResult.data)
    setNodeName('')
  }

  return (
    <Dialog open={open} onClose={handleClose} fullWidth maxWidth="sm">
      <form onSubmit={handleSubmit} noValidate>
        <DialogTitle>ノードを追加</DialogTitle>
        <DialogContent>
          <Stack spacing={3} sx={{ mt: 1 }}>
            <TextField
              label="ノード名"
              value={nodeName}
              onChange={(event) => setNodeName(event.target.value)}
              error={Boolean(errors.nodeName)}
              helperText={errors.nodeName || apiError}
              disabled={loading}
              required
              fullWidth
            />
          </Stack>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={handleClose} disabled={loading}>
            キャンセル
          </Button>
          <Button type="submit" variant="contained" disabled={loading}>
            追加
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
