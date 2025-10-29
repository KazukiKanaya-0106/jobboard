import RefreshIcon from '@mui/icons-material/Refresh'
import {
  Alert,
  Box,
  CircularProgress,
  IconButton,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Toolbar,
  Typography,
} from '@mui/material'
import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../../auth/auth-context'
import { fetchJobs } from '../api'

export default function JobsPage() {
  const { auth } = useAuth()

  const jobsQuery = useQuery({
    queryKey: ['jobs'],
    queryFn: () => fetchJobs(auth!),
    enabled: Boolean(auth?.token),
  })

  return (
    <Paper elevation={0} sx={{ p: 3 }}>
      <Toolbar disableGutters sx={{ justifyContent: 'space-between', mb: 2 }}>
        <Typography variant="h5">ジョブ履歴</Typography>
        <IconButton onClick={() => jobsQuery.refetch()} disabled={jobsQuery.isFetching}>
          {jobsQuery.isFetching ? <CircularProgress size={20} /> : <RefreshIcon />}
        </IconButton>
      </Toolbar>

      {jobsQuery.isLoading ? (
        <Box sx={{ py: 8, textAlign: 'center' }}>
          <CircularProgress />
        </Box>
      ) : jobsQuery.isError ? (
        <Alert severity="error" sx={{ my: 4 }}>
          ジョブの取得に失敗しました
        </Alert>
      ) : (
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>ノードID</TableCell>
              <TableCell>ステータス</TableCell>
              <TableCell>開始</TableCell>
              <TableCell>終了</TableCell>
              <TableCell>所要時間 (h)</TableCell>
              <TableCell>タグ</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {jobsQuery.data && jobsQuery.data.length > 0 ? (
              jobsQuery.data.map((job) => (
                <TableRow key={job.id} hover>
                  <TableCell>{job.id}</TableCell>
                  <TableCell>{job.nodeId}</TableCell>
                  <TableCell sx={{ textTransform: 'capitalize' }}>{job.status}</TableCell>
                  <TableCell>{job.startedAt ? job.startedAt.toLocaleString() : '-'}</TableCell>
                  <TableCell>{job.finishedAt ? job.finishedAt.toLocaleString() : '-'}</TableCell>
                  <TableCell>{job.durationHours ? job.durationHours.toFixed(2) : '-'}</TableCell>
                  <TableCell>{job.tag ?? '-'}</TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={7}>
                  <Stack sx={{ py: 6, textAlign: 'center', color: 'text.secondary' }}>
                    ジョブ履歴がありません
                  </Stack>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      )}
    </Paper>
  )
}
