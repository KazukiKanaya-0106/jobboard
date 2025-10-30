import { alpha, Box, Paper, Stack, Typography } from "@mui/material";
import { Outlet } from "react-router-dom";

export default function AuthLayout() {
  return (
    <Box
      sx={(theme) => ({
        minHeight: "100vh",
        position: "relative",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        overflow: "hidden",
        px: { xs: 3, sm: 6 },
        py: { xs: 6, sm: 8 },
        background: alpha(theme.palette.background.default, 0.9),
      })}
    >
      <Box
        sx={(theme) => ({
          position: "absolute",
          inset: 0,
          background: `linear-gradient(135deg, ${alpha(theme.palette.primary.light, 0.92)} 0%, ${alpha(theme.palette.primary.main, 0.94)} 55%, ${alpha(theme.palette.primary.dark, 0.96)} 100%)`,
          filter: "blur(24px)",
          transform: "scale(1.05)",
        })}
      />
      <Box
        sx={(theme) => ({
          position: "absolute",
          inset: { xs: "8%", md: "6%" },
          borderRadius: 32,
          background:
            "radial-gradient(circle at 18% 18%, rgba(255,255,255,0.22), transparent 52%), radial-gradient(circle at 82% 82%, rgba(255,255,255,0.18), transparent 45%)",
          border: `1px solid ${alpha(theme.palette.common.white, 0.08)}`,
          backdropFilter: "blur(18px)",
          boxShadow: "0 40px 120px rgba(57, 45, 118, 0.35)",
        })}
      />
      <Box
        sx={{
          position: "absolute",
          inset: 0,
          backgroundColor: "rgba(18, 10, 38, 0.25)",
        }}
      />
      <Box
        sx={{
          position: "relative",
          zIndex: 1,
          width: "100%",
          maxWidth: 1040,
        }}
      >
        <Box
          sx={{
            display: "grid",
            gridTemplateColumns: { xs: "1fr", md: "1.1fr 0.9fr" },
            gap: { xs: 4, md: 6 },
            alignItems: "stretch",
          }}
        >
          <Paper
            elevation={0}
            sx={(theme) => ({
              display: { xs: "none", md: "flex" },
              position: "relative",
              overflow: "hidden",
              borderRadius: 5,
              background: `linear-gradient(140deg, ${alpha(theme.palette.primary.dark, 0.9)}, ${alpha(theme.palette.primary.main, 0.85)})`,
              color: "common.white",
              px: { md: 6.5, lg: 8.5 },
              py: { md: 8.5, lg: 11 },
              boxShadow: "0 34px 90px rgba(41, 32, 102, 0.35)",
              backdropFilter: "blur(12px)",
              "&::before": {
                content: '""',
                position: "absolute",
                inset: 0,
                opacity: 0.14,
                background:
                  "radial-gradient(circle at 24% 28%, rgba(255,255,255,0.65), transparent 58%), radial-gradient(circle at 80% 78%, rgba(255,255,255,0.35), transparent 48%)",
              },
              "&::after": {
                content: '""',
                position: "absolute",
                right: -120,
                top: -140,
                width: 420,
                height: 420,
                background: alpha("#ffffff", 0.08),
                borderRadius: "42% 58% 46% 54% / 58% 34% 66% 42%",
                filter: "blur(1px)",
              },
            })}
          >
            <Stack spacing={5} sx={{ position: "relative", zIndex: 1, maxWidth: 440 }}>
              <Box>
                <Typography variant="overline" sx={{ letterSpacing: 4, opacity: 0.7 }}>
                  OPERATIONS DASHBOARD
                </Typography>
                <Typography variant="h3" component="h1" fontWeight={700} gutterBottom sx={{ mt: 1 }}>
                  Jobboard Hub
                </Typography>
                <Typography variant="h6" sx={{ opacity: 0.9, lineHeight: 1.6 }}>
                  ノードの状況とジョブの進行を一つの画面で。
                  <br />
                  運用をスムーズにつなぐハブを体験してください。
                </Typography>
              </Box>
              <Stack spacing={2.8}>
                {[
                  "リアルタイムなノードモニタリング",
                  "ジョブ進捗と履歴の可視化",
                  "トークンベースのセキュアなアクセス",
                ].map((item) => (
                  <Stack key={item} direction="row" spacing={2.4} alignItems="center">
                    <Box
                      component="span"
                      sx={{
                        width: 11,
                        height: 11,
                        borderRadius: "50%",
                        backgroundColor: alpha("#ffffff", 0.9),
                        boxShadow: "0 0 18px rgba(255,255,255,0.5)",
                      }}
                    />
                    <Typography variant="body1" sx={{ opacity: 0.92, letterSpacing: 0.2 }}>
                      {item}
                    </Typography>
                  </Stack>
                ))}
              </Stack>
            </Stack>
          </Paper>

          <Paper
            elevation={0}
            sx={(theme) => ({
              display: "flex",
              flexDirection: "column",
              justifyContent: "center",
              p: { xs: 4, sm: 5.5 },
              borderRadius: 5,
              border: `1px solid ${alpha(theme.palette.primary.main, 0.08)}`,
              boxShadow: "0 30px 90px rgba(18, 14, 40, 0.18)",
              background: alpha(theme.palette.background.paper, 0.96),
              backdropFilter: "blur(18px)",
            })}
          >
            <Stack spacing={4.5} sx={{ maxWidth: 430, width: "100%", mx: "auto" }}>
              <Box sx={{ display: { xs: "block", md: "none" } }}>
                <Typography variant="overline" sx={{ letterSpacing: 3, opacity: 0.7 }}>
                  OPERATIONS DASHBOARD
                </Typography>
                <Typography variant="h4" component="h1" fontWeight={700} gutterBottom sx={{ mt: 1 }}>
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
    </Box>
  );
}
