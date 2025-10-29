import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { ReactNode } from 'react'
import { AuthProvider } from '../../features/auth/AuthContext'

type AppProvidersProps = {
  children: ReactNode
}

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 1000 * 10,
    },
  },
})

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      light: '#b8b0f0',
      main: '#6958d3',
      dark: '#4539a3',
      contrastText: '#ffffff',
    },
    secondary: {
      light: '#f5cde2',
      main: '#e3a6c9',
      dark: '#b97da0',
      contrastText: '#3f1f2f',
    },
    background: {
      default: '#f5f3fb',
      paper: '#ffffff',
    },
    text: {
      primary: '#201c3a',
      secondary: '#5a547b',
    },
  },
  typography: {
    fontFamily:
      "'Inter', 'Noto Sans JP', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Helvetica Neue', 'Arial', sans-serif",
    h1: {
      fontWeight: 700,
      letterSpacing: '-0.03em',
    },
    h2: {
      fontWeight: 700,
      letterSpacing: '-0.025em',
    },
    h3: {
      fontWeight: 600,
      letterSpacing: '-0.02em',
    },
    h4: {
      fontWeight: 600,
      letterSpacing: '-0.015em',
    },
    button: {
      fontWeight: 600,
      textTransform: 'none',
      letterSpacing: '-0.01em',
    },
    subtitle1: {
      fontWeight: 500,
      letterSpacing: '-0.01em',
    },
  },
  shape: {
    borderRadius: 16,
  },
  components: {
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 22,
          boxShadow:
            '0 14px 32px rgba(69, 56, 163, 0.14), 0 6px 16px rgba(105, 88, 211, 0.12), inset 0 0 0 1px rgba(105, 88, 211, 0.05)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 999,
          paddingInline: 26,
          paddingBlock: 12,
          boxShadow: '0 10px 22px rgba(69, 56, 163, 0.22)',
          transition: 'transform 0.18s ease, box-shadow 0.18s ease',
          '&:hover': {
            transform: 'translateY(-1px)',
            boxShadow: '0 12px 26px rgba(69, 56, 163, 0.28)',
          },
        },
      },
    },
  },
})

export default function AppProviders({ children }: AppProvidersProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <AuthProvider>{children}</AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  )
}
