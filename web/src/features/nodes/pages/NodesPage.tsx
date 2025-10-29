import AddIcon from '@mui/icons-material/Add'
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline'
import RefreshIcon from '@mui/icons-material/Refresh'
import {
  Alert,
  Box,
  Button,
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
  Tooltip,
  Typography,
} from '@mui/material'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useAuth } from '../../auth/AuthContext'
import NodeCreateDialog from '../components/NodeCreateDialog'
import NodeDeleteDialog from '../components/NodeDeleteDialog'
import NodeTokenDialog from '../components/NodeTokenDialog'
import { createNode, deleteNode, fetchNodes, type Node } from '../api'
import type { CreateNodeRequest } from '../schemas'

export default function NodesPage() {
  const { auth } = useAuth()
  const queryClient = useQueryClient()
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [createdToken, setCreatedToken] = useState<string | null>(null)
  const [createError, setCreateError] = useState<string | null>(null)
  const [nodeToDelete, setNodeToDelete] = useState<Node | null>(null)
  const [deleteError, setDeleteError] = useState<string | null>(null)

  const nodesQuery = useQuery({
    queryKey: ['nodes'],
    queryFn: () => fetchNodes(auth!),
    enabled: Boolean(auth?.token),
  })

  const createMutation = useMutation({
    mutationFn: (input: CreateNodeRequest) => createNode(auth!, input),
    onSuccess: (result) => {
      setIsCreateOpen(false)
      if (result.token) {
        setCreatedToken(result.token)
      }
      queryClient.invalidateQueries({ queryKey: ['nodes'] })
    },
    onError: (error: unknown) => {
      setCreateError(error instanceof Error ? error.message : 'ノードの作成に失敗しました')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (nodeId: number) => deleteNode(auth!, nodeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['nodes'] })
      setNodeToDelete(null)
    },
    onError: (error: unknown) => {
      setDeleteError(error instanceof Error ? error.message : 'ノードの削除に失敗しました')
    },
  })

  const handleCreate = async (input: CreateNodeRequest) => {
    setCreateError(null)
    await createMutation.mutateAsync(input)
  }

  const handleRefresh = () => {
    nodesQuery.refetch()
  }

  const handleDeleteConfirm = async () => {
    if (!nodeToDelete) return
    setDeleteError(null)
    await deleteMutation.mutateAsync(nodeToDelete.id)
  }

  return (
    <Stack spacing={3}>
      <Paper elevation={0} sx={{ p: 3 }}>
        <Toolbar disableGutters sx={{ justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h5">ノード一覧</Typography>
          <Stack direction="row" spacing={2}>
            <IconButton onClick={handleRefresh} disabled={nodesQuery.isFetching}>
              {nodesQuery.isFetching ? <CircularProgress size={20} /> : <RefreshIcon />}
            </IconButton>
            <Button variant="contained" startIcon={<AddIcon />} onClick={() => setIsCreateOpen(true)}>
              ノードを追加
            </Button>
          </Stack>
        </Toolbar>

        {nodesQuery.isLoading ? (
          <Box sx={{ py: 8, textAlign: 'center' }}>
            <CircularProgress />
          </Box>
        ) : nodesQuery.isError ? (
          <Alert severity="error" sx={{ my: 4 }}>
            ノードの取得に失敗しました
          </Alert>
        ) : (
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>ノード名</TableCell>
                <TableCell>現在のジョブID</TableCell>
                <TableCell>作成日時</TableCell>
                <TableCell align="right">操作</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {nodesQuery.data && nodesQuery.data.length > 0 ? (
                nodesQuery.data.map((node) => (
                  <TableRow key={node.id} hover>
                    <TableCell>{node.id}</TableCell>
                    <TableCell>{node.nodeName}</TableCell>
                    <TableCell>{node.currentJobId ?? 'なし'}</TableCell>
                    <TableCell>{node.createdAt.toLocaleString()}</TableCell>
                    <TableCell align="right">
                      <Tooltip title="ノードを削除">
                        <span>
                          <IconButton
                            color="error"
                            onClick={() => {
                              setNodeToDelete(node)
                              setDeleteError(null)
                            }}
                            disabled={deleteMutation.isPending && nodeToDelete?.id === node.id}
                          >
                            <DeleteOutlineIcon />
                          </IconButton>
                        </span>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={5}>
                    <Box sx={{ py: 6, textAlign: 'center', color: 'text.secondary' }}>ノードがありません</Box>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        )}
      </Paper>

      <NodeCreateDialog
        open={isCreateOpen}
        onClose={() => {
          setIsCreateOpen(false)
          setCreateError(null)
        }}
        onSubmit={handleCreate}
        loading={createMutation.isPending}
        apiError={createError}
      />

      <NodeTokenDialog
        token={createdToken}
        onClose={() => {
          setCreatedToken(null)
        }}
      />

      <NodeDeleteDialog
        open={Boolean(nodeToDelete)}
        nodeName={nodeToDelete?.nodeName}
        onClose={() => {
          if (deleteMutation.isPending) return
          setNodeToDelete(null)
          setDeleteError(null)
        }}
        onConfirm={handleDeleteConfirm}
        loading={deleteMutation.isPending}
        error={deleteError}
      />
    </Stack>
  )
}
