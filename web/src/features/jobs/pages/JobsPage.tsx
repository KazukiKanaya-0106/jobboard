import RefreshIcon from "@mui/icons-material/Refresh";
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Toolbar,
  Tooltip,
  Typography,
} from "@mui/material";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useAuth } from "../../auth/AuthContext";
import { fetchJobs, type Job } from "../api";

type StatusKey = "running" | "completed" | "failed" | "other";

const STATUS_STYLES: Record<StatusKey, { label: string; bgcolor: string; color: string }> = {
  running: {
    label: "running",
    bgcolor: "info.light",
    color: "info.contrastText",
  },
  completed: {
    label: "completed",
    bgcolor: "success.light",
    color: "success.contrastText",
  },
  failed: {
    label: "failed",
    bgcolor: "error.main",
    color: "error.contrastText",
  },
  other: {
    label: "unknown",
    bgcolor: "grey.600",
    color: "common.white",
  },
};

function StatusChip({ status, onClick }: { status: string; onClick?: () => void }) {
  const normalized = status.toLowerCase();
  let key: StatusKey = "other";

  if (normalized === "running" || normalized === "completed" || normalized === "failed") {
    key = normalized;
  }

  const config = STATUS_STYLES[key];

  return (
    <Chip
      label={key === "other" ? status : config.label}
      size="small"
      sx={{
        px: 1.5,
        fontWeight: 600,
        textTransform: "capitalize",
        bgcolor: config.bgcolor,
        color: config.color,
        cursor: onClick ? "pointer" : "default",
      }}
      onClick={onClick}
      clickable={Boolean(onClick)}
    />
  );
}

export default function JobsPage() {
  const { auth } = useAuth();
  const [selectedJob, setSelectedJob] = useState<Job | null>(null);

  const requireAuth = () => {
    if (!auth) {
      throw new Error("Authentication is required to load jobs");
    }
    console.log("Using auth:", auth);
    return auth;
  };

  const jobsQuery = useQuery({
    queryKey: ["jobs"],
    queryFn: () => fetchJobs(requireAuth()),
    enabled: Boolean(auth?.token),
  });

  const handleCloseErrorDialog = () => setSelectedJob(null);

  return (
    <Paper elevation={0} sx={{ p: 3 }}>
      <Toolbar disableGutters sx={{ justifyContent: "space-between", mb: 2 }}>
        <Typography variant="h5">ジョブ履歴</Typography>
        <IconButton onClick={() => jobsQuery.refetch()} disabled={jobsQuery.isFetching}>
          {jobsQuery.isFetching ? <CircularProgress size={20} /> : <RefreshIcon />}
        </IconButton>
      </Toolbar>

      {jobsQuery.isLoading ? (
        <Box sx={{ py: 8, textAlign: "center" }}>
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
              jobsQuery.data.map((job) => {
                const isFailed = job.status.toLowerCase() === "failed";
                const statusChip = (
                  <StatusChip status={job.status} onClick={isFailed ? () => setSelectedJob(job) : undefined} />
                );

                return (
                  <TableRow key={job.id} hover>
                    <TableCell>{job.id}</TableCell>
                    <TableCell>{job.nodeId}</TableCell>
                    <TableCell>
                      {isFailed ? (
                        <Tooltip title="エラー詳細を表示" arrow>
                          <span style={{ display: "inline-flex" }}>{statusChip}</span>
                        </Tooltip>
                      ) : (
                        statusChip
                      )}
                    </TableCell>
                    <TableCell>{job.startedAt ? job.startedAt.toLocaleString() : "-"}</TableCell>
                    <TableCell>{job.finishedAt ? job.finishedAt.toLocaleString() : "-"}</TableCell>
                    <TableCell>{job.durationHours != null ? job.durationHours.toFixed(2) : "-"}</TableCell>
                    <TableCell>{job.tag ?? "-"}</TableCell>
                  </TableRow>
                );
              })
            ) : (
              <TableRow>
                <TableCell colSpan={7}>
                  <Stack sx={{ py: 6, textAlign: "center", color: "text.secondary" }}>ジョブ履歴がありません</Stack>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      )}

      <Dialog open={Boolean(selectedJob)} onClose={handleCloseErrorDialog} fullWidth maxWidth="sm">
        <DialogTitle>ジョブ{selectedJob ? ` #${selectedJob.id}` : ""}のエラー詳細</DialogTitle>
        <DialogContent dividers>
          {selectedJob?.errorText ? (
            <Box
              component="pre"
              sx={{
                whiteSpace: "pre-wrap",
                fontFamily: "Roboto Mono, monospace",
                fontSize: 14,
                m: 0,
              }}
            >
              {selectedJob.errorText}
            </Box>
          ) : (
            <Typography color="text.secondary">エラー詳細は記録されていません。</Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseErrorDialog}>閉じる</Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
}
