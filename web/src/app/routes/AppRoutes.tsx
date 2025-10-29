import { Navigate, Route, Routes } from "react-router-dom";
import { BrowserRouter } from "react-router-dom";
import AuthLayout from "../../layouts/AuthLayout";
import DashboardLayout from "../../layouts/DashboardLayout";
import LoginPage from "../../features/auth/pages/LoginPage";
import RegisterPage from "../../features/auth/pages/RegisterPage";
import NodesPage from "../../features/nodes/pages/NodesPage";
import JobsPage from "../../features/jobs/pages/JobsPage";
import ProtectedRoute from "../../features/auth/components/ProtectedRoute";
import RedirectIfAuthenticated from "../../features/auth/components/RedirectIfAuthenticated";

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/auth"
          element={
            <RedirectIfAuthenticated>
              <AuthLayout />
            </RedirectIfAuthenticated>
          }
        >
          <Route index element={<Navigate to="login" replace />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />
        </Route>

        <Route
          path="/"
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="nodes" replace />} />
          <Route path="nodes" element={<NodesPage />} />
          <Route path="jobs" element={<JobsPage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
