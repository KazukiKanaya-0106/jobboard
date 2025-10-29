import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  Typography,
} from '@mui/material'

type NodeDeleteDialogProps = {
  open: boolean
  nodeName?: string
  onClose: () => void
  onConfirm: () => Promise<void> | void
  loading?: boolean
  error?: string | null
}

export default function NodeDeleteDialog({
  open,
  nodeName,
  onClose,
  onConfirm,
  loading,
  error,
}: NodeDeleteDialogProps) {
  return (
    <Dialog open={open} onClose={loading ? undefined : onClose} maxWidth="xs" fullWidth>
      <DialogTitle>ノードを削除</DialogTitle>
      <DialogContent dividers>
        <Stack spacing={2}>
          <Typography>
            ノード
            <Typography component="span" fontWeight={600} color="primary" sx={{ mx: 0.5 }}>
              {nodeName}
            </Typography>
            を削除しますか？
          </Typography>
          <Typography variant="body2" color="text.secondary">
            この操作は取り消せません。ノードの実行中ジョブがある場合はクリアされます。
          </Typography>
          {error ? (
            <Alert severity="error" variant="filled">
              {error}
            </Alert>
          ) : null}
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button onClick={onClose} disabled={Boolean(loading)}>
          キャンセル
        </Button>
        <Button onClick={onConfirm} color="error" variant="contained" disabled={Boolean(loading)}>
          削除する
        </Button>
      </DialogActions>
    </Dialog>
  )
}
