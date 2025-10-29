import { Box, Paper, Stack, Typography } from "@mui/material";
import { Outlet } from "react-router-dom";

export default function AuthLayout() {
  return (
    <Box
      sx={(theme) => ({
        minHeight: "100vh",
        display: "grid",
        gridTemplateColumns: {
          xs: "1fr",
          md: "1.05fr 0.95fr",
        },
        background: `linear-gradient(135deg, ${theme.palette.primary.light} 0%, ${theme.palette.primary.main} 60%, ${theme.palette.primary.dark} 100%)`,
      })}
    >
      <Box
        sx={{
          position: "relative",
          display: { xs: "none", md: "flex" },
          flexDirection: "column",
          justifyContent: "center",
          px: { md: 8, lg: 12 },
          py: { md: 12, lg: 16 },
          color: "common.white",
          background: (theme) =>
            `linear-gradient(135deg, ${theme.palette.primary.dark}, ${theme.palette.primary.main})`,
          overflow: "hidden",
          "&::before": {
            content: '""',
            position: "absolute",
            inset: 0,
            opacity: 0.18,
            background:
              "radial-gradient(circle at 20% 20%, rgba(255,255,255,0.6), transparent 55%), radial-gradient(circle at 80% 80%, rgba(255,255,255,0.4), transparent 45%)",
          },
          "&::after": {
            content: '""',
            position: "absolute",
            top: -80,
            right: -120,
            width: 360,
            height: 360,
            background: "rgba(255,255,255,0.08)",
            borderRadius: "36% 64% 61% 39% / 43% 42% 58% 57%",
            filter: "blur(0.5px)",
          },
        }}
      >
        <Stack spacing={5} sx={{ position: "relative", maxWidth: 420 }}>
          <div>
            <Typography variant="h3" component="h1" fontWeight={700} gutterBottom>
              Jobboard Hub
            </Typography>
            <Typography variant="h6" sx={{ opacity: 0.9, lineHeight: 1.5 }}>
              クラスタとジョブ管理をスマートに。運用に集中できるダッシュボード体験を提供します。
            </Typography>
          </div>

          <Stack spacing={2.5}>
            {["リアルタイムなノードモニタリング", "ジョブの進行状況を可視化", "トークンベースのセキュアなアクセス"].map(
              (item) => (
                <Stack key={item} direction="row" spacing={2} alignItems="center">
                  <Box
                    component="span"
                    sx={{
                      width: 10,
                      height: 10,
                      borderRadius: "50%",
                      backgroundColor: "rgba(255,255,255,0.85)",
                      boxShadow: "0 0 12px rgba(255,255,255,0.45)",
                    }}
                  />
                  <Typography variant="body1" sx={{ opacity: 0.85 }}>
                    {item}
                  </Typography>
                </Stack>
              ),
            )}
          </Stack>
        </Stack>
      </Box>

      <Box
        sx={(theme) => ({
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          px: { xs: 2, sm: 4, md: 6 },
          py: { xs: 6, sm: 8, md: 10 },
          background: "transparent",
          [theme.breakpoints.up("md")]: {
            background: `linear-gradient(140deg, ${theme.palette.primary.light}1c 0%, ${theme.palette.background.default} 40%, ${theme.palette.primary.main}14 100%)`,
            backdropFilter: "blur(20px)",
          },
        })}
      >
        <Paper
          elevation={12}
          sx={(theme) => ({
            width: "100%",
            maxWidth: 440,
            p: { xs: 4, sm: 6 },
            borderRadius: 4,
            border: `1px solid ${theme.palette.primary.main}22`,
            boxShadow: theme.shadows[10],
            background: theme.palette.background.paper,
          })}
        >
          <Stack spacing={4}>
            <Box sx={{ display: { xs: "block", md: "none" } }}>
              <Typography variant="h4" component="h1" fontWeight={700} gutterBottom>
                Jobboard Hub
              </Typography>
              <Typography variant="subtitle1" color="text.secondary">
                クラスター管理ダッシュボードへようこそ
              </Typography>
            </Box>
            <Outlet />
          </Stack>
        </Paper>
      </Box>
    </Box>
  );
}
