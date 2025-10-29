import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  IconButton,
  InputAdornment,
  Tooltip,
  Snackbar,
  Alert,
} from "@mui/material";
import { useState } from "react";

type NodeTokenDialogProps = {
  token: string | null;
  onClose: () => void;
};

export default function NodeTokenDialog({ token, onClose }: NodeTokenDialogProps) {
  const [copied, setCopied] = useState(false);

  if (!token) return null;

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(token);
      setCopied(true);
    } catch {
      setCopied(false);
    }
  };

  return (
    <>
      <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
        <DialogTitle>NodeToken を保存してください</DialogTitle>
        <DialogContent dividers>
          <TextField
            label="NodeToken"
            value={token}
            fullWidth
            InputProps={{
              readOnly: true,
              endAdornment: (
                <InputAdornment position="end">
                  <Tooltip title="クリップボードにコピー">
                    <IconButton onClick={handleCopy}>
                      <ContentCopyIcon />
                    </IconButton>
                  </Tooltip>
                </InputAdornment>
              ),
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} variant="contained">
            閉じる
          </Button>
        </DialogActions>
      </Dialog>
      <Snackbar
        open={copied}
        autoHideDuration={2000}
        onClose={() => setCopied(false)}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
      >
        <Alert severity="success" variant="filled">
          NodeToken をコピーしました
        </Alert>
      </Snackbar>
    </>
  );
}
